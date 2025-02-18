package safe

import (
	"reflect"
	"unsafe"
)

var cleanSlice = make([]byte, 1024)

func Cleanup[T any](v *T) {
	var (
		data []byte
		t    T
	)

	size := int(unsafe.Sizeof(t))
	h := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	h.Len = size
	h.Cap = size
	h.Data = uintptr(Noescape(v))

	for i := 0; i < size; i += 1024 {
		copy(data[i:], cleanSlice)
	}
}
