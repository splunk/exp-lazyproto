package lazyproto

import "unsafe"

type PointerSliceArena struct {
	arena []unsafe.Pointer
}

func NewPointerSliceArena(capacity int) *PointerSliceArena {
	return &PointerSliceArena{
		arena: make([]unsafe.Pointer, 0, capacity),
	}
}

func (a *PointerSliceArena) Alloc(ptrCount int) unsafe.Pointer {
	a.arena = append(a.arena, make([]unsafe.Pointer, ptrCount)...)
	return unsafe.Pointer(&a.arena[len(a.arena)-ptrCount])
	//s := make([]unsafe.Pointer, ptrCount)
	//return unsafe.Pointer(&s[0])
}

func (a *PointerSliceArena) FreeAll() {
	a.arena = nil
}
