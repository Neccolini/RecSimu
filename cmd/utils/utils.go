package utils

func Remove[T comparable](slice []T, target T) []T {
	result := slice
	var zero T
	for i, v := range slice {
		if v == target {
			result[i] = result[len(result)-1]
			result[len(result)-1] = zero
			result = result[:len(result)-1]
		}
	}
	return result
}
