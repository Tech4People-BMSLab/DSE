package controller

import (
	"net/http"
	"os"

	"bms.dse/src/services/crawler"
	"bms.dse/src/utils/httputil"
	"bms.dse/src/utils/logutil"
	"github.com/go-chi/chi/v5"
)

// ------------------------------------------------------------
// : Aliases
// ------------------------------------------------------------
type JSON = map[string]any
type Writer  = http.ResponseWriter
type Request = http.Request
type User    = crawler.User

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var logger = logutil.NewLogger("API")
// ------------------------------------------------------------
// : Handler
// ------------------------------------------------------------
func HandleStartUser(w Writer, r *Request) {
	// http://localhost:5000/api/users/bd7f48026f63/start
	envToken   := os.Getenv("API_TOKEN")
	queryToken := r.URL.Query().Get("token")

	// Check if query token is valid
	if queryToken != envToken {
		logger.Warn().Msg("Unauthorized")
		httputil.WriteJSON(w, http.StatusUnauthorized, JSON{"error": "Unauthorized"})
		return
	}
	
	// Check if user is not empty
	userToken := chi.URLParam(r, "token")
	if userToken == "" {
		logger.Warn().Msg("Missing token")
		httputil.WriteJSON(w, http.StatusBadRequest, JSON{"error": "Missing token"})
		return
	}

	
	// Check if user exists in the database
	user  := &User{}
	users := crawler.GetUsers()
	for _, u := range users {
		if u.Token == userToken {
			user = u
			break
		}
	}
	if user == nil {
		logger.Warn().Str("token", userToken).Msg("User not found")
		httputil.WriteJSON(w, http.StatusNotFound, JSON{"error": "User not found"})
		return
	}


	// Start search for user
	logger.Info().Msg("Request to Start Search for User")
	go user.Start()

	httputil.WriteJSON(w, http.StatusOK, JSON{"status": "ok"})
}

func HandleStopUser(w Writer, r *Request) {
	// http://localhost:5000/api/users/bd7f48026f63/stop
	envToken  := os.Getenv("API_TOKEN")
	queryToken := r.URL.Query().Get("token")

	// Check if query token is valid
	if queryToken != envToken {
		logger.Warn().Msg("Unauthorized")
		httputil.WriteJSON(w, http.StatusUnauthorized, JSON{"error": "Unauthorized"})
		return
	}

	// Check if user is not empty
	userToken := chi.URLParam(r, "token")
	if userToken == "" {
		logger.Warn().Msg("Missing token")
		httputil.WriteJSON(w, http.StatusBadRequest, JSON{"error": "Missing token"})
		return
	}

	// Check if user exists in the database
	user  := &User{}
	users := crawler.GetUsers()
	for _, u := range users {
		if u.Token == userToken {
			user = u
			break
		}
	}
	if user == nil {
		logger.Warn().Str("token", userToken).Msg("User not found")
		httputil.WriteJSON(w, http.StatusNotFound, JSON{"error": "User not found"})
		return
	}

	// Stop search for user
	logger.Info().Msg("Request to Stop Search for User")
	go user.Stop()

	httputil.WriteJSON(w, http.StatusOK, JSON{"status": "ok"})
}
