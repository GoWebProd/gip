package stack

import (
	"testing"
)

func TestStack(t *testing.T) {
	var s Stack[int]

	_, ok := s.Pop()
	if ok {
		t.Fatal("pop empty stack returns non-nil")
	}

	i1 := 1
	i2 := 2
	i3 := 3

	s.Push(i1)
	s.Push(i2)

	v, _ := s.Pop()
	if v != i2 {
		t.Fatal("pop returns bad value")
	}

	s.Push(i3)

	v, _ = s.Pop()
	if v != i3 {
		t.Fatal("pop returns bad value")
	}

	v, _ = s.Pop()
	if v != i1 {
		t.Fatal("pop returns bad value")
	}

	_, ok = s.Pop()
	if ok {
		t.Fatal("pop empty stack returns non-nil")
	}
}

func TestStack2(t *testing.T) {
	var s Stack[int]

	for i := 0; i < 5; i++ {
		s.Push(i)
	}

	for i := 4; i >= 0; i-- {
		v, ok := s.Pop()
		if !ok {
			t.Fatal("no value")
		}

		if v != i {
			t.Fatalf("bad value: %d, expected: %d", v, i)
		}
	}
}

func BenchmarkLockFree(b *testing.B) {
	var s Stack[int]

	for i := 0; i < b.N; i++ {
		s.Push(i)
		s.Pop()
	}
}

func BenchmarkMutex(b *testing.B) {
	s := NewMutexStack[int]()

	for i := 0; i < b.N; i++ {
		s.Push(i)
		s.Pop()
	}
}
