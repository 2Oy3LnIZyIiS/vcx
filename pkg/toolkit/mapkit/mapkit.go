// Package mapkit provides type-safe value extraction from map[string]any.
//
// Eliminates boilerplate type assertions when working with generic maps.
// Returns zero values if key doesn't exist or type doesn't match.
// Supports: string, int, int64, []byte, bool.
package mapkit


// GetString extracts a string value from map, returns empty string if not found.
func GetString(data map[string]any, key string) string {
	if val, ok := data[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

// GetInt extracts an int value from map, returns 0 if not found.
func GetInt(data map[string]any, key string) int {
	if val, ok := data[key]; ok && val != nil {
		if intVal, ok := val.(int); ok {
			return intVal
		}
		if int64Val, ok := val.(int64); ok {
			return int(int64Val)
		}
	}
	return 0
}

// GetBytes extracts a []byte value from map, returns nil if not found.
func GetBytes(data map[string]any, key string) []byte {
	if val, ok := data[key]; ok && val != nil {
		if bytes, ok := val.([]byte); ok {
			return bytes
		}
	}
	return nil
}

// GetBool extracts a bool value from map, returns false if not found.
func GetBool(data map[string]any, key string) bool {
	if val, ok := data[key]; ok && val != nil {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false
}
