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
