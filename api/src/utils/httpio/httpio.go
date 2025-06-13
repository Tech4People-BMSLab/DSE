package httpio

import (
	"dse/src/utils/arrays"
	"dse/src/utils/datetime"
	"dse/src/utils/json"
	"io"
	"net/http"

	"github.com/cohesivestack/valgo"
	"github.com/tidwall/gjson"
)

func WriteJSON(w http.ResponseWriter, code int, data json.JSON) error {
	b, err := json.ToBytes(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(b))
	return nil
}

func WriteError(w http.ResponseWriter, code int, message string) error {
	return WriteJSON(w, code, json.JSON{
		"timestamp": datetime.Now(),
		"error"    : message,
	})
}

func WriteValidationError(w http.ResponseWriter, v *valgo.Validation) error {
	index    := 0
	response := make(json.JSON)

	for _, err := range v.Errors() {
		field   := err.Name()
		message := *arrays.Last[string](err.Messages())
		response[field] = message
		index += 1
	}

	return WriteJSON(w, http.StatusBadRequest, json.JSON{
		"timestamp": datetime.Now(),
		"error":    response,
	})
}

func StreamJSON(w http.ResponseWriter, data json.JSON) error {
	b, err := json.ToBytes(data)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(b))
	return nil
}

func WriteStatus(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func ReadJSON(req *http.Request) (*gjson.Result, error) {
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	p := json.ParseBytes(b)
	return &p, nil
}

