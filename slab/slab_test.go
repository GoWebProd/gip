package slab

import (
	"runtime"
	"sync"
	"testing"

	"github.com/couchbase/go-slab"

	"github.com/GoWebProd/gip/allocator"
	"github.com/GoWebProd/gip/pool"
)

func TestMin(t *testing.T) {
	s := New(8, 64, 2)

	data := s.Get(2)
	d := *data

	if len(d) != 2 || cap(d) != 8 {
		t.Fatalf("bad slice: %d-%d, must be 2-8", len(d), cap(d))
	}

	s.Put(data)
}

func TestMax(t *testing.T) {
	s := New(8, 64, 2)

	data := s.Get(128)
	d := *data

	if len(d) != 128 || cap(d) != 128 {
		t.Fatalf("bad slice: %d-%d, must be 128-128", len(d), cap(d))
	}

	s.Put(data)
}

func TestMed(t *testing.T) {
	s := New(8, 64, 2)

	data := s.Get(27)
	d := *data

	if len(d) != 27 || cap(d) != 32 {
		t.Fatalf("bad slice: %d-%d, must be 27-32", len(d), cap(d))
	}

	s.Put(data)
}

func TestSmallGrowFactor(t *testing.T) {
	s := New(8, 64, 1.01)

	for i := 0; i < len(s.pools); i++ {
		if s.sizes[i] != i+8 {
			t.Fatalf("bad increments on small grow factor: %+v", &s.pools)
		}
	}
}

func BenchmarkOurPool(b *testing.B) {
	s := pool.Pool[[]byte]{}
	wg := sync.WaitGroup{}

	for i := 0; i < runtime.GOMAXPROCS(0)*10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := 0; i < b.N; i++ {
				d := s.Get()
				if d == nil {
					d = new([]byte)
					*d = allocator.Alloc(4096)
				}

				_ = (*d)[:i%4096]

				s.Put(d)
			}
		}()
	}

	wg.Wait()
}

func BenchmarkOur(b *testing.B) {
	s := New(4096, 4096, 1.1)
	wg := sync.WaitGroup{}

	for i := 0; i < runtime.GOMAXPROCS(0)*10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := 0; i < b.N; i++ {
				d := s.Get(i % 4096)

				s.Put(d)
			}
		}()
	}

	wg.Wait()
}

func BenchmarkPool(b *testing.B) {
	wg := sync.WaitGroup{}
	pool := sync.Pool{
		New: func() interface{} {
			return make([]byte, 4096)
		},
	}

	for i := 0; i < runtime.GOMAXPROCS(0)*10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := 0; i < b.N; i++ {
				d := pool.Get().([]byte)
				_ = d[:i%4096]

				pool.Put(d)
			}
		}()
	}

	wg.Wait()
}

func BenchmarkMalloc(b *testing.B) {
	wg := sync.WaitGroup{}

	for i := 0; i < runtime.GOMAXPROCS(0)*10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := 0; i < b.N; i++ {
				d := allocator.Alloc(1 + i%4096)

				allocator.Free(d)
			}
		}()
	}

	wg.Wait()
}

func BenchmarkChannel(b *testing.B) {
	wg := sync.WaitGroup{}
	n := runtime.GOMAXPROCS(0) * 10
	pool := make(chan []byte, n)

	for i := 0; i < n; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := 0; i < b.N; i++ {
				var d []byte

				select {
				case d = <-pool:

				default:
					d = make([]byte, 4096)
				}

				_ = d[:i%4096]

				pool <- d
			}
		}()
	}

	wg.Wait()
}

func BenchmarkCouchbase(b *testing.B) {
	s := slab.NewArena(48, 4096, 1.1, nil)
	m := sync.Mutex{}
	wg := sync.WaitGroup{}

	for i := 0; i < runtime.GOMAXPROCS(0)*10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := 0; i < b.N; i++ {
				m.Lock()

				d := s.Alloc(i % 4096)

				s.DecRef(d)
				m.Unlock()
			}
		}()
	}

	wg.Wait()
}

func BenchmarkMake(b *testing.B) {
	wg := sync.WaitGroup{}

	for i := 0; i < runtime.GOMAXPROCS(0)*10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for i := 0; i < b.N; i++ {
				d := make([]byte, i%4096)
				_ = d
			}
		}()
	}

	wg.Wait()
}
