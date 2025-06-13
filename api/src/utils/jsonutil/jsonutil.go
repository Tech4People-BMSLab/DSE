package jsonutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// JSON represents a JSON object with thread-safe operations.
type JSON struct {
    data interface{}
    mu   sync.RWMutex
}

// NewJSON creates a new instance of JSON.
func NewJSON() *JSON {
    return &JSON{data: make(map[string]interface{})}
}

// FromJSON initializes the JSON object from a JSON byte slice.
func (j *JSON) FromJSON(data []byte) error {
    var temp interface{}
    if err := json.Unmarshal(data, &temp); err != nil {
        return err
    }
    j.mu.Lock()
    defer j.mu.Unlock()
    j.data = temp
    return nil
}

// ToJSON serializes the JSON object to a JSON byte slice.
func (j *JSON) ToJSON() ([]byte, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return json.Marshal(j.data)
}

// Value returns the underlying data as map[string]interface{}.
func (j *JSON) Value() (map[string]interface{}, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    if m, ok := j.data.(map[string]interface{}); ok {
        return m, nil
    }
    return nil, errors.New("data is not a JSON object")
}

// Set sets a value in the JSON object at the specified path.
func (j *JSON) Set(path string, value interface{}) error {
    j.mu.Lock()
    defer j.mu.Unlock()
    return set(j.data, path, value)
}

// Get retrieves a value from the JSON object at the specified path.
func (j *JSON) Get(path string) (interface{}, bool) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return get(j.data, path)
}

// Delete removes a value from the JSON object at the specified path.
func (j *JSON) Delete(path string) error {
    j.mu.Lock()
    defer j.mu.Unlock()
    return deleteKey(j.data, path)
}

// Exists checks if a key exists at the specified path.
func (j *JSON) Exists(path string) bool {
    _, ok := j.Get(path)
    return ok
}

// Keys returns the keys at the specified path.
func (j *JSON) Keys(path string) ([]string, error) {
    j.mu.RLock()
    defer j.mu.RUnlock()
    node, ok := get(j.data, path)
    if !ok {
        return nil, errors.New("path does not exist")
    }
    if m, ok := node.(map[string]interface{}); ok {
        keys := make([]string, 0, len(m))
        for k := range m {
            keys = append(keys, k)
        }
        return keys, nil
    }
    return nil, errors.New("not an object at path")
}

// Verify checks if the JSON object matches the provided schema.
func (j *JSON) Verify(schema interface{}) error {
    j.mu.RLock()
    defer j.mu.RUnlock()
    return verifySchema(j.data, schema)
}

//
// Helper functions
//

func set(data interface{}, path string, value interface{}) error {
    parts := strings.Split(path, ".")
    last := len(parts) - 1
    current := data

    for i, part := range parts {
        if i == last {
            switch curr := current.(type) {
            case map[string]interface{}:
                curr[part] = value
                return nil
            default:
                return fmt.Errorf("invalid path: %s", strings.Join(parts[:i+1], "."))
            }
        } else {
            switch curr := current.(type) {
            case map[string]interface{}:
                if _, ok := curr[part]; !ok {
                    curr[part] = make(map[string]interface{})
                }
                current = curr[part]
            default:
                return fmt.Errorf("invalid path: %s", strings.Join(parts[:i+1], "."))
            }
        }
    }
    return nil
}

func get(data interface{}, path string) (interface{}, bool) {
    parts := strings.Split(path, ".")
    current := data

    for _, part := range parts {
        switch curr := current.(type) {
        case map[string]interface{}:
            if val, ok := curr[part]; ok {
                current = val
            } else {
                return nil, false
            }
        default:
            return nil, false
        }
    }
    return current, true
}

func deleteKey(data interface{}, path string) error {
    parts := strings.Split(path, ".")
    last := len(parts) - 1
    current := data

    for i, part := range parts {
        switch curr := current.(type) {
        case map[string]interface{}:
            if i == last {
                if _, ok := curr[part]; ok {
                    delete(curr, part)
                    return nil
                }
                return fmt.Errorf("key does not exist: %s", path)
            } else {
                if val, ok := curr[part]; ok {
                    current = val
                } else {
                    return fmt.Errorf("invalid path: %s", strings.Join(parts[:i+1], "."))
                }
            }
        default:
            return fmt.Errorf("invalid path: %s", strings.Join(parts[:i+1], "."))
        }
    }
    return nil
}

func verifySchema(data, schema interface{}) error {
    switch s := schema.(type) {
    case map[string]interface{}:
        d, ok := data.(map[string]interface{})
        if !ok {
            return errors.New("data is not an object")
        }
        for key, val := range s {
            if dv, ok := d[key]; ok {
                if err := verifySchema(dv, val); err != nil {
                    return fmt.Errorf("key '%s': %v", key, err)
                }
            } else {
                return fmt.Errorf("missing key: %s", key)
            }
        }
    case []interface{}:
        d, ok := data.([]interface{})
        if !ok {
            return errors.New("data is not an array")
        }
        if len(s) == 0 {
            return nil // Empty schema array matches any array
        }
        for i, item := range d {
            if err := verifySchema(item, s[0]); err != nil {
                return fmt.Errorf("index %d: %v", i, err)
            }
        }
    case string:
        if s == "any" {
            return nil
        }
        if reflect.TypeOf(data).String() != s {
            return fmt.Errorf("expected type %s, got %T", s, data)
        }
    default:
        if !reflect.DeepEqual(data, schema) {
            return fmt.Errorf("expected value %v, got %v", schema, data)
        }
    }
    return nil
}
