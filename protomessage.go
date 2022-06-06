package lazyproto

import (
	"reflect"
	"unsafe"
)

const FlagsMessageModified = 1

type ProtoMessage struct {
	// TODO: implement a custom byte view that only uses Data and Len. We don't need Cap.
	Bytes  []byte
	Flags  uint64
	Parent *ProtoMessage
}

func (m *ProtoMessage) IsModified() bool {
	return m.Flags&FlagsMessageModified != 0
}

func (m *ProtoMessage) MarkModified() {
	if m.Flags&FlagsMessageModified == 0 {
		m.markModified()
	}
}

func (m *ProtoMessage) markModified() {
	m.Flags |= FlagsMessageModified
	parent := m.Parent
	for parent != nil {
		if parent.Flags&FlagsMessageModified != 0 {
			break
		}
		parent.Flags |= FlagsMessageModified
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
