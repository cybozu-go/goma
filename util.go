// Utility functions for plugins.

package goma

const (
	EPSILON = 0.00000001
)

// FloatEquals compares two floats allowing error within EPSILON.
func FloatEquals(a, b float64) bool {
	if a == b {
		return true
	}
	return (a-b) < EPSILON && (b-a) < EPSILON
}

// GetBool extracts a boolean from TOML decoded map.
// If m[key] does not exist or is not a bool, non-nil error is returned.
func GetBool(key string, m map[string]interface{}) (bool, error) {
	v, ok := m[key]
	if !ok {
		return false, ErrNoKey
	}
	b, ok := v.(bool)
	if !ok {
		return false, ErrInvalidType
	}
	return b, nil
}

// GeInt extracts an integer from TOML decoded map.
// If m[key] does not exist or is not an integer, non-nil error is returned.
func GetInt(key string, m map[string]interface{}) (int, error) {
	v, ok := m[key]
	if !ok {
		return 0, ErrNoKey
	}
	i, ok := v.(int)
	if !ok {
		return 0, ErrInvalidType
	}
	return i, nil
}

// GetFloat extracts a float from TOML decoded map.
// If m[key] does not exist or is not a float, non-nil error is returned.
func GetFloat(key string, m map[string]interface{}) (float64, error) {
	v, ok := m[key]
	if !ok {
		return 0, ErrNoKey
	}
	f, ok := v.(float64)
	if !ok {
		return 0, ErrInvalidType
	}
	return f, nil
}

// GetString extracts a string from TOML decoded map.
// If m[key] does not exist or is not a string, non-nil error is returned.
func GetString(key string, m map[string]interface{}) (string, error) {
	v, ok := m[key]
	if !ok {
		return "", ErrNoKey
	}
	s, ok := v.(string)
	if !ok {
		return "", ErrInvalidType
	}
	return s, nil
}

// GetStringList constructs a string list from TOML decoded map.
// If m[key] does not exist or is not a string list, non-nil error is returned.
func GetStringList(key string, m map[string]interface{}) ([]string, error) {
	v, ok := m[key]
	if !ok {
		return nil, ErrNoKey
	}

	if sl, ok := v.([]string); ok {
		return sl, nil
	}

	l, ok := v.([]interface{})
	if !ok {
		return nil, ErrInvalidType
	}
	ret := make([]string, 0, len(l))
	for _, t := range l {
		s, ok := t.(string)
		if !ok {
			return nil, ErrInvalidType
		}
		ret = append(ret, s)
	}
	return ret, nil
}

// GetStringMap constructs a map[string]string from TOML decoded map.
// If m[key] does not exist or is not a string map, non-nil error is returned.
func GetStringMap(key string, m map[string]interface{}) (map[string]string, error) {
	v, ok := m[key]
	if !ok {
		return nil, ErrNoKey
	}

	if sm, ok := v.(map[string]string); ok {
		return sm, nil
	}

	m2, ok := v.(map[string]interface{})
	if !ok {
		return nil, ErrInvalidType
	}
	ret := make(map[string]string)
	for k, v2 := range m2 {
		s, ok := v2.(string)
		if !ok {
			return nil, ErrInvalidType
		}
		ret[k] = s
	}
	return ret, nil
}
