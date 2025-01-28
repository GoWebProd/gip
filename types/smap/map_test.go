package smap

import (
	"sync"
	"testing"
)

func TestMap(t *testing.T) {
	var m Map[string, int]

	i1 := 1
	i2 := 2

	m.Store("test", i1)
	m.Store("test2", i2)

	val, ok := m.Load("test")
	if !ok {
		t.Fatal("value not found in map")
	}

	if val != i1 {
		t.Fatal("invalid value:", val)
	}
}

func BenchmarkMap(b *testing.B) {
	var m Map[string, int]

	for i := 0; i < b.N; i++ {
		m.Store("test", i)

		val, ok := m.Load("test")
		if !ok {
			b.Fatal("value not found in map")
		}

		if val != i {
			b.Fatal("invalid value:", val)
		}

		m.Delete("test")
	}
}

func BenchmarkSyncMap(b *testing.B) {
	var m sync.Map

	i1 := 1

	for i := 0; i < b.N; i++ {
		m.Store("test", &i1)

		val, ok := m.Load("test")
		if !ok {
			b.Fatal("value not found in map")
		}

		if *(val.(*int)) != i1 {
			b.Fatal("invalid value:", *(val.(*int)))
		}

		m.Delete("test")
	}
}
