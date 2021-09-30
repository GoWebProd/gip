package spinlock

import (
	"runtime"
	"sync/atomic"
)

type noCopy struct{}

// Locker is a spinlock implementation.
//
// A Locker must not be copied after first use.
type Locker struct {
	noCopy noCopy

	lock uintptr
}

// Lock locks l.
// If the lock is already in use, the calling goroutine
// blocks until the locker is available.
//go:nosplit
func (l *Locker) Lock() {
	for !atomic.CompareAndSwapUintptr(&l.lock, 0, 1) {
		runtime.Gosched()
	}
}

// Unlock unlocks l.
//go:nosplit
func (l *Locker) Unlock() {
	atomic.StoreUintptr(&l.lock, 0)
}
