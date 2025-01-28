package stack

import "sync"

type MutexStack[T any] struct {
	v  []T
	mu sync.Mutex
}

func NewMutexStack[T any]() MutexStack[T] {
	return MutexStack[T]{v: make([]T, 0, 8)}
}

func (s *MutexStack[T]) Push(v T) {
	s.mu.Lock()
	s.PushUnlocked(v)
	s.mu.Unlock()
}

func (s *MutexStack[T]) PushUnlocked(v T) {
	s.v = append(s.v, v)
}

func (s *MutexStack[T]) Pop() (T, bool) {
	s.mu.Lock()

	v, ok := s.PopUnlocked()

	s.mu.Unlock()

	return v, ok
}

func (s *MutexStack[T]) PopUnlocked() (T, bool) {
	var v T

	if len(s.v) == 0 {
		return v, false
	}

	v = s.v[len(s.v)-1]
	s.v = s.v[:len(s.v)-1]

	return v, true
}
