package controller

import (
	"dse/src/core/models"
	"dse/src/core/services/db"
	"dse/src/utils"
	"net/http"
)

// ------------------------------------------------------------
// : Locals
// ------------------------------------------------------------
var (
	logger = utils.NewLogger()
)

func HandleReset(w http.ResponseWriter, r *http.Request) {
	defer recover()

	token      := r.URL.Query().Get("token")
	users, err := db.GetUsers()
	if err != nil {
		logger.Error().Err(err).Msg("")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if token == "" {
		users.Each(func(index int, user *models.User) bool {
			user.Reset()
			return true
		})
	} else {
		users.Each(func(index int, user *models.User) bool {
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
