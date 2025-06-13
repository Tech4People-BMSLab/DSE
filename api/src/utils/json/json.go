package json

import (
	"bytes"
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type JSON = map[string]any

// ToBytes converts an interface{} value to a JSON byte array.
// It uses jsoniter library with an indentation step of 4 spaces for formatting.
//
// Parameters:
//   - v: interface{} - The value to be converted to JSON bytes
//
// Returns:
//   - []byte - The JSON formatted byte array
//   - error - An error if JSON marshaling fails, nil otherwise
func ToBytes(v interface{}) ([]byte, error) {
	var json = jsoniter.Config{
		IndentionStep: 4,
	}.Froze()
	var data, err = json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// ToString converts an interface{} value to its JSON string representation.
// It first marshals the value into bytes using ToByte and then converts the bytes to a string.
//
// Parameters:
//   - v: interface{} - The value to be converted to a JSON string.
//
// Returns:
//   - string: The JSON string representation of the input value.
//   - error: An error if marshaling fails, nil otherwise.
func ToString(v interface{}) (string, error) {
	data, err := ToBytes(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FromBytes unmarshals JSON data from a byte slice into a provided interface.
// It uses jsoniter library configured to be compatible with the standard library.
//
// Parameters:
//   - data: byte slice containing JSON data to be unmarshaled
//   - v: pointer to the variable where unmarshaled data will be stored
//
// Returns:
//   - error: nil if successful, error otherwise if unmarshaling fails
func FromBytes(data []byte, v interface{}) (error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var err  = json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

// MustFromBytes unmarshals JSON data from a byte slice into a provided interface.
// If unmarshaling fails, the program will panic.
//
// Parameters:
//   - data: byte slice containing JSON data to be unmarshaled
//
// Returns:
//   - interface{}: The unmarshaled value
func MustFromBytes(data []byte) interface{} {
	var v interface{}
	if err := FromBytes(data, &v); err != nil {
		panic(err)
	}
	return v
}


// ParseBytes unmarshalls JSON into Result object.
// It uses the gjson library to parse JSON data and return a Result object.
//
// Parameters:
//   - data: byte slice containing JSON data to be parsed	
//
// Returns:
//   - gjson.Result: The parsed JSON data as a gjson.Result object
func ParseBytes(data []byte) gjson.Result {
	return gjson.ParseBytes(data)
}

// NewDecoder creates a new json.Decoder from a byte slice.
// It wraps the byte slice in a bytes.Reader and returns a new decoder
// that can be used to decode JSON data.
//
// Parameters:
//   - data: byte slice containing JSON data to be decoded
//
// Returns:
//   - *json.Decoder: a new JSON decoder instance
func NewDecoder(data []byte) *json.Decoder {
	return json.NewDecoder(bytes.NewReader(data))
}


// NewEncoder creates and returns a new JSON encoder with a bytes buffer.
// The encoder can be used to write JSON data to the buffer.
// Note: The returned encoder uses a new, empty bytes buffer.
func NewEncoder() *json.Encoder {
	return json.NewEncoder(&bytes.Buffer{})
}

// Set assigns a value to a specific path within a JSON string using dot notation.
// It returns the modified JSON string and any error that occurred during the operation.
//
// Parameters:
//   - json: The input JSON string to modify
//   - path: The path where to set the value, using dot notation (e.g., "user.name")
//   - value: The value to set at the specified path
//
// Returns:
//   - string: The modified JSON string
//   - error: An error if the operation fails, nil otherwise
func Set(json string, path string, value interface{}) (string, error) {
	return sjson.Set(json, path, value)
}

// SetBytes assigns a value to a specific path within a JSON byte slice using dot notation.
// It converts the byte slice to a string, sets the value, and converts back to bytes.
//
// Parameters:
//   - json: The input JSON byte slice to modify
//   - path: The path where to set the value, using dot notation (e.g., "user.name")
//   - value: The value to set at the specified path
//
// Returns:
//   - []byte: The modified JSON as bytes
//   - error: An error if the operation fails, nil otherwise
func SetBytes(json []byte, path string, value interface{}) ([]byte, error) {
	b, err := sjson.SetBytes(json, path, value)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Get retrieves a value at a specific path within a JSON string using dot notation.
// It uses the github.com/tidwall/gjson package to safely extract values.
//
// Parameters:
//   - json: The input JSON string to read from
//   - path: The path to get the value from, using dot notation (e.g., "user.name")
//
// Returns:
//   - interface{}: The value found at the specified path
//   - error: An error if the operation fails, nil otherwise
func Get(json string, path string) gjson.Result {
	return gjson.Get(json, path)
}

// GetBytes retrieves a value at a specific path within a JSON byte slice using dot notation.
// It converts the byte slice to a string and uses gjson to extract values.
//
// Parameters:
//   - json: The input JSON byte slice to read from
//   - path: The path to get the value from, using dot notation (e.g., "user.name")
//
// Returns:
//   - gjson.Result: The value found at the specified path
func GetBytes(json []byte, path string) gjson.Result {
	return gjson.GetBytes(json, path)
}

// Keys returns all top-level keys from a JSON string.
// It parses the JSON string and extracts the keys from the root object.
//
// Parameters:
//   - json: The JSON string to extract keys from
//
// Returns:
//   - []string: A slice containing all top-level keys
//   - error: An error if parsing fails, nil otherwise
func Keys(json string) ([]string, error) {
	var keys []string
	gjson.Parse(json).ForEach(func(key, value gjson.Result) bool {
		keys = append(keys, key.String())
		return true
	})
	return keys, nil
}

// KeysBytes returns all top-level keys from a JSON byte slice.
// It converts the byte slice to a string and extracts the keys from the root object.
//
// Parameters:
//   - json: The JSON byte slice to extract keys from
//
// Returns:
//   - []string: A slice containing all top-level keys
//   - error: An error if parsing fails, nil otherwise
func KeysBytes(json []byte) ([]string, error) {
	var keys []string
	gjson.ParseBytes(json).ForEach(func(key, value gjson.Result) bool {
		keys = append(keys, key.String())
		return true
	})
	return keys, nil
}

// Values returns all top-level values from a JSON string.
// It parses the JSON string and extracts the values from the root object.
//
// Parameters:
//   - json: The JSON string to extract values from
//
// Returns:
//   - []interface{}: A slice containing all top-level values
//   - error: An error if parsing fails, nil otherwise
func Values(json string) ([]interface{}, error) {
	var values []interface{}
	gjson.Parse(json).ForEach(func(key, value gjson.Result) bool {
		values = append(values, value.Value())
		return true
	})
	return values, nil
}

// ValuesBytes returns all top-level values from a JSON byte slice.
// It converts the byte slice to a string and extracts the values from the root object.
//
// Parameters:
//   - json: The JSON byte slice to extract values from
//
// Returns:
//   - []interface{}: A slice containing all top-level values
//   - error: An error if parsing fails, nil otherwise
func ValuesBytes(json []byte) ([]interface{}, error) {
	var values []interface{}
	gjson.ParseBytes(json).ForEach(func(key, value gjson.Result) bool {
		values = append(values, value.Value())
		return true
	})
	return values, nil
}

// Format takes an interface{} and returns a formatted JSON string with proper indentation.
// It first converts the input to JSON bytes and then formats them.
//
// Parameters:
//   - v: interface{} - The value to be formatted as JSON
//
// Returns:
//   - string - The formatted JSON string
//   - error - An error if formatting fails, nil otherwise
func Format(v interface{}) (string, error) {
	bytes, err := ToBytes(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// FormatBytes takes a JSON byte slice and returns a formatted JSON byte slice with proper indentation.
// It first unmarshals the input bytes into an interface{} and then formats it using ToByte.
//
// Parameters:
//   - data: []byte - The JSON bytes to be formatted
//
// Returns:
//   - []byte - The formatted JSON bytes
//   - error - An error if formatting fails, nil otherwise
func FormatBytes(data []byte) ([]byte, error) {
	var v interface{}
	if err := FromBytes(data, &v); err != nil {
		return nil, err
	}
	return ToBytes(v)
}
