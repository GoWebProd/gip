package queue

import (
	"sync"
	"testing"
)

func TestBuffer(t *testing.T) {
	buffer := New[string](2)
	test1 := "test1"
	test2 := "test2"
	test3 := "test3"

	if str, ok := buffer.Take(); ok {
		t.Fatalf("buffer must be empty, but returned: %s", *str)
	}

	if buffer.Size() != 0 {
		t.Fatalf("buffer must have size 0 but now %d", buffer.Size())
	}

	if !buffer.Put(&test1) {
		t.Fatal("can't put in buffer, but buffer must be not full")
	}

	if buffer.Size() != 1 {
		t.Fatalf("buffer must have size 1 but now %d", buffer.Size())
	}

	if !buffer.Put(&test2) {
		t.Fatal("can't put in buffer, but buffer must be not full")
	}

	if buffer.Size() != 2 {
		t.Fatalf("buffer must have size 2 but now %d", buffer.Size())
	}

	if buffer.Put(&test3) {
		t.Fatal("can put in buffer, but buffer must be full")
	}

	if buffer.Size() != 2 {
		t.Fatalf("buffer must have size 2 but now %d", buffer.Size())
	}

	str, ok := buffer.Take()
	if !ok {
		t.Fatalf("buffer must be not empty, but data not returned")
	}

	if *str != test1 {
		t.Fatalf("string must be %q but returned %q", test1, *str)
	}

	if buffer.Size() != 1 {
		t.Fatalf("buffer must have size 1 but now %d", buffer.Size())
	}

	str, ok = buffer.Take()
	if !ok {
		t.Fatalf("buffer must be not empty, but data not returned")
	}

	if *str != test2 {
		t.Fatalf("string must be %q but returned %q", test2, *str)
	}

	if buffer.Size() != 0 {
		t.Fatalf("buffer must have size 0 but now %d", buffer.Size())
	}

	if str, ok = buffer.Take(); ok {
		t.Fatalf("buffer must be empty, but returned: %s", *str)
	}

	if !buffer.Put(&test3) {
		t.Fatal("can't put in buffer, but buffer must be not full")
	}

	if buffer.Size() != 1 {
		t.Fatalf("buffer must have size 1 but now %d", buffer.Size())
	}
}

func TestParallel(t *testing.T) {
	const n = 100
	const elements = 1000

	queue := New[int](n * elements)
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	m := make(map[int]struct{})

	for i := 0; i < n; i++ {
		wg.Add(1)

		n := elements * 100 * (i + 1)

		if i == 0 {
			wg.Add(1)

			go func(base int) {
				defer wg.Done()

				for i := 0; i < n*elements; i++ {
					item := base + i

					queue.Put(&item)
				}
			}(n)
		}

		go func() {
			defer wg.Done()

			localM := make(map[int]struct{})

			for i := 0; i < elements; i++ {
				item, ok := queue.Take()
				if !ok {
					i--
					continue
				}

				localM[*item] = struct{}{}
			}

			mu.Lock()

			for k, v := range localM {
				m[k] = v
			}

			mu.Unlock()
		}()
	}

	wg.Wait()

	if len(m) != elements*n {
		t.Fatalf("bad elements count in map: %d, expected: %d", len(m), elements*n)
	}
}

func BenchmarkNew(b *testing.B) {
	queue := New[int](10)

	for i := 0; i < b.N; i++ {
		queue.Put(&i)

		v, ok := queue.Take()
		if !ok || *v != i {
			b.Fatal("bad value")
		}
	}
}
