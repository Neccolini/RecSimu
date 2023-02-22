package random

import (
	"math/rand"
	"time"
)

func RandomChoice[T any](array []T) T {
	return array[randGenerator(len(array))]
}

func randGenerator(len int) int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len)
	return r
}
