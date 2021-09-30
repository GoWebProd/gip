package spinlock

import (
	"sync"
	"testing"
)

func testLock(b *testing.B, threads int, l sync.Locker) {
	var wg sync.WaitGroup
	wg.Add(threads)

	var count1 int
	var count2 int

	for i := 0; i < threads; i++ {
		go func() {
			for i := 0; i < b.N; i++ {
				l.Lock()
				count1++
				count2 += 2
				l.Unlock()
			}
			wg.Done()
		}()
	}

	wg.Wait()

	if count1 != threads*b.N {
		b.Fatal("mismatch")
	}
	if count2 != threads*b.N*2 {
		b.Fatal("mismatch")
	}
}

func BenchmarkSpinlock_1(b *testing.B) {
	testLock(b, 1, &Locker{})
}

func BenchmarkSpinlock_6(b *testing.B) {
	testLock(b, 6, &Locker{})
}

func BenchmarkSpinlock_12(b *testing.B) {
	testLock(b, 12, &Locker{})
}

func BenchmarkMutex_1(b *testing.B) {
	testLock(b, 1, &sync.Mutex{})
}

func BenchmarkMutex_6(b *testing.B) {
	testLock(b, 6, &sync.Mutex{})
}

func BenchmarkMutex_12(b *testing.B) {
	testLock(b, 12, &sync.Mutex{})
}
