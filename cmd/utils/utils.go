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

type Comparable interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~float32 | ~float64
}

func Min[T Comparable](a, b T) T {
	if a > b {
		return b
	} else {
		return a
	}
}

func Max[T Comparable](a, b T) T {
	if a < b {
		return b
	} else {
		return a
	}
}

func SearchList[T comparable](arr []T, target T) bool {
	for _, val := range arr {
		if val == target {
			return true
		}
	}
	return false
}
