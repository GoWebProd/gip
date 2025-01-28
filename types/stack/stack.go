package stack

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/GoWebProd/gip/allocator"
)

type node[T any] struct {
	next uintptr
	data T
}

type Stack[T any] uintptr

func (head *Stack[T]) Init() {
	runtime.SetFinalizer(head, func(head *Stack[T]) {
		for {
			_, ok := head.Pop()
			if !ok {
				break
			}
		}
	})
}

func (head *Stack[T]) Push(v T) {
	node := allocator.AllocObject[node[T]]()
	node.data = v
	new := uintptr(unsafe.Pointer(node))

	for {
		old := atomic.LoadUintptr((*uintptr)(head))
		node.next = old

		if atomic.CompareAndSwapUintptr((*uintptr)(head), old, new) {
			break
		}
	}
}

func (head *Stack[T]) Pop() (T, bool) {
	var v T

	for {
		old := atomic.LoadUintptr((*uintptr)(head))
		if old == 0 {
			return v, false
		}

		node := (*node[T])(unsafe.Pointer(old))
		next := atomic.LoadUintptr(&node.next)

		if atomic.CompareAndSwapUintptr((*uintptr)(head), old, next) {
			allocator.FreeObject(node)

			return node.data, true
		}
	}
}
