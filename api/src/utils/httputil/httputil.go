package httputil

import (
	"encoding/json"
	"net/http"
)

// ------------------------------------------------------------
// : Aliases
// ------------------------------------------------------------
type JSON = map[string]any
// ------------------------------------------------------------
// : Helpers
// ------------------------------------------------------------
func WriteJSON(w http.ResponseWriter, code int, data JSON) error {
	response, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(response))
	return nil
}
