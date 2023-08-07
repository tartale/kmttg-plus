package model

func Filter[T any](slice []T, f func(T) bool) []T {

	var result []T
	for _, item := range slice {
		if f(item) {
			result = append(result, item)
		}
	}

	return result
}
