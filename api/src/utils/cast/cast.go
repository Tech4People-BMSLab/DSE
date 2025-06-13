package cast

import (
	"strconv"

	c "github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

// ToBytes converts a given value to a byte slice ([]byte).
// It supports the following types:
//   - []byte: Returns the value as-is.
//   - string: Converts the string to a byte slice.
//   - gjson.Result: Converts the raw JSON result to a byte slice.
//
// Parameters:
//   - v: The value to convert.
//
// Returns:
//   - []byte: The converted byte slice, or nil if the type is unsupported.
func ToBytes(v any) []byte {
	switch v := v.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	case gjson.Result:
		return []byte(v.Raw)
	default:
		return []byte(c.ToString(v))
	}
}

// ToBytesSlice converts a given value to a slice of byte slices ([][]byte).
// It supports the following types:
//   - [][]byte: Returns the value as-is.
//   - []string: Converts each string to a byte slice.
//   - []any: Converts each element to a byte slice if it is a string or []byte.
//
// Parameters:
//   - v: The value to convert.
//
// Returns:
//   - [][]byte: The converted slice of byte slices, or nil if the type is unsupported.
func ToBytesSlice(v any) [][]byte {
	switch v := v.(type) {
	case [][]byte:
		return v
	case []string:
		var slices [][]byte
		for _, str := range v {
			slices = append(slices, []byte(str))
		}
		return slices
	case []any:
		var slices [][]byte
		for _, item := range v {
			switch item := item.(type) {
			case []byte:
				slices = append(slices, item)
			case string:
				slices = append(slices, []byte(item))
			default:
				return nil
			}
		}
		return slices
	default:
		return nil
	}
}

// ToString converts a given value to a string.
// It uses the cast package for safe conversion.
//
// Parameters:
//   - v: The value to convert.
//
// Returns:
//   - string: The converted string.
func ToString(v any) string {
	switch v := v.(type) {
	case gjson.Result:
		return v.Raw
	default:
		return c.ToString(v)
	}
}

// ToInt converts a given value to an integer.
// It uses the cast package for safe conversion.
//
// Parameters:
//   - v: The value to convert.
//
// Returns:
//   - int: The converted integer.
func ToInt(v any) int {
	switch v := v.(type) {
		case string:
			i, _ := strconv.Atoi(v)
			return i
		default:
			return c.ToInt(v)
	}
}

// ToUint converts a given value to an unsigned integer.
// It uses the cast package for safe conversion.
//
// Parameters:
//   - v: The value to convert.
//
// Returns:
//   - uint: The converted unsigned integer.
func ToUint(v any) uint {
	return c.ToUint(v)
}

// ToFloat32 converts a given value to a 32-bit floating point number.
// It uses the cast package for safe conversion.
//
// Parameters:
//   - v: The value to convert.
//
// Returns:
//   - float32: The converted float32 value.
func ToFloat32(v any) float32 {
	return c.ToFloat32(v)
}

// ToFloat64 converts a given value to a 64-bit floating point number.
// It uses the cast package for safe conversion.
//
// Parameters:
//   - v: The value to convert.
//
// Returns:
//   - float64: The converted float64 value.
func ToFloat64(v any) float64 {
	return c.ToFloat64(v)
}

// ToBool converts a given value to a boolean.
// It uses the cast package for safe conversion.
//
// Parameters:
//   - v: The value to convert.
//
// Returns:
//   - bool: The converted boolean value.
func ToBool(v any) bool {
	return c.ToBool(v)
}

// ToMap converts a given value to a map[string]interface{}.
// It uses the cast package for safe conversion.
//
// Parameters:
//   - v: The value to convert.
//
// Returns:
//   - map[string]interface{}: The converted map.
func ToMap(v any) map[string]interface{} {
	return c.ToStringMap(v)
}

// ToSlice converts a given value to a slice of interface{}.
// It uses the cast package for safe conversion.
//
// Parameters:
//   - v: The value to convert.
//
// Returns:
//   - []interface{}: The converted slice of interface{}.
func ToSlice(v any) []interface{} {
	return c.ToSlice(v)
}
