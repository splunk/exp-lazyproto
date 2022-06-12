package protomessage

import "unsafe"

type DecodeBlock struct {
	RepeatPtrCount int
	RepeatPtrSlice []unsafe.Pointer
}
