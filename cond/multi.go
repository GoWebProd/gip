package cond

import (
	"sync"
	"unsafe"

	"github.com/GoWebProd/gip/rtime"
	"github.com/GoWebProd/gip/stack"
)

type Multi struct {
	ready   bool
	waiters stack.MutexStack[unsafe.Pointer]
	mu      sync.Mutex
}

func (c *Multi) park(p1, p2 unsafe.Pointer) bool {
	c.waiters.PushUnlocked(p1)
	c.mu.Unlock()

	return true
}

func (c *Multi) Wait() {
	c.mu.Lock()

	if c.ready {
		c.mu.Unlock()

		return
	}

	rtime.GoPark(c.park, unsafe.Pointer(&c.mu), 0, 0, 1)
}

func (c *Multi) Done() {
	c.mu.Lock()

	c.ready = true

	for {
		g, ok := c.waiters.PopUnlocked()
		if !ok {
			break
		}

		rtime.GoReady(g, 1)
	}

	c.mu.Unlock()
}
