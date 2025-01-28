package iface

import (
	"unsafe"

	"github.com/GoWebProd/gip/safe"
)

type emptyInterface struct {
	typ unsafe.Pointer
	ptr unsafe.Pointer
}

//go:nosplit
func Interface[T any](p *T) any {
	var t any = (*T)(nil)

	iface := (*emptyInterface)(unsafe.Pointer(&t))
	iface.ptr = safe.Noescape(p)

	return t
}

func GetPointer(i any) unsafe.Pointer {
	iface := (*emptyInterface)(unsafe.Pointer(&i))

	return iface.ptr
}

func SetPointer(i *any, ptr unsafe.Pointer) {
	iface := (*emptyInterface)(unsafe.Pointer(i))
	iface.ptr = ptr
}

func Unpack(i any) (unsafe.Pointer, unsafe.Pointer) {
	iface := (*emptyInterface)(unsafe.Pointer(&i))

	return iface.typ, iface.ptr
}

func Build(typ unsafe.Pointer, ptr unsafe.Pointer) any {
	var i any

	iface := (*emptyInterface)(unsafe.Pointer(&i))
	iface.typ = typ
	iface.ptr = ptr

	return i
}
