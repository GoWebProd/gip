package cond

import (
	"sync"
	"testing"
	"time"
)

func TestSingle(t *testing.T) {
	var (
		c Single
		v int
	)

	go func() {
		time.Sleep(10 * time.Millisecond)

		v = 1

		c.Done()
	}()

	c.Wait()

	if v != 1 {
		t.Fatalf("bad value: %d", v)
	}

	c.Wait()
}

func TestMulti(t *testing.T) {
	var (
		c   Multi
		v   int
		wg  sync.WaitGroup
		err bool
	)

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			c.Wait()

			if v != 1 {
				err = true
			}
		}()
	}

	time.Sleep(10 * time.Millisecond)

	v = 1

	c.Done()
	wg.Wait()

	if err {
		t.Fatal("one of goroutines failed")
	}

	c.Wait()
}
