package allocator

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/GoWebProd/gip/safe"

	"github.com/ebitengine/purego"
)

var (
	mallocx func(size uint32, flags uint32) unsafe.Pointer
	rallocx func(ptr unsafe.Pointer, size uint32, flags uint32) unsafe.Pointer
	free    func(ptr unsafe.Pointer)
)

func init() {
	var name string

	switch runtime.GOOS {
	case "darwin":
		name = "libjemalloc.dylib"
	case "unix":
		name = "libjemalloc.so"
	default:
		panic("jemalloc no supported on " + runtime.GOOS)
	}

	lib, err := purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		panic("failed to load jemalloc: " + err.Error())
	}

	purego.RegisterLibFunc(&mallocx, lib, "mallocx")
	purego.RegisterLibFunc(&rallocx, lib, "rallocx")
	purego.RegisterLibFunc(&free, lib, "free")
}

func Alloc(size uint32) []byte {
	ptr := mallocx(size, 0x40)
	if ptr == nil {
		panic("out of memory")
	}

	var data []byte

	h := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	h.Len = int(size)
	h.Cap = int(size)
	h.Data = uintptr(ptr)

	return data
}

func Realloc(data *[]byte, size uint32) {
	h := (*reflect.SliceHeader)(unsafe.Pointer(data))
	if h.Cap == 0 {
		*data = Alloc(size)[:0]

		return
	}

	ptr := rallocx(unsafe.Pointer(h.Data), size, 0x40)
	if ptr == nil {
		panic("out of memory")
	}

	h.Data = uintptr(ptr)
	h.Cap = int(size)
}

func AllocObject[T any]() *T {
	var obj T

	size := unsafe.Sizeof(obj)

	ptr := mallocx(uint32(size), 0x40)
	if ptr == nil {
		panic("out of memory")
	}

	return (*T)(ptr)
}

func Free(data []byte) {
	free(safe.Noescape(&data[0]))
}

func FreeObject[T any](o *T) {
	free(safe.Noescape(o))
}
