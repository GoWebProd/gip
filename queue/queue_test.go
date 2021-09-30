package queue

import (
	"sync"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := New[string]()
	test1 := "test1"
	test2 := "test2"
	test3 := "test3"

	if str, ok := queue.Take(); ok {
		t.Fatalf("queue must be empty, but returned: %s", *str)
	}

	queue.Put(&test1)
	queue.Put(&test2)

	str, ok := queue.Take()
	if !ok {
		t.Fatalf("queue must be not empty, but data not returned")
	}

	if *str != test1 {
		t.Fatalf("string must be %q but returned %q", test1, *str)
	}

	queue.Put(&test3)

	str, ok = queue.Take()
	if !ok {
		t.Fatalf("queue must be not empty, but data not returned")
	}

	if *str != test2 {
		t.Fatalf("string must be %q but returned %q", test2, *str)
	}

	str, ok = queue.Take()
	if !ok {
		t.Fatalf("queue must be not empty, but data not returned")
	}

	if *str != test3 {
		t.Fatalf("string must be %q but returned %q", test2, *str)
	}

	if str, ok = queue.Take(); ok {
		t.Fatalf("queue must be empty, but returned: %s", *str)
	}
}

func TestParallel(t *testing.T) {
	const n = 100
	const elements = 100000

	queue := New[int]()
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	m := make(map[int]struct{})

	for i := 0; i < n; i++ {
		wg.Add(2)

		n := elements * 100 * (i + 1)

		go func(base int) {
			defer wg.Done()

			for i := 0; i < elements; i++ {
				item := base + i

				queue.Put(&item)
			}
		}(n)

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
	queue := New[int]()

	for i := 0; i < b.N; i++ {
		queue.Put(&i)

		v, ok := queue.Take()
		if !ok || *v != i {
			b.Fatal("bad value")
		}
	}
}
