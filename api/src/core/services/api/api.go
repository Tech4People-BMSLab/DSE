package api

import (
	"compress/gzip"
	"dse/src/core/global"
	"dse/src/core/log"
	"dse/src/core/models"
	"dse/src/core/services/api/controller"
	"dse/src/core/services/api/download"
	"dse/src/core/services/api/metrics"
	"dse/src/core/services/api/ws"
	"dse/src/core/services/db"
	"dse/src/core/services/extractor"
	"dse/src/utils/datetime"
	"dse/src/utils/event"
	"dse/src/utils/gatekeeper"
	"dse/src/utils/hashmap"
	"dse/src/utils/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/dromara/carbon/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tidwall/gjson"
	"github.com/vmihailenco/msgpack/v5"
)

// ------------------------------------------------------------
// : Aliases
// ------------------------------------------------------------
type Packet = models.Packet
type State  = models.State
type User   = models.User
type Task   = models.Task
type Form   = models.Form
type Search = models.Search

type ClientState = models.ClientState
type ServerState = models.ServerState
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	router  = chi.NewRouter()
	
	addr string = fmt.Sprintf("%s:%s", "localhost", "5000")
	host string = "localhost"
	port string = "5000"

	started carbon.Carbon = carbon.Now(carbon.UTC)
	version string        = global.VERSION
	mutex   sync.Mutex    = sync.Mutex{}

	gk = gatekeeper.NewGateKeeper(true)

	// Errors
	ErrMissingToken         = fmt.Errorf("missing token")
	ErrInvalidToken         = fmt.Errorf("invalid token")

	ErrStreamingUnsupported = fmt.Errorf("streaming unsupported")
	ErrInvalidContentType   = fmt.Errorf("invalid content type")
	ErrFileTooLarge         = fmt.Errorf("file too large")
	ErrUnauthorized         = fmt.Errorf("unauthorized")
	ErrInvalidDays          = fmt.Errorf("invalid days parameter")
	ErrInternal             = fmt.Errorf("internal server error")
)
// ------------------------------------------------------------
// : Helpers
// ------------------------------------------------------------
func validateToken(w http.ResponseWriter, r *http.Request) (string, error) {
	token := r.URL.Query().Get("token")

	if token == ""      { 
		http.Error(w, ErrMissingToken.Error(), http.StatusBadRequest)
		return "", ErrMissingToken
	}

	if len(token) != 12 {
		http.Error(w, ErrInvalidToken.Error(), http.StatusBadRequest)
		return "", ErrInvalidToken
	}

	return token, nil
}

func parseVersion(w http.ResponseWriter, r *http.Request) (string, error) {
	version := r.URL.Query().Get("version")
	
	if version == "" {
		http.Error(w, "Missing version", http.StatusBadRequest)
		return "", fmt.Errorf("missing version")
	}

	return version, nil
}

func getFlusher(w http.ResponseWriter) (http.Flusher, error) {
	var flusher, ok = w.(http.Flusher)
	if !ok {
		http.Error(w, ErrStreamingUnsupported.Error(), http.StatusInternalServerError)
		return nil, ErrStreamingUnsupported
	}
	return flusher, nil
}

func inflate(data io.Reader) ([]byte, error) {
	reader, err := gzip.NewReader(data)
	if err != nil { return nil, err }
	defer reader.Close()

	inflated, err := io.ReadAll(reader)
	if err != nil { return nil, err }

	return inflated, nil
}

func unpack(data []byte) (Packet, error) {
	var packet Packet
	err := msgpack.Unmarshal(data, &packet)
	return packet, err
}

// ------------------------------------------------------------
// : Methods
// ------------------------------------------------------------
func Wait() {
	gk.Wait()
}

// ------------------------------------------------------------
// : Handlers
// ------------------------------------------------------------
func GetRoot(w http.ResponseWriter, r *http.Request) {
	defer recover()

	response := hashmap.NewHashMap[string, any]()
	response.Set("timestamp", datetime.ToISO(datetime.Now()))
	response.Set("version"  , version)
	response.Set("started"  , started.ToIso8601String())
	response.Set("uptime"   , started.DiffForHumans(carbon.Now(carbon.UTC)))

	b, err := response.ToJSON()
	if err != nil {
		log.Error().Err(err).Msg("Failed to encode response")
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func GetReload(w http.ResponseWriter, r *http.Request) {
	defer recover()

	users, err := db.GetUsers()
	if err != nil {
		log.Error().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users.Each(func(index int, user *User) bool {
		if user.IsOnline() {
			user.Send("reload", nil)
		}
		return true
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func getReset(w http.ResponseWriter, r *http.Request) {
	defer recover()

	token      := r.URL.Query().Get("token")
	users, err := db.GetUsers()
	if err != nil {
		log.Error().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if token == "" {
		users.Each(func(index int, user *User) bool {
			user.Reset()
			return true
		})
	} else {
		users.Each(func(index int, user *User) bool {
			if user.Token == token {
				user.Reset()
			}
			return true
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
// ------------------------------------------------------------
// : Client Handlers
// ------------------------------------------------------------
func GetEvent(w http.ResponseWriter, r *http.Request) {
	qtoken, err := validateToken(w, r)
	if  err != nil { return }

    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

	flusher, err := getFlusher(w)
	if err != nil { 
		log.Error().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	flusher.Flush()

	// Check if user exists
	exists, err := db.HasUser(qtoken)
	if err != nil {
		log.Error().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create user if not exists
	if !exists {	
		event := event.On(fmt.Sprintf("user.%s.created", qtoken))
		<- event // Wait for user creation (in the event handler)
	}

	// Get user
	user, err := db.GetUser(qtoken)
	if err != nil {
		log.Error().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.SetSSE(w, r)

	event.Emit(event.UserConnected, user)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	defer event.Emit(event.UserDisconnected, user)

	for {
		select {
			case <-r.Context().Done(): // Wait for disconnect
				return
			case <-time.After(5 * time.Second): // Heartbeat
				if !user.IsOnline() {
					return
				}
				user.Send("heartbeat", nil)
		}
	}
}


func PostEvent(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/octet-stream" {
		http.Error(w, ErrInvalidContentType.Error(), http.StatusUnsupportedMediaType)
		return
	}

	inflated, err := inflate(r.Body)
	if err != nil { 
		log.Error().Err(err).Msg("Failed inflation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	packet, err := unpack(inflated)
	if err != nil { 
		log.Error().Err(err).Msg("Failed unpack")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if packet.Version < version {
		log.Error().Str("version", packet.Version).Str("token", packet.From).Msg("Unsupported version")
		http.Error(w, "Unsupported version", http.StatusBadRequest)
		return
	}

	exists, err := db.HasUser(packet.From)
	if err != nil {
		log.Error().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var user *User

	if !exists {
		user, err = db.CreateUser(packet.From)
		if err != nil {
			log.Error().Err(err).Msg("")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		event.Emit(fmt.Sprintf("user.%s.created", user.Token))
	}

	user, err = db.GetUser(packet.From)
	if err != nil {
		log.Error().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.ToBytes(packet.Data)
	if err != nil {
		log.Error().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.State.Server.Online   = true
	user.State.Server.LastPing = carbon.Now(carbon.UTC).ToIso8601String()
	user.Save()

	mutex.Lock()
	defer mutex.Unlock()

	switch packet.Action {
		case "update": {
			func() {
				defer func() {
					if r := recover(); r != nil {
						err, _ := r.(error)
						log.Error().Err(err).Msg("Recovered")
					}
				}()
				
				parsed := gjson.ParseBytes(b)
	
				user.SetVersion(packet.Version)
				
				user.State.Client.Background.State = parsed.Get("background.state").String()
	
				user.State.Client.Browser = map[string]interface{}{
					"name"       : parsed.Get("browser.name").String(),
					"os"         : parsed.Get("browser.os").String(),
					"os_version" : parsed.Get("browser.os_version").String(),
					"version"    : parsed.Get("browser.version").String(),
				}
	
				user.State.Client.Crawler.State       = parsed.Get("crawler.state").String()
				user.State.Client.Crawler.Window      = parsed.Get("crawler.window").String()
				user.State.Client.Crawler.StartedAt   = parsed.Get("crawler.started_at").String()
				user.State.Client.Crawler.CompletedAt = parsed.Get("crawler.completed_at").String()
	
				user.State.Client.Extension.State    = parsed.Get("extension.state").String()
				user.State.Client.Extension.Version  = parsed.Get("extension.version").String()
				user.State.Client.Extension.Language = parsed.Get("extension.language").String()
	
				user.State.Client.User.Type  = parsed.Get("user.type").String()
				user.State.Client.User.Token = parsed.Get("user.token").String()
				user.State.Client.User.Popup = parsed.Get("user.popup").Bool()
	
				user.State.Client.User.Form.Age = parsed.Get("user.form.age").String()
				user.State.Client.User.Form.Sex = parsed.Get("user.form.sex").String()
	
				user.State.Client.User.Form.Browser.Brave         = parsed.Get("user.form.browser.brave").Bool()
				user.State.Client.User.Form.Browser.Chrome        = parsed.Get("user.form.browser.chrome").Bool()
				user.State.Client.User.Form.Browser.Firefox       = parsed.Get("user.form.browser.edge").Bool()
				user.State.Client.User.Form.Browser.MicrosoftEdge = parsed.Get("user.form.browser.firefox").Bool()
				user.State.Client.User.Form.Browser.Opera         = parsed.Get("user.form.browser.opera").Bool()
				user.State.Client.User.Form.Browser.Safari        = parsed.Get("user.form.browser.safari").Bool()
				user.State.Client.User.Form.Browser.Unselected	  = parsed.Get("user.form.browser.unselected").Bool()
	
				user.State.Client.User.Form.Social.TV	          = parsed.Get("user.form.social.tv").Bool()
				user.State.Client.User.Form.Social.DeKrant        = parsed.Get("user.form.social.de-krant").Bool()
				user.State.Client.User.Form.Social.Nieuwswebsites = parsed.Get("user.form.social.nieuwswebsites").Bool()
				user.State.Client.User.Form.Social.YouTube        = parsed.Get("user.form.social.youtube").Bool()
				user.State.Client.User.Form.Social.Facebook       = parsed.Get("user.form.social.facebook").Bool()
				user.State.Client.User.Form.Social.Instagram      = parsed.Get("user.form.social.instagram").Bool()
				user.State.Client.User.Form.Social.WhatsApp       = parsed.Get("user.form.social.whatsapp").Bool()
				user.State.Client.User.Form.Social.Linkedin       = parsed.Get("user.form.social.linkedin").Bool()
				user.State.Client.User.Form.Social.Twitter        = parsed.Get("user.form.social.twitter").Bool()
				user.State.Client.User.Form.Social.Telegram       = parsed.Get("user.form.social.telegram").Bool()
				user.State.Client.User.Form.Social.Reddit         = parsed.Get("user.form.social.reddit").Bool()
				user.State.Client.User.Form.Social.Radio          = parsed.Get("user.form.social.radio").Bool()
				user.State.Client.User.Form.Social.Anders         = parsed.Get("user.form.social.anders").Bool()
				user.State.Client.User.Form.Social.Unselected     = parsed.Get("user.form.social.unselected").Bool()
	
				user.State.Client.User.Form.Education  = parsed.Get("user.form.education").String()
				user.State.Client.User.Form.Employment = parsed.Get("user.form.employment").String()
				user.State.Client.User.Form.Income     = parsed.Get("user.form.income").String()
	
				user.State.Client.User.Form.Language.Duits      = parsed.Get("user.form.language.duits").Bool()
				user.State.Client.User.Form.Language.Engels     = parsed.Get("user.form.language.engels").Bool()
				user.State.Client.User.Form.Language.Frans      = parsed.Get("user.form.language.frans").Bool()
				user.State.Client.User.Form.Language.Italiaans  = parsed.Get("user.form.language.italiaans").Bool()
				user.State.Client.User.Form.Language.Nederlands = parsed.Get("user.form.language.nederlands").Bool()
				user.State.Client.User.Form.Language.Spaans     = parsed.Get("user.form.language.spaans").Bool()
				user.State.Client.User.Form.Language.Unselected = parsed.Get("user.form.language.unselected").Bool()
	
				user.State.Client.User.Form.Political      = parsed.Get("user.form.political").String()
				user.State.Client.User.Form.Postcode.Value = parsed.Get("user.form.postcode.value").String()
				user.State.Client.User.Form.Resident       = parsed.Get("user.form.resident").String()
	
				user.State.Client.User.Form.SearchEngine.Anders     = parsed.Get("user.form.search_engine.anders").Bool()
				user.State.Client.User.Form.SearchEngine.Bing       = parsed.Get("user.form.search_engine.bing").Bool()
				user.State.Client.User.Form.SearchEngine.Duckduckgo = parsed.Get("user.form.search_engine.duckduckgo").Bool()
				user.State.Client.User.Form.SearchEngine.Ecosia     = parsed.Get("user.form.search_engine.ecosia").Bool()
				user.State.Client.User.Form.SearchEngine.Google     = parsed.Get("user.form.search_engine.google").Bool()
				user.State.Client.User.Form.SearchEngine.Startpage  = parsed.Get("user.form.search_engine.startpage").Bool()
				user.State.Client.User.Form.SearchEngine.Yahoo      = parsed.Get("user.form.search_engine.yahoo").Bool()
				user.State.Client.User.Form.SearchEngine.Unselected = parsed.Get("user.form.search_engine.unselected").Bool()
				user.Save()
			}()
		}

		case "upload": { 
			go extractor.OnUpload(user, b)
		}
		case "reset" : {
			event.Emit("user.reset", user)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ------------------------------------------------------------
// : Init
// ------------------------------------------------------------
func Init() {
	if value, ok := os.LookupEnv("API_HOST"); ok { host = value }
	if value, ok := os.LookupEnv("API_PORT"); ok { port = value }
	addr = fmt.Sprintf("%s:%s", host, port)

	// Middlewares
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))
	
	// Routes
	router.Get("/ws", ws.HandleWS)
	
	router.Get("/api",              GetRoot)
	router.Get("/api/event",        GetEvent)
	router.Post("/api/event",       PostEvent)
	router.Get("/api/debug/reload", GetReload)

	router.Get("/api/download/logs",           download.GetLogs)
	router.Get("/api/download/users",          download.GetUsers)
	router.Get("/api/download/searches", 	   download.GetSearches)
	router.Get("/api/download/searches/full",  download.GetSearches)

	router.Get("/api/users/reset", controller.HandleReset)

	router.Get("/api/metrics",                metrics.HandleHealthCheck)
	router.Get("/api/metrics/users",          metrics.GetMetricUsers)
	router.Get("/api/metrics/searches",       metrics.GetMetricSearch)
	router.Get("/api/metrics/searches/size",  metrics.GetMetricsSearchSize)
	router.Get("/api/metrics/searches/total", metrics.GetMetricsSearchTotal)

	// Serve assets at /assets
	router.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/assets"))))
	router.Handle("/images/*", http.StripPrefix("/images/", http.FileServer(http.Dir("./public/images"))))
	router.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch (true) { 
			case r.URL.Path == "/"        : { http.Redirect(w, r, "/dse/consent", http.StatusTemporaryRedirect) }
			case r.URL.Path == "/consent" : { http.ServeFile(w, r, "./public/index.html") }
			case r.URL.Path == "/form"    : { http.ServeFile(w, r, "./public/index.html") }
			case r.URL.Path == "/complete": { http.ServeFile(w, r, "./public/index.html") }
			
			default: { http.FileServer(http.Dir("./public")).ServeHTTP(w, r) }
		}
	}))
	
	db.Wait()

	// Start server
	go func() {
		err := http.ListenAndServe(addr, router)
		if err != nil {
			log.Fatal().Err(err).Str("addr", addr).Msg("Failed to start server")
			return
		}
	}()
	log.Info().Str("addr", addr).Msg("Ready")
	gk.Unlock()
	
	// TODO: Remove
	// go func() {
	// 	http.Get("http://localhost:5000/api/download/searches/merge?token=dse2024&start=2024-11-01&end=2024-11-10")
	// }()

	// TODO: Remove
	// tool.Prepare()
}
