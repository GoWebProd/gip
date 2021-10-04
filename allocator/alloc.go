package allocator

/*
#cgo LDFLAGS: -ljemalloc -lm -lstdc++ -pthread -ldl
#include <stdlib.h>
#include <jemalloc/jemalloc.h>
*/
import "C"

import (
	"reflect"
	"unsafe"
)

//go:linkname throw runtime.throw
func throw(s string)

func Alloc(size int) []byte {
	ptr := C.mallocx(C.size_t(size), 0x40)
	if ptr == nil {
		throw("out of memory")
	}

	var data []byte

	h := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	h.Len = size
	h.Cap = size
	h.Data = uintptr(ptr)

	return data
}

func AllocObject[T any]() *T {
	var obj T

	size := unsafe.Sizeof(obj)

	ptr := C.mallocx(C.size_t(size), 0x40)
	if ptr == nil {
		throw("out of memory")
	}

	return (*T)(ptr)
}

func Free(data []byte) {
	C.free(unsafe.Pointer(&data[0]))
}

func FreeObject[T any](o *T) {
	C.free(unsafe.Pointer(o))
}
