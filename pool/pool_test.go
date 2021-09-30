package pool

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
	_ "unsafe"

	"github.com/GoWebProd/gip/allocator"
	"github.com/GoWebProd/gip/rtime"
)

func TestParallel(t *testing.T) {
	proc := runtime.GOMAXPROCS(0)
	p := Pool[[]byte]{}
	zeroBytes := make([]byte, 128)

	var (
		wg sync.WaitGroup
		e  bool
	)

	for i := 0; i < proc*100; i++ {
		wg.Add(1)

		number := i
		template := fmt.Sprintf("goroutine %d", number)

		go func() {
			defer wg.Done()

			for i := 0; i < 10000; i++ {
				b := p.Get()
				if b == nil {
					b = new([]byte)
					*b = allocator.Alloc(128)
				}

				if !bytes.Equal(*b, zeroBytes) {
					t.Errorf("bad value: %v", string((*b)[:len(template)]))

					e = true
				}

				copy(*b, []byte(template))
				time.Sleep(time.Millisecond)
				copy(*b, zeroBytes)
				p.Put(b)
			}
		}()
	}

	wg.Wait()

	if e {
		t.Fatal("error raised")
	}
}

func TestPool(t *testing.T) {
	p := Pool[[]byte]{}

	b1 := p.Get()
	if b1 != nil {
		t.Fatalf("bad length: %d", len(*b1))
	}

	b1 = new([]byte)
	*b1 = allocator.Alloc(5)

	(*b1)[0] = 'h'
	(*b1)[1] = 'e'
	(*b1)[2] = 'l'
	(*b1)[3] = 'l'
	(*b1)[4] = 'o'

	b2 := p.Get()
	if b2 != nil {
		t.Fatalf("bad length: %d", len(*b2))
	}

	p.Put(b1)

	b3 := p.Get()
	if len(*b3) != 5 {
		t.Fatalf("bad length: %d", len(*b1))
	}

	if !bytes.Equal(*b1, *b3) {
		t.Fatalf("b1 != b3: %v != %v", b1, b3)
	}
}

func BenchmarkDefault(b *testing.B) {
	t := sync.Pool{
		New: func() interface{} {
			return make([]byte, 128)
		},
	}

	for i := 0; i < b.N; i++ {
		x := t.Get().([]byte)
		copy(x, "test")
		t.Put(x)
	}
}

func BenchmarkOur(b *testing.B) {
	t := Pool[[]byte]{}

	for i := 0; i < b.N; i++ {
		x := t.Get()
		if x == nil {
			x = new([]byte)
			*x = allocator.Alloc(128)
		}

		copy(*x, "test")
		t.Put(x)
	}
}

func BenchmarkPin(b *testing.B) {
	g := runtime.GOMAXPROCS(0)

	for i := 0; i < b.N; i++ {
		n := rtime.ProcPin()
		if n >= g {
			b.Fatal(n)
		}

		rtime.ProcUnpin()
	}
}

func BenchmarkLock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		runtime.LockOSThread()
		runtime.UnlockOSThread()
	}
}
