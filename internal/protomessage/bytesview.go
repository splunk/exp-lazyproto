package protomessage

import (
	"reflect"
	"unsafe"
)

// BytesView as an immutable view into a byte sequence. This is very similar to []byte
// native slice. The difference is that it does not maintain a "capacity" value
// because the view is immutable, and it is impossible to grow the "length",
// so we don't need to keep "capacity" at all. This saves us some memory space.
type BytesView struct {
	Data unsafe.Pointer
	Len  int
}

// BytesViewFromBytes creates a BytesView from []byte slice. Note that the information
// about the capacity of the []byte slice is lost.
func BytesViewFromBytes(b []byte) (dest BytesView) {
	src := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest.Data = unsafe.Pointer(src.Data)
	dest.Len = src.Len
	return dest
}

// BytesFromBytesView creates a []byte from BytesView
func BytesFromBytesView(src BytesView) (b []byte) {
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest.Data = uintptr(src.Data)
	dest.Len = src.Len
	dest.Cap = src.Len
	return b
}
