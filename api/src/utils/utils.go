package utils

import "os"

// ------------------------------------------------------------
// : Functions
// ------------------------------------------------------------
func IsEmpty(v interface{}) bool {
	switch v.(type) {
		case string:
			return v == ""
		case int:
			return v == 0
		case float64:
			return v == 0.0
		case bool:
			return v == false

		case []any:
			return len(v.([]any)) == 0

		default:
			return v == nil

	}
}

func IsProd() bool {
	return os.Getenv("DEBUG") != "1"
}

func IsDebugMode() bool {
	return os.Getenv("DEBUG") == "1"
}
