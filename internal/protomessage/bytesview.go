package protomessage

import (
	"reflect"
	"unsafe"
)

type BytesView struct {
	Data unsafe.Pointer
	Len  int
}

func BytesViewFromBytes(b []byte) (dest BytesView) {
	src := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest.Data = unsafe.Pointer(src.Data)
	dest.Len = src.Len
	return dest
}

func BytesFromBytesView(src BytesView) (b []byte) {
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest.Data = uintptr(src.Data)
	dest.Len = src.Len
	dest.Cap = src.Len
	return b
}
