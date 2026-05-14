package drift

import "fmt"

// FlattenMap converts a nested map[string]interface{} into a flat map
// with dot-separated keys, making deep comparison straightforward.
//
// Example:
//
//	{"spec": {"replicas": 3}} -> {"spec.replicas": 3}
func FlattenMap(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	flattenRecursive("", input, result)
	return result
}

// UnflattenMap converts a flat map with dot-separated keys back into
// a nested map[string]interface{}. It is the inverse of FlattenMap.
//
// Example:
//
//	{"spec.replicas": 3} -> {"spec": {"replicas": 3}}
func UnflattenMap(input map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range input {
		setNested(result, key, value)
	}
	return result
}

// setNested sets a value in a nested map using a dot-separated key path.
func setNested(m map[string]interface{}, key string, value interface{}) {
	for i := 0; i < len(key); i++ {
		if key[i] == '.' {
			prefix := key[:i]
			rest := key[i+1:]
			if _, ok := m[prefix]; !ok {
				m[prefix] = make(map[string]interface{})
			}
			if nested, ok := m[prefix].(map[string]interface{}); ok {
				setNested(nested, rest, value)
			}
			return
		}
	}
	m[key] = value
}

func flattenRecursive(prefix string, current map[string]interface{}, result map[string]interface{}) {
	for key, value := range current {
		fullKey := key
		if prefix != "" {
			fullKey = fmt.Sprintf("%s.%s", prefix, key)
		}

		switch v := value.(type) {
		case map[string]interface{}:
			flattenRecursive(fullKey, v, result)
		default:
			result[fullKey] = value
		}
	}
}
