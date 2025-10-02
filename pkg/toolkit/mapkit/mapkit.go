package mapkit


func GetString(data map[string]any, key string) string {
	if val, ok := data[key]; ok && val != nil {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

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

func GetBytes(data map[string]any, key string) []byte {
	if val, ok := data[key]; ok && val != nil {
		if bytes, ok := val.([]byte); ok {
			return bytes
		}
	}
	return nil
}

func GetBool(data map[string]any, key string) bool {
	if val, ok := data[key]; ok && val != nil {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false
}
