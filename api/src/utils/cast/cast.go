package cast

import (
	c "github.com/spf13/cast"
)

// ------------------------------------------------------------
// : Safe Functions
// ------------------------------------------------------------
func ToBytes(v any) []byte {
	switch v := v.(type) {
	case []byte:
		return v
	case string:
		return []byte(v)
	default:
		return nil
	}
}

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

func ToString(v any) string {
	return c.ToString(v)
}

func ToInt(v any) int {
	return c.ToInt(v)
}

func ToUint(v any) uint {
	return c.ToUint(v)
}

func ToFloat32(v any) float32 {
	return c.ToFloat32(v)
}

func ToFloat64(v any) float64 {
	return c.ToFloat64(v)
}

func ToBool(v any) bool {
	return c.ToBool(v)
}

func ToMap(v any) map[string]interface{} {
	return c.ToStringMap(v)
}

func ToSlice(v any) []interface{} {
	return c.ToSlice(v)
}
