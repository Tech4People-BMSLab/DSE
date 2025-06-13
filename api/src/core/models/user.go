package models

import (
	"dse/src/utils/datetime"
	"dse/src/utils/event"
	"dse/src/utils/gatekeeper"
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/natefinch/lumberjack"
	"github.com/olebedev/emitter"
	"github.com/rs/zerolog"
)

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	version string = "3.0.5"
)

// ------------------------------------------------------------
// : User
// ------------------------------------------------------------
type User struct {
	Token   string         `json:"token"`
	State   State          `json:"state"`
	
	version string `json:"-"` // Version of the extensionm
	

	logger *zerolog.Logger `json:"-"` // Logger for user

	w  http.ResponseWriter `json:"-"` // ResponseWriter
	r *http.Request        `json:"-"` // Request
	conn *websocket.Conn   `json:"-"` // WS Connection

	mutex          sync.Mutex  `json:"-"` // Mutex for user
	flag_saving    atomic.Bool `json:"-"` // Saving flag
	flag_crawling  atomic.Bool `json:"-"` // Crawling flag
	flag_resetting atomic.Bool `json:"-"` // Resetting flag
}
// ------------------------------------------------------------
// : Init
// ------------------------------------------------------------
func (u *User) Init() {
	if u.State == (State{}) {
		u.State = *NewState()
	}

	if u.State.Client == nil {
		u.State.Client = NewClientState()
	}

	if u.State.Server == nil {
		u.State.Server = NewServerState()
	}
	
	u.State.Server.GK = gatekeeper.NewGateKeeper(false)
	u.State.Server.Online = false
	u.State.Client.Crawler.mutex = sync.Mutex{}
	u.State.Client.Crawler.cond  = *sync.NewCond(&u.State.Client.Crawler.mutex)

	u.mutex        = sync.Mutex{}
	u.flag_saving  = atomic.Bool{}

	event.On(fmt.Sprintf("user.%s.reload", u.Token), func(e *emitter.Event) { })
}

// ------------------------------------------------------------
// : Communication
// ------------------------------------------------------------
func (u *User) SetSSE(w http.ResponseWriter, r *http.Request) {
	u.w = w
	u.r = r

	go func() {
		version := r.URL.Query().Get("version")
		u.version = version
		u.logger  = newLogger(u)
		u.logger.Info().Msg("Connected SSE")

		var server = u.State.Server
		
		server.Online   = true
		server.LastPing = datetime.ToISO(datetime.Now())
		u.Save()
	
		 <- r.Context().Done()
	
		u.logger.Info().Msg("Disconnected SSE")
	
		u.w = nil
		u.r = nil
	
		server.Online = false
		u.Save()
		
	}()
}

func (u *User) SetWS(conn *websocket.Conn) {
	if u.conn != nil { return }
	u.conn = conn

	go func() {
		if u.logger != nil {
			u.logger.Info().Msg("Connected WS")
		}

		var server = u.State.Server

		server.Online   = true
		server.LastPing = datetime.ToISO(datetime.Now())
		u.Save()
	}()
}

func (u *User) Send(action string, data interface{}) {
    if u.w == nil || u.r == nil { return }

    flusher, ok := u.w.(http.Flusher)
    if !ok {
        u.logger.Error().Msg("Streaming unsupported by ResponseWriter")
        http.Error(u.w, "streaming unsupported", http.StatusInternalServerError)
        return
    }

    packet, err := NewPacket(version, "api", u.Token, action, data)
    if err != nil {
        u.logger.Error().Err(err).Msg("Failed to create packet")
        return
    }

    str := fmt.Sprintf("data: %s\n\n", packet.ToJSON())
    if str == "" {
        u.logger.Error().Msg("Failed to format packet data")
        return
    }

    _, err = u.w.Write([]byte(str))
    if err != nil {
        u.logger.Error().Err(err).Msg("Failed to write data to ResponseWriter")
        return
    }

    flusher.Flush()
}

// ------------------------------------------------------------
// : Controls
// ------------------------------------------------------------
func (u *User) Reload() {
	u.Send("reload", nil)
}

func (u *User) Start() {

	if u.flag_crawling.Load() {
		return
	}

	u.flag_crawling.Store(true)
	defer u.flag_crawling.Store(false)


    server := u.State.Server

    // Collect tasks that need processing (stale tasks)
    var todo = []*Task{}
    for _, task := range *server.Tasks {
        if task.IsStale() {
            todo = append(todo, task)
        }
    }

    // If there are no tasks, return early
    if len(todo) == 0 {
        return
    }

    u.Send("crawler.start", nil)
    u.logger.Info().Str("token", u.Token).Msg("Starting")

    server.StartedAt = datetime.ToISO(datetime.Now())
    u.Save()

    // Wait for the crawler to be ready
    if !u.WaitForState("ready", 5*time.Second) {
        return
    }

    // Group tasks in batches of 3 and process them
    for i := 0; i < len(todo); i += 5 {
        end := i + 5
        if end > len(todo) {
            end = len(todo)
        }

        batch := todo[i:end] // Take up to 3 tasks
        u.logger.Info().Str("token", u.Token).Int("batch_size", len(batch)).Msg("Processing batch of tasks")

        // Mark tasks as started
        for _, task := range batch {
            u.logger.Info().Str("token", u.Token).Str("task", task.Keyword).Str("status", "starting").Msg("Task")
            task.StartedAt = datetime.ToISO(datetime.Now())
        }

        // Send batch of tasks to crawler
        u.Send("crawler.scrape", batch)

        // Wait for the scraper to transition to "scraping"
        if !u.WaitForState("scraping", 5*time.Second) {
            continue // Skip to next batch if timeout occurs
        }

        // Wait for the scraper to complete and transition back to "ready"
        if !u.WaitForState("ready", 20*time.Second) {
            continue // Skip to next batch if timeout occurs
        }

        // Mark tasks as completed
        for _, task := range batch {
            task.CompletedAt = datetime.ToISO(datetime.Now())
            u.logger.Info().Str("token", u.Token).Str("task", task.Keyword).Str("status", "completed").Msg("Task")
        }

        u.Save()
    }

    // Mark the entire process as completed
    server.CompletedAt = datetime.ToISO(datetime.Now())
    u.Save()

    u.Send("crawler.complete", nil)
    u.logger.Info().Str("token", u.Token).Msg("All tasks completed")

    // Wait for the crawler to return to idle state
    if !u.WaitForState("idle", 3*time.Second) {
        return
    }
}

func (u *User) Reset() {
	u.logger.Info().Str("token", u.Token).Msg("Reset")

	if u.flag_resetting.Load() == true {
		return
	}

	u.flag_resetting.Store(true)
	defer u.flag_resetting.Store(false)

	server := u.State.Server
	server.CompletedAt = ""

	for _, task := range *server.Tasks {
		task.Reset()
	}

	u.Save()
}

// ------------------------------------------------------------
// : Setters
// ------------------------------------------------------------
func (u *User) SetVersion(version string) {
	u.version = version
}

// ------------------------------------------------------------
// : Getters
// ------------------------------------------------------------
func (u *User) IsOnline() bool {
	return u.State.Server.Online
}

func (u *User) ValidVersion() bool {
	if u.State.Client.Extension.Version == "" { return false }
	return u.State.Client.Extension.Version >= version
}


func (u *User) RequiersScraping() bool {
	var server   = u.State.Server
	var requires = false

	if server.TaskHash == "" { return true }

	for _, task := range *server.Tasks {
		if task.IsStale() {
			requires = true
			break
		}
	}

	return requires
}

func (u *User) CrawlingFlag() *atomic.Bool {
	return &u.flag_crawling
}
// ------------------------------------------------------------
// : Misc
// ------------------------------------------------------------
func (u *User) WaitForState(target string, timeout time.Duration) bool {
	var c = u.State.Client.Crawler

	c.mutex.Lock()
	defer c.mutex.Unlock()

	deadline := time.Now().Add(timeout)
	for c.State != target {
		remaining := deadline.Sub(time.Now())
		if remaining <= 0 {
			return false // Timeout
		}
		c.cond.Wait()
	}

	return true
}

func (u *User) Save() {
	u.State.Client.Crawler.cond.Broadcast()

	if u.flag_saving.Load() { 
		return 
	}

	u.flag_saving.Store(true)
	defer u.flag_saving.Store(false)

	ch := make(chan error)
	event.Emit("user.updated", u, ch)
	
	err := <- ch

	if err != nil {
		u.logger.Error().Err(err).Msg("")
	}
}

// ------------------------------------------------------------
// : Statics
// ------------------------------------------------------------
func newLogger(u *User) *zerolog.Logger {
	var logger  zerolog.Logger
	var multi   zerolog.LevelWriter
	writer_console := &zerolog.ConsoleWriter{}
 	writer_file    := &lumberjack.Logger{}

	location, err := time.LoadLocation("Europe/Amsterdam")
	if err != nil {
		panic(err)
	}

	writer_console.Out             = os.Stdout
	writer_console.TimeLocation    = location
	writer_console.FormatTimestamp = func(i interface{}) string {
		return time.Now().Format(time.RFC3339)
	}
	writer_console.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s=", i)
	}
	writer_console.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%v", i)
	}
	writer_console.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	writer_file.Filename   = "./logs/api.log"
	writer_file.MaxSize    = 100
	writer_file.MaxBackups = 100
	writer_file.MaxAge     = 30
	writer_file.Compress   = true

	multi  = zerolog.MultiLevelWriter(writer_console, writer_file)
	logger = zerolog.New(multi).With().
		Caller().
		Str("user"   , u.Token).
		Str("version", u.version).
		Timestamp().
		Logger()

	return &logger
}
