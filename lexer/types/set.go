package types

import (
	"maps"
	"slices"
)

type Set[T comparable] struct {
	storage map[T]bool
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		storage: make(map[T]bool),
	}
}

func (s *Set[T]) Put(k T) {
	if _, exists := s.storage[k]; exists {
		return
	}
	s.storage[k] = true
}

func (s *Set[T]) Get(k T) bool {
	if _, exists := s.storage[k]; exists {
		return exists
	}
	return false
}

func (s *Set[T]) GetLastValue() T {
	var t T
	if len(s.storage) > 0 {
		keys := slices.Collect(maps.Keys(s.storage))
		return keys[len(keys)-1]
	}
	return t
}

func (s *Set[T]) HasValue() bool {
	return len(s.storage) > 0
}

func (s *Set[T]) Delete(k T) {
	delete(s.storage, k)
}
