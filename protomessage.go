package lazyproto

import (
	"reflect"
	"unsafe"
)

type ProtoMessage struct {
	// TODO: implement a custom byte view that only uses Data and Len. We don't need Cap
	// and can save memory by eliminating it.
	// Bytes are set to nil when the message is modified (i.e. marshaling will do
	// full field-by-field encoding).
	Bytes []byte
	//Flags  uint64
	Parent *ProtoMessage
}

func (m *ProtoMessage) IsModified() bool {
	return m.Bytes == nil
}

func (m *ProtoMessage) MarkModified() {
	if m.Bytes != nil {
		m.markModified()
	}
}

func (m *ProtoMessage) markModified() {
	m.Bytes = nil
	parent := m.Parent
	for parent != nil {
		if parent.IsModified() {
			break
		}
		parent.Bytes = nil
		parent = parent.Parent
	}
}

type BytesView struct {
	string
}

func BytesViewFromBytes(b []byte) (bv BytesView) {
	src := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest := (*reflect.StringHeader)(unsafe.Pointer(&bv))
	dest.Data = src.Data
	dest.Len = src.Len
	return bv
}

func BytesFromBytesView(bv BytesView) (b []byte) {
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	src := (*reflect.StringHeader)(unsafe.Pointer(&bv))
	dest.Data = src.Data
	dest.Len = src.Len
	dest.Cap = src.Len
	return b
}
