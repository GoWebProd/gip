package pool

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/GoWebProd/gip/rtime"
)

var allPoolsMu sync.Mutex

//go:linkname runtime_LoadAcquintptr runtime/internal/atomic.LoadAcquintptr
func runtime_LoadAcquintptr(ptr *uintptr) uintptr

//go:linkname runtime_StoreReluintptr runtime/internal/atomic.StoreReluintptr
func runtime_StoreReluintptr(ptr *uintptr, val uintptr) uintptr

type noCopy struct{}

type Pool[T any] struct {
	noCopy noCopy

	local     unsafe.Pointer // local fixed-size per-P pool, actual type is [P]poolLocal
	localSize uintptr        // size of the local array
}

func (p *Pool[T]) pin() (*poolLocal[T], int) {
	pid := rtime.ProcPin()
	// In pinSlow we store to local and then to localSize, here we load in opposite order.
	// Since we've disabled preemption, GC cannot happen in between.
	// Thus here we must observe local at least as large localSize.
	// We can observe a newer/larger local, it is fine (we must observe its zero-initialized-ness).
	s := runtime_LoadAcquintptr(&p.localSize) // load-acquire
	l := p.local                              // load-consume
	if uintptr(pid) < s {
		return indexLocal[T](l, pid), pid
	}
	return p.pinSlow()
}

func (p *Pool[T]) pinSlow() (*poolLocal[T], int) {
	// Retry under the mutex.
	// Can not lock the mutex while pinned.
	rtime.ProcUnpin()
	allPoolsMu.Lock()
	defer allPoolsMu.Unlock()

	pid := rtime.ProcPin()
	// poolCleanup won't be called while we are pinned.
	s := p.localSize
	l := p.local

	if uintptr(pid) < s {
		return indexLocal[T](l, pid), pid
	}
	size := runtime.GOMAXPROCS(0)
	local := make([]poolLocal[T], size)

	atomic.StorePointer(&p.local, unsafe.Pointer(&local[0])) // store-release
	runtime_StoreReluintptr(&p.localSize, uintptr(size))     // store-release

	return &local[pid], pid
}

// Put adds x to the pool.
func (p *Pool[T]) Put(x *T) {
	if x == nil {
		return
	}

	l, _ := p.pin()

	if l.private == nil {
		l.private = x
		x = nil
	} else {
		l.shared.pushHead(x)
	}

	rtime.ProcUnpin()
}

// Get selects an arbitrary item from the Pool, removes it from the
// Pool, and returns it to the caller.
// Get may choose to ignore the pool and treat it as empty.
// Callers should not assume any relation between values passed to Put and
// the values returned by Get.
//
// If Get would otherwise return nil and p.New is non-nil, Get returns
// the result of calling p.New.
func (p *Pool[T]) Get() *T {
	l, pid := p.pin()
	x := l.private
	l.private = nil

	if x == nil {
		// Try to pop the head of the local shard. We prefer
		// the head over the tail for temporal locality of
		// reuse.
		x, _ = l.shared.popHead()
		if x == nil {
			x = p.getSlow(pid)
		}
	}

	rtime.ProcUnpin()

	return x
}

func (p *Pool[T]) getSlow(pid int) *T {
	// See the comment in pin regarding ordering of the loads.
	size := runtime_LoadAcquintptr(&p.localSize) // load-acquire
	locals := p.local                            // load-consume

	// Try to steal one element from other procs.
	for i := 0; i < int(size); i++ {
		l := indexLocal[T](locals, (pid+i+1)%int(size))

		if x, _ := l.shared.popTail(); x != nil {
			return x
		}
	}

	return nil
}
