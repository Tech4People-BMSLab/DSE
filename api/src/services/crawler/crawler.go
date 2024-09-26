package crawler

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"bms.dse/src/debug"
	"bms.dse/src/services/db"
	"bms.dse/src/services/ipc"
	"bms.dse/src/utils/gatekeeper"
	"bms.dse/src/utils/httputil"
	"bms.dse/src/utils/logging"
	"bms.dse/src/utils/semaphore"
	"github.com/nats-io/nats.go"
	"github.com/robfig/cron/v3"
	"github.com/tidwall/gjson"
	"github.com/vmihailenco/msgpack/v5"
)

// ------------------------------------------------------------
// : Aliases
// ------------------------------------------------------------
type JSON    = map[string]any
type Writer  = http.ResponseWriter
type Request = http.Request

// ------------------------------------------------------------
// : Types
// ------------------------------------------------------------
type Website struct {
	Name string
	URL  string
}

type Keyword struct {
	Keyword string
}

type Payload struct {
	Timestamp    string `json:"timestamp"    msgpack:"timestamp"`
	Token        string `json:"token"        msgpack:"token"`
	URL          string `json:"url"          msgpack:"url"`
	Browser      JSON   `json:"browser"      msgpack:"browser"`
	Localization string `json:"localization" msgpack:"localization"`
	Keyword      string `json:"keyword"      msgpack:"keyword"`
	Website      string `json:"website"      msgpack:"website"`
	HTML         string `json:"html"         msgpack:"html"`
}
// ------------------------------------------------------------
// : Constants
// ------------------------------------------------------------
const (
	USER_PING      = "bms.dse.users.%s.ping"
	CRAWLER_SCRAPE = "bms.dse.%s.crawler.scrape"
	CRAWLER_START  = "bms.dse.%s.crawler.start"
	CRAWLER_STOP   = "bms.dse.%s.crawler.stop"
)

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var logger = logging.NewLogger("crawler")
var mUsers = make(map[string]*User)
var lWebsites = []Website{
	{Name: "Google"    , URL: "https://www.google.com/search?q=%s"}, 
	{Name: "DuckDuckGo", URL: "https://duckduckgo.com/?q=%s"},
	{Name: "Bing"      , URL: "https://www.bing.com/search?q=%s"},
}
var lKeywords = []Keyword{
	{Keyword: "vluchtelingen"},
	{Keyword: "migratie in nederland"},
	{Keyword: "asiel en migratiebeleid"},
}
var wg  = sync.WaitGroup{}
var sem = semaphore.New(100)

var gkUsersLoaded = gatekeeper.NewGateKeeper(true)
// ------------------------------------------------------------
// : Events
// ------------------------------------------------------------
func OnConnect(m *nats.Msg) {
	if debug.IsDebugMode() { return }
	go func() {
		payload := gjson.ParseBytes(m.Data)
		user    := &User{}

		token := payload.Get("token").String()
		_, ok := mUsers[token]

		if !ok {
			user.Token = token
			user.Load()
			mUsers[user.GetID()] = user
		} else {
			user = mUsers[token]
			user.Token = token
			user.Load()
		}

		user.OnHeartbeat(m)
	}()
}

func OnPing(m *nats.Msg) {
	if debug.IsDebugMode() { return }
	go func() {
		payload := gjson.ParseBytes(m.Data)
		user    := &User{}

		token := payload.Get("token").String()
		_, ok := mUsers[token]

		if !ok {
			user.Token = token
			user.Load()
			mUsers[user.GetID()] = user
		} else {
			user = mUsers[token]
			user.Token = token
			user.Load()
		}

		user.OnHeartbeat(m)
	}()
	
}

func OnUpload(m *nats.Msg) {
	if debug.IsDebugMode() { return }
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error().Interface("recover", r).Msg("Recovered")
			}
		}()

		sem.Acquire()
		defer sem.Release()
		
		payload := gjson.ParseBytes(m.Data)
		uToken  := payload.Get("token").String()
		user    := mUsers[uToken]
		
		taskURL          := payload.Get("url").String()
		taskParsedURL, _ := url.Parse(taskURL)
		taskID           := taskParsedURL.Query().Get("dse_id")

		hasSent := false
		for _, t := range user.Tasks {
			if t.ID == taskID {
				t.Uploaded <- m.Data
				hasSent = true
				break
			}
		}

		if !hasSent {
			logger.Warn().Str("token", uToken).Str("id", taskID).Msg("Task not found")
		}

		m.Respond([]byte("{}"))
	}()
}

func HandleUpload(w Writer, r *Request) {

	// Get data from request
	data, err := io.ReadAll(r.Body)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, JSON{
			"error": "FAILED_TO_READ_REQUEST_BODY",
		})
		return
	}
	defer r.Body.Close()

	// Create payload
	payload := &Payload{}

	// Decode payload
	err = msgpack.Unmarshal(data, payload)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, JSON{
			"error": "FAILED_TO_UNMARSHAL_PAYLOAD",
		})
		return
	}

	// Lock the semaphore
	sem.Acquire()
	defer sem.Release()

	// Get the user
	user := mUsers[payload.Token]

	// Get the task ID
	taskURL          := payload.URL
	taskParsedURL, _ := url.Parse(taskURL)
	taskID           := taskParsedURL.Query().Get("dse_id")

	// Convert data into bytes using JSON
	data, err = json.Marshal(payload)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, JSON{
			"error": "FAILED_TO_MARSHAL_PAYLOAD",
		})
	}
	
	// Find task
	hasSent := false
	for _, t := range user.Tasks {
		if t.ID == taskID {
			t.Uploaded <- data
			hasSent = true
			break
		}
	}

	if !hasSent {
		logger.Warn().Str("token", payload.Token).Str("id", taskID).Msg("Task not found")
	}

	httputil.WriteJSON(w, http.StatusOK, JSON{
		"status": "OK",
	})
}

func OnUserStart(m *nats.Msg) {
	go func() {
		payload := gjson.ParseBytes(m.Data)
		user    := &User{}

		token := payload.Get("token").String()
		_, ok := mUsers[token]

		if !ok {
			user.Token = token
			user.Load()
			mUsers[user.GetID()] = user
		} else {
			user = mUsers[token]
			user.Token = token
			user.Load()
		}

		user.OnHeartbeat(m)
		go user.Start()

		logger.Info().Str("token", token).Msg("User Test")

		m.Respond([]byte("{}"))
	}()
}

func OnUserStop(m *nats.Msg) {
	go func() {
		payload := gjson.ParseBytes(m.Data)
		user    := &User{}

		token := payload.Get("token").String()
		_, ok := mUsers[token]

		if !ok {
			user.Token = token
			user.Load()
			mUsers[user.GetID()] = user
		} else {
			user = mUsers[token]
			user.Token = token
			user.Load()
		}

		user.OnHeartbeat(m)
		go user.Stop()

		m.Respond([]byte("{}"))
	}()
}

// ------------------------------------------------------------
// : Functions
// ------------------------------------------------------------
func GetUsers() map[string]*User {
	gkUsersLoaded.Wait()
	
	return mUsers
}
// ------------------------------------------------------------
// : Initializers
// ------------------------------------------------------------
func InitDB() {
	defer wg.Done()
	logger.Info().Msg("Initializing database")
	
	if debug.IsDebugMode() { return }

	db.AutoMigrate(&User{})
}

func InitEvents() {
	defer wg.Done()
	logger.Info().Msg("Initializing events")
	ipc.Subscribe("bms.dse.users.*.connected" , OnConnect)
	ipc.Subscribe("bms.dse.users.connected"   , OnConnect)
	ipc.Subscribe("bms.dse.users.pong"        , OnPing)
	ipc.Subscribe("bms.dse.v1.5.search.upload", OnUpload)
	ipc.Subscribe("bms.dse.v1.5.crawler.start", OnUserStart)
	ipc.Subscribe("bms.dse.v1.5.crawler.stop" , OnUserStop)
}

func InitPinger() {
	defer wg.Done()
	logger.Info().Msg("Initializing pinger")

	// if debug.IsDebugMode() { return } // Comment this if needed during debugging

	scheduler := cron.New()

	users := []*User{}
	db.Table("users").Where("connected = ?", true).Update("connected", false)
	db.Table("users").Find(&users)

	for _, u := range users {
		mUsers[u.GetID()] = u
	}

	gkUsersLoaded.Unlock()
	
	ipc.Publish("bms.dse.users.ping", []byte("{}"))
	scheduler.AddFunc("@every 1m", func() {
		ipc.Publish("bms.dse.users.ping", []byte("{}"))
	})

	scheduler.AddFunc("@every 2m", func() {
		for _, user := range mUsers {
			if time.Since(user.LastPing) > time.Minute {
				if !user.Connected { continue }

				user.Connected = false
				user.Save()
			}
		}
	})

	scheduler.Start()
}
// ------------------------------------------------------------
// : Debug
// ------------------------------------------------------------
func Debug() {
	if !debug.IsDebugMode() { return }
}

// ------------------------------------------------------------
// : Init
// ------------------------------------------------------------
func Init() {
	wg.Add(3)
	go InitDB()
	go InitEvents()
	go InitPinger()
	wg.Wait()

	go Debug()
}
