package crawler

import (
	"fmt"
	"net/url"
	"sync"
	"time"

	"bms.dse/src/services/db"
	"bms.dse/src/services/ipc"
	"bms.dse/src/services/parser"
	"github.com/Masterminds/semver/v3"
	"github.com/golang-module/carbon"
	"github.com/jackc/pgtype"
	"github.com/nats-io/nats.go"
	"github.com/rs/xid"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"gorm.io/gorm"
)

// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
type JSONB = pgtype.JSONB
type Time  = time.Time
// ------------------------------------------------------------
// : User
// ------------------------------------------------------------
type User struct {
	Token       string        `json:"token"      gorm:"primaryKey;unique;notNull;type:text"`
	Connected   bool          `json:"connected"  gorm:"type:boolean"`
	Form        JSONB         `json:"form"       gorm:"type:jsonb"`
	Manifest    JSONB         `json:"manifest"   gorm:"type:jsonb"`
	Storage     JSONB         `json:"storage"    gorm:"type:jsonb"`
	Version     string        `json:"version"    gorm:"type:text"`

	LastPing    Time          `json:"last_ping"   gorm:"type:timestamptz"`
	LastSearch  Time          `json:"last_search" gorm:"type:timestamptz"`

	Tasks       []*Task       `json:"-"          gorm:"-"`
	Mutex       sync.Mutex    `json:"-"          gorm:"-"`
}

type Task struct {
	ID          string      `json:"id"         gorm:"primaryKey;unique;notNull;type:text"`
	User        string		`json:"user"       gorm:"type:text"`

	Status	    string      `json:"status"     gorm:"type:text"`

	URL		    string      `json:"url"        gorm:"type:text"`
	Keyword     string      `json:"keyword"    gorm:"type:text"`
	Website     string      `json:"website"    gorm:"type:text"`

	CreatedAt   Time        `json:"createdAt"  gorm:"type:timestamptz"`
	StartedAt   Time        `json:"startedAt"  gorm:"type:timestamptz"`
	FinishedAt  Time        `json:"finishedAt" gorm:"type:timestamptz"`

	Uploaded chan []byte    `json:"-"          gorm:"-"`
}

type Search struct {
    ID           uint        `gorm:"primaryKey;autoIncrement" json:"id"`
    User         string      `gorm:"type:text"                json:"user"`
    Timestamp    time.Time   `gorm:"type:time stamptz"        json:"timestamp"`
    Browser      JSONB       `gorm:"type:jsonb"               json:"browser"`
    Localization string      `gorm:"type:text"                json:"localization"`
    URL          string      `gorm:"type:text"                json:"url"`
    Keyword      string      `gorm:"type:text"                json:"keyword"`
    Website      string      `gorm:"type:text"                json:"website"`
    Results      JSONB       `gorm:"type:jsonb"               json:"results"`
}

// ------------------------------------------------------------
// : Methods
// ------------------------------------------------------------
func (u *User) GetID() string {
	return u.Token
}

func (user *User) IsConnected() bool {
	subject := fmt.Sprintf(USER_PING, user.GetID())
	_, err := ipc.Request(subject, []byte("{}"), 10 * time.Second)
	if err != nil {
		user.Connected = false
	} else {
		user.Connected = true 
	}

	return user.Connected
}

func (u *User) HasAcceptedVersion() bool {
    vBase, err := semver.NewVersion("1.0.162")
    if err != nil {
        return false
    }
 
    uBase, err := semver.NewVersion(u.Version)
    if err != nil {
        return false
    }

	return uBase.GreaterThan(vBase)
}

func (u *User) HasSearchedThisWeek() bool { 
	timeMonday := carbon.Now().SetWeekStartsAt(carbon.Monday).StartOfWeek()
	timeSearch := carbon.FromStdTime(u.LastSearch)

	return timeSearch.Gte(timeMonday)
}

func (u *User) HasTasks() bool {
	return len(u.Tasks) > 0
}

func (u *User) Load() {
	db.Table("users").Where("token = ?", u.Token).First(u)
}

func (u *User) Save() {
    if len(u.Manifest.Bytes) == 0 { u.Manifest.Status = pgtype.Null }
    if len(u.Form.Bytes)     == 0 { u.Form.Status     = pgtype.Null }
    if len(u.Storage.Bytes)  == 0 { u.Storage.Status  = pgtype.Null }

	userExisting := User{}
    tx := db.Table("users").Where("token = ?", u.Token).First(&userExisting)

    if tx.Error != nil {
        if tx.Error == gorm.ErrRecordNotFound {
            tx = db.Table("users").Create(u)
        } else {
            logger.Error().Err(tx.Error).Msg("Failed to search for user")
            return
        }
    } else {
        tx = db.Table("users").Save(u)
    }

    if tx.Error != nil {
        logger.Error().Err(tx.Error).Msg("User failed to save")
        return
    }
}

// ------------------------------------------------------------
// : Events
// ------------------------------------------------------------
func (u *User) OnHeartbeat(m *nats.Msg) {
	payload := gjson.ParseBytes(m.Data)
	
	u.Version    = payload.Get("manifest.version").String()
	u.Manifest   = pgtype.JSONB{Bytes: []byte(payload.Get("manifest").Raw), Status: pgtype.Present}
	u.Form		 = pgtype.JSONB{Bytes: []byte(payload.Get("form").Raw)    , Status: pgtype.Present}
	u.Storage	 = pgtype.JSONB{Bytes: []byte(payload.Get("storage").Raw) , Status: pgtype.Present}

	u.Connected  = true
	u.LastPing   = time.Now()
	u.Save()

	time.Sleep(3 * time.Second)
	 
	logger.Info().
	Str("token"       , u.GetID()).
	Str("version"     , u.Version).
	Time("last_search", u.LastSearch).
	Msg("User Heartbeat")

	go func() {
		u.Mutex.Lock()
		defer u.Mutex.Unlock()

		if !u.IsConnected()         { return }
		if !u.HasAcceptedVersion()  { return }
		if  u.HasSearchedThisWeek() { return }

		logger.Info().
		Str("token"              , u.GetID()).
		Str("version"            , u.Version).
		Time("last_search"       , u.LastSearch).
		Bool("connected"         , u.Connected).
		Bool("version_accepted"  , u.HasAcceptedVersion()).
		Bool("has_searched"      , u.HasSearchedThisWeek()).
		Msg("User Ready to Start")

		err := u.Start()
		if err != nil {
			logger.Error().Err(err).Msg("User failed to start")
			u.LastSearch = time.Time{}
			return
		}
	}()
}
// ------------------------------------------------------------
// : Controller
// ------------------------------------------------------------
func (u *User) Start() error {
	var err error

	// Close user window (if opened)
	subject   := fmt.Sprintf(CRAWLER_STOP, u.GetID())
	_, _ = ipc.Request(subject, []byte("{}"), 60 * time.Second)

	// Generate task from keywords and websites (if no tasks)
	if u.HasTasks() == false {
		logger.Debug().
		Str("user"   , u.GetID()).
		Str("version", u.Version).
		Msg("Generating tasks")

		for _, k := range lKeywords {
			for _, w := range lWebsites {
				t := Task{
					ID:        xid.New().String(),
					User:      u.Token,
	
					Status:   "pending",
	
					URL:       w.URL,
					Keyword:   k.Keyword,
					Website:   w.Name,
	
					CreatedAt : time.Now(),
					StartedAt : time.Time{},
					FinishedAt: time.Time{},

					Uploaded: make(chan []byte, 1),
				}
	
				logger.Debug().
				Str("id"     , t.ID).
				Str("user"   , t.User).
				Str("keyword", t.Keyword).
				Str("website", t.Website).
				Msg("Task generated")
	
				u.Tasks = append(u.Tasks, &t)
			}
		}
	}

	// Start crawler engine on user
	subject = fmt.Sprintf(CRAWLER_START, u.GetID())
	_, err  = ipc.Request(subject, []byte("{}"), 60 * time.Second)
	if err != nil {
		logger.Error().
		Str("user"   , u.GetID()).
		Str("version", u.Version).
		Err(err).
 		Msg("Crawler engine failed to start") 
		return err
	}

	// Start all tasks
	for _, task := range u.Tasks {
		if task.Status != "pending" { continue }
		if task.Status == "skip"    { continue }

		task.StartedAt = time.Now()
		err := u.Scrape(task)
		if err != nil {
			task.Status = "failed"
			logger.Error().
			Str("user"   , u.GetID()).
			Str("version", u.Version).
			Str("task"   , task.ID).
			Str("url"    , task.URL).
			Str("website", task.Website).
			Str("keyword", task.Keyword).
			Err(err).
			Msg("Scraping failed") 
			continue
		}

		select {
			case data := <-task.Uploaded: {
				payload      := gjson.ParseBytes(data)
				html         := payload.Get("html").String()
				results, err := parser.Parse(task.URL, html)

				if err != nil {
					logger.Error().
					Err(err).
					Str("user"   , u.GetID()).	
					Str("task"   , task.ID).
					Str("url"    , task.URL).
					Str("website", task.Website).
					Str("keyword", task.Keyword).
					Msg("Parser failed")
					task.Status = "failed"
					continue
				}

				search := &Search{}
				search.Timestamp = time.Now()
				search.User		 = u.GetID()

				search.Localization = payload.Get("localization").String()
				search.URL		    = payload.Get("url").String()
				search.Keyword	    = payload.Get("keyword").String()
				search.Website	    = payload.Get("website").String()

				search.Browser = pgtype.JSONB{Bytes: []byte(payload.Get("browser").Raw), Status: pgtype.Present}
				search.Results = pgtype.JSONB{Bytes: []byte(results), Status: pgtype.Present}

				task.FinishedAt = time.Now()
				task.Status	    = "completed"

				db.Table("searches").Save(search)

				logger.Info().
				Str("user"     , u.GetID()). 
				Str("task"     , task.ID).
				Str("url"      , task.URL).
				Str("website"  , task.Website).
				Str("keyword"  , task.Keyword).
				Msg("Task Complete")
			}

			case <-time.After(60 * time.Second): {
				task.Status = "failed"
				logger.Error().
				Str("user"   , u.GetID()).	
				Str("task"   , task.ID).
				Str("url"    , task.URL).
				Str("website", task.Website).
				Str("keyword", task.Keyword).
				Msg("Scraping timeout")
			}
		}
	}

	u.Tasks      = []*Task{}
	u.LastSearch = time.Now()
	u.Save()

	// Stop crawler engine on user
	subject = fmt.Sprintf(CRAWLER_STOP, u.GetID())
	_, err  = ipc.Request(subject, []byte("{}"), 60 * time.Second)
	if err != nil {
		logger.Error().Err(err).Msg("Crawler engine failed to stop")
		return err
	}

	logger.Info().
	Str("user"     , u.GetID()).
	Msg("Search Complete")

	return nil
}

func (u *User) Scrape(task *Task) error { 
	taskURL, _    := url.Parse(task.URL)
	taskParsedURL := taskURL.Query()
	taskParsedURL.Set("q"		   , task.Keyword)
	taskParsedURL.Add("dse"        , "1")
	taskParsedURL.Add("dse_id"	   , task.ID)
	taskParsedURL.Add("dse_keyword", task.Keyword)
	taskParsedURL.Add("dse_website", task.Website)
	taskURL.RawQuery = taskParsedURL.Encode()

	payload   := ""
	payload, _ = sjson.Set(payload, "id"       , task.ID)
	payload, _ = sjson.Set(payload, "timestamp", time.Now().Format(time.RFC3339))
	payload, _ = sjson.Set(payload, "keyword"  , task.Keyword)
	payload, _ = sjson.Set(payload, "website"  , task.Website)
	payload, _ = sjson.Set(payload, "url"      , taskURL.String())

	logger.Info().
	Str("user", u.GetID()).
	Str("task", task.ID).
	Str("url" , taskURL.String()).
	Str("website", task.Website).
	Str("keyword", task.Keyword).
	Msg("Scraping")

	s := fmt.Sprintf(CRAWLER_SCRAPE, u.GetID())
	ipc.Publish(s, []byte(payload))
	return nil
}

func (u *User) Stop() {
	for _, task := range u.Tasks {
		task.Status = "skip"
	}

	subject := fmt.Sprintf(CRAWLER_STOP, u.GetID())
	_, err  := ipc.Request(subject, []byte("{}"), 60 * time.Second)
	if err != nil {
		logger.Error().Err(err).Msg("Crawler engine failed to stop")
	}
}
