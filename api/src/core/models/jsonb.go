package models

import (
	"encoding/json"
	"fmt"
)

// ------------------------------------------------------------
// : JSONBMap
// ------------------------------------------------------------
type JSONBMap map[string]interface{}

func (j JSONBMap) Value() ([]byte, error) {
	return json.Marshal(j)
}

func (j *JSONBMap) Scan(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("type assertion to string failed")
	}
	return json.Unmarshal([]byte(str), j)
}
