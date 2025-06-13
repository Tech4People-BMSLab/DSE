package object

import (
	"dse/src/utils/json"
	"fmt"
	"strconv"
	"strings"
)

// Get retrieves a value from a nested map[string]any or []any using a dot-separated key path.
// It returns the value and a boolean indicating whether the retrieval was successful.
func Get(data map[string]any, path string) (any, bool) {
    keys := strings.Split(path, ".")
    var current any = data

    for _, key := range keys {
        switch curr := current.(type) {
        case map[string]any:
            var ok bool
            current, ok = curr[key] // Use '=' to update the outer 'current' variable
            if !ok {
                return nil, false
            }
            // Continue to next key
        case []any:
            // Attempt to parse the key as an integer index
            index, err := strconv.Atoi(key)
            if err != nil || index < 0 || index >= len(curr) {
                return nil, false
            }
            current = curr[index]
        default:
            // Current level is neither map nor slice
            return nil, false
        }
    }

    return current, true
}

// GetString retrieves a string value from the nested path.
func GetString(data map[string]any, path string) (string, bool) {
    val, ok := Get(data, path)
    if !ok {
        return "", false
    }
    str, ok := val.(string)
    return str, ok
}

// GetInt retrieves an integer value from the nested path.
func GetInt(data map[string]any, path string) (int, bool) {
    val, ok := Get(data, path)
    if !ok {
        return 0, false
    }
    switch v := val.(type) {
    case int:
        return v, true
    case float64:
        return int(v), true
    case string:
        i, err := strconv.Atoi(v)
        if err != nil {
            return 0, false
        }
        return i, true
    default:
        return 0, false
    }
}

// GetFloat retrieves a float64 value from the nested path.
func GetFloat(data map[string]any, path string) (float64, bool) {
    val, ok := Get(data, path)
    if !ok {
        return 0, false
    }
    switch v := val.(type) {
    case float64:
        return v, true
    case int:
        return float64(v), true
    case string:
        f, err := strconv.ParseFloat(v, 64)
        if err != nil {
            return 0, false
        }
        return f, true
    default:
        return 0, false
    }
}

// GetBool retrieves a boolean value from the nested path.
func GetBool(data map[string]any, path string) (bool, bool) {
    val, ok := Get(data, path)
    if !ok {
        return false, false
    }
    switch v := val.(type) {
    case bool:
        return v, true
    case string:
        b, err := strconv.ParseBool(v)
        if err != nil {
            return false, false
        }
        return b, true
    case int:
        return v != 0, true
    case float64:
        return v != 0, true
    default:
        return false, false
    }
}

// GetSlice retrieves a slice value from the nested path.
func GetSlice(data map[string]any, path string) ([]any, bool) {
    val, ok := Get(data, path)
    if !ok {
        return nil, false
    }
    slice, ok := val.([]any)
    return slice, ok
}

// GetMap retrieves a map value from the nested path.
func GetMap(data map[string]any, path string) (map[string]any, bool) {
    val, ok := Get(data, path)
    if !ok {
        return nil, false
    }
    m, ok := val.(map[string]any)
    return m, ok
}

// Set assigns a value to a nested key path within the data structure.
// It creates maps or slices along the path if they do not exist.
func Set(data map[string]any, path string, value any) error {
    keys := strings.Split(path, ".")
    var current any = data

    for i, key := range keys {
        last := i == len(keys)-1
        switch curr := current.(type) {
        case map[string]any:
            if last {
                curr[key] = value
                return nil
            }
            next, exists := curr[key]
            if !exists {
                // Determine if the next key is an index or a map key
                if _, err := strconv.Atoi(keys[i+1]); err == nil {
                    next = []any{}
                } else {
                    next = map[string]any{}
                }
                curr[key] = next
            }
            current = curr[key]
        case []any:
            index, err := strconv.Atoi(key)
            if err != nil {
                return fmt.Errorf("invalid index '%s' at '%s'", key, strings.Join(keys[:i], "."))
            }
            if index < 0 {
                return fmt.Errorf("negative index '%d' at '%s'", index, strings.Join(keys[:i], "."))
            }
            // Expand the slice if necessary
            for len(curr) <= index {
                curr = append(curr, nil)
            }
            if last {
                curr[index] = value
                return nil
            }
            if curr[index] == nil {
                if _, err := strconv.Atoi(keys[i+1]); err == nil {
                    curr[index] = []any{}
                } else {
                    curr[index] = map[string]any{}
                }
            }
            current = curr[index]
        default:
            return fmt.Errorf("unexpected type at '%s'", strings.Join(keys[:i], "."))
        }
    }

    return nil
}

// Delete removes a value from the nested key path within the data structure.
func Delete(data map[string]any, path string) bool {
    keys := strings.Split(path, ".")
    var current any = data

    for i, key := range keys {
        last := i == len(keys)-1
        switch curr := current.(type) {
        case map[string]any:
            if last {
                if _, exists := curr[key]; exists {
                    delete(curr, key)
                    return true
                }
                return false
            }
            var ok bool
            current, ok = curr[key]
            if !ok {
                return false
            }
        case []any:
            index, err := strconv.Atoi(key)
            if err != nil || index < 0 || index >= len(curr) {
                return false
            }
            if last {
                curr[index] = nil
                return true
            }
            current = curr[index]
        default:
            return false
        }
    }

    return false
}


// Has checks whether a nested key path exists within the map.
// It returns true if the path leads to a valid value.
func Has(data map[string]any, path string) bool {
    _, ok := Get(data, path)
    return ok
}

// ListKeys lists all immediate keys at the specified path.
func ListKeys(data map[string]any, path string) ([]string, bool) {
    val, ok := Get(data, path)
    if !ok {
        return nil, false
    }

    switch v := val.(type) {
    case map[string]any:
        keys := make([]string, 0, len(v))
        for key := range v {
            keys = append(keys, key)
        }
        return keys, true
    case []any:
        keys := make([]string, len(v))
        for i := range v {
            keys[i] = strconv.Itoa(i)
        }
        return keys, true
    default:
        return nil, false
    }
}

// ToStruct converts a nested map[string]any or []any to a given struct.
// It returns an error if the conversion fails.
func ToStruct(data any, result any) error {
    b, err := json.ToBytes(data)
    if err != nil {
        return err
    }
    
    err = json.FromBytes(b, result)
    if err != nil {
        return err
    }

    return nil
}
