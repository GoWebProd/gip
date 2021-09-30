package stack

import (
	"sync"
	"testing"
)

func TestStack(t *testing.T) {
	var s Stack[int]

	if s.Pop() != nil {
		t.Fatal("pop empty stack returns non-nil")
	}

	i1 := 1
	i2 := 2
	i3 := 3

	s.Push(&i1)
	s.Push(&i2)

	v := s.Pop()
	if v == nil || *v != i2 {
		t.Fatal("pop returns bad value")
	}

	s.Push(&i3)

	v = s.Pop()
	if v == nil || *v != i3 {
		t.Fatal("pop returns bad value")
	}

	v = s.Pop()
	if v == nil || *v != i1 {
		t.Fatal("pop returns bad value")
	}

	if s.Pop() != nil {
		t.Fatal("pop empty stack returns non-nil")
	}
}

func BenchmarkLockFree(b *testing.B) {
	var s Stack[int]

	for i := 0; i < b.N; i++ {
		s.Push(&i)
		s.Pop()
	}
}

func BenchmarkMutex(b *testing.B) {
	s := newMutexStack[int]()

	for i := 0; i < b.N; i++ {
		s.Push(&i)
		s.Pop()
	}
}

type mutexStack[T any] struct {
	v  []*T
	mu sync.Mutex
}

func newMutexStack[T any]() *mutexStack[T] {
	return &mutexStack[T]{v: make([]*T, 0)}
}

func (s *mutexStack[T]) Push(v *T) {
	s.mu.Lock()
	s.v = append(s.v, v)
	s.mu.Unlock()
}

func (s *mutexStack[T]) Pop() *T {
	s.mu.Lock()
	v := s.v[len(s.v)-1]
	s.v = s.v[:len(s.v)-1]
	s.mu.Unlock()
	return v
}
