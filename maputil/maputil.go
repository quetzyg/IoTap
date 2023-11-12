package maputil

import "strings"

// KeyExists checks if a given key is present in the provided map or within any nested maps it might contain.
// The function takes two parameters: a map of type map[string]any and a key name string that supports dot notation for nested keys.
// For example, calling KeyExists(m, "a.b") will check if 'm' contains a map mapped to 'a', which in turn contains a key 'b'.
// KeyExists returns true if the key (or nested key) is found in the map and false otherwise.
// Note: This function only checks for the key's existence and does not consider its corresponding value.
func KeyExists(payload map[string]any, key string) bool {
	levels := strings.Split(key, ".")

	for _, level := range levels {
		val, exists := payload[level]
		if !exists {
			return false
		}

		if nested, ok := val.(map[string]any); ok {
			payload = nested
		}
	}

	return true
}
