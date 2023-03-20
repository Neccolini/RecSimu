package set

import (
	"math/rand"
)

type Set[T comparable] struct {
	hashMap map[T]struct{}
}

func NewSet[T comparable](t T) *Set[T] {
	return &Set[T]{map[T]struct{}{t: {}}}
}

func (s *Set[T]) Add(t T) {
	s.hashMap[t] = struct{}{}
}

func (s *Set[T]) Remove(t T) {
	delete(s.hashMap, t)
}

func (s *Set[T]) Size() int {
	return len(s.hashMap)
}

func (s *Set[T]) RandomChoice() T {
	randomKey := func() T {
		keys := make([]T, 0, len(s.hashMap))
		for k := range s.hashMap {
			keys = append(keys, k)
		}
		return keys[rand.Intn(len(keys))]
	}()
	return randomKey
}

func (s *Set[T]) ToList() []T {
	list := make([]T, 0, s.Size())
	for k := range s.hashMap {
		list = append(list, k)
	}
	return list
}

func (s *Set[T]) Contains(t T) bool {
	_, ok := s.hashMap[t]
	return ok
}