package utils

// ToInterfaces converts a slice of any type to []interface{}.
func ToInterfaces[T any](values []T) []interface{} {
	result := make([]interface{}, len(values))
	for i, v := range values {
		result[i] = v
	}
	return result
}
