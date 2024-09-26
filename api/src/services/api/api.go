package api

import (
	"net/http"
	"os"

	"bms.dse/src/services/api/controller"
	"bms.dse/src/services/api/download"
	"bms.dse/src/services/api/livelog"
	"bms.dse/src/services/crawler"
	"bms.dse/src/utils/logutil"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// ------------------------------------------------------------
// : Aliases
// ------------------------------------------------------------
type JSON    = map[string]any
type Writer  = http.ResponseWriter
type Request = http.Request
// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var logger = logutil.NewLogger("API")
var token string

var this = struct {
	Router *chi.Mux
}{

}
// ------------------------------------------------------------
// : Helpers
// ------------------------------------------------------------
func GetBuild() string {
	var build []byte

	_, err := os.Stat("build.txt")
	if err != nil {
		return "unknown"
	}

	build, err = os.ReadFile("build.txt")
	if err != nil {
		return "unknown"
	}

	return string(build)
}

// ------------------------------------------------------------
// : Inits
// ------------------------------------------------------------
func InitMiddlewares() {
	this.Router.Use(middleware.RequestID)
	this.Router.Use(middleware.RealIP)
	this.Router.Use(middleware.Recoverer)
	this.Router.Use(middleware.Logger)
	this.Router.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-store")
			h.ServeHTTP(w, r)
		})
	})

	this.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge        : 300, 
	}))
}

func InitRoutes() {
	this.Router.Get("/api/logs", livelog.HandleLog)

	this.Router.Get("/api/download/users"           , download.HandleDownloadUsers)
	this.Router.Get("/api/download/searches"        , download.HandleDownloadSearches)
	this.Router.Get("/api/download/searches/bypass" , download.HandleDownloadSearchesBypass)

	this.Router.Post("/api/upload", crawler.HandleUpload)

	this.Router.Get("/api/users/{token}/start", controller.HandleStartUser)
	this.Router.Get("/api/users/{token}/stop" , controller.HandleStopUser)

	this.Router.Get("/api/system/shutdown", func(w http.ResponseWriter, r *http.Request) {
		envToken   := os.Getenv("API_TOKEN")
		queryToken := r.URL.Query().Get("token")

		if queryToken == envToken {
			logger.Info().Msg("Shutting down")
			os.Exit(0)
			return
		}
	})
}

func InitStatic() {
	// Serve assets at /assets
	fsAssets := http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/assets")))
	fsPublic := http.FileServer(http.Dir("./public"))

	this.Router.Handle("/assets/*", fsAssets)
	this.Router.Handle("/*", fsPublic)

	this.Router.Handle("/", http.HandlerFunc(func(w Writer, r *Request) {
		http.ServeFile(w, r, "./public/index.html")
	}))
	this.Router.Handle("/consent", http.HandlerFunc(func(w Writer, r *Request) {
		http.ServeFile(w, r, "./public/index.html")
	}))
	this.Router.Handle("/form", http.HandlerFunc(func(w Writer, r *Request) {
		http.ServeFile(w, r, "./public/index.html")
	}))
	this.Router.Handle("/complete", http.HandlerFunc(func(w Writer, r *Request) {
		http.ServeFile(w, r, "./public/index.html")
	}))

	// Logo
	this.Router.Handle("/logo.png", http.HandlerFunc(func(w Writer, r *Request) {
		http.ServeFile(w, r, "./public/logo.png")
	}))
}

// ------------------------------------------------------------
// : API
// ------------------------------------------------------------
func Init() {
	this.Router = chi.NewRouter()

	envHost := os.Getenv("API_HOST")
	envPort := os.Getenv("API_PORT")
 
	// TODO: Remove these later in prod
	if envHost == "" { envHost = "0.0.0.0" }
	if envPort == "" { envPort = "5000" }

	if envHost == "" { logger.Panic().Msg("API_HOST is not set") }
	if envPort == "" { logger.Panic().Msg("API_PORT is not set") }
	
	InitMiddlewares()
	InitRoutes()
	InitStatic()

	logger.Info().Str("host", envHost).Str("port", envPort).Msg("Starting API")

	err := http.ListenAndServe(envHost + ":" + envPort, this.Router)
	if err != nil {
		logger.Panic().Err(err).Msg("Failed to start API")
	}
}
