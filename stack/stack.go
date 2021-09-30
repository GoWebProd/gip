package stack

import (
	"sync/atomic"
	"unsafe"
)

type node[T any] struct {
	next    uintptr
	data    *T
}

type Stack[T any] uintptr

func (head *Stack[T]) Push(v *T) {
	node := &node[T]{
		data: v,
	}
	new := uintptr(unsafe.Pointer(node))

	for {
		old := atomic.LoadUintptr((*uintptr)(head))
		node.next = old
		
		if atomic.CompareAndSwapUintptr((*uintptr)(head), old, new) {
			break
		}
	}
}

func (head *Stack[T]) Pop() *T{
	for {
		old := atomic.LoadUintptr((*uintptr)(head))
		if old == 0 {
			return nil
		}

		node := (*node[T])(unsafe.Pointer(old))
		next := atomic.LoadUintptr(&node.next)

		if atomic.CompareAndSwapUintptr((*uintptr)(head), old, next) {
			return node.data
		}
	}
}
