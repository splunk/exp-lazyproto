package lazyproto

import (
	"reflect"
	"unsafe"
)

type OneOf struct {
	lenAndFieldIdx int64
	capOrVal       int64
	ptr            unsafe.Pointer
}

// Number of bits to use for field index. This should be wide enough to fit all field indexes.
const fieldIdxBitCount = 5

// Bit mask for field index part of lenAndFieldIdx field.
const fieldIdxMask = (1 << fieldIdxBitCount) - 1

// MaxSliceLen is the maximum length of a slice-type that can be stored in OneOf.
// The length of Go slices can be at most maxint, however OneOf is not able to
// store lengths of maxint. Len field in OneOf uses typeFieldBitCount bits less
// than int, i.e. the maximum length of a slice stored in OneOf is
// maxint / (2^fieldIdxBitCount), which we calculate below.
const MaxSliceLen = int((^uint(0))>>1) >> fieldIdxBitCount

func NewOneOfString(v string, fieldIdx int) OneOf {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}

	return OneOf{
		ptr:            unsafe.Pointer(hdr.Data),
		lenAndFieldIdx: int64((hdr.Len << fieldIdxMask) | fieldIdx),
	}
}

func NewOneOfBytes(v []byte, fieldIdx int) OneOf {
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&v))
	if hdr.Len > MaxSliceLen {
		panic("maximum len exceeded")
	}

	return OneOf{
		ptr:            unsafe.Pointer(hdr.Data),
		lenAndFieldIdx: int64((hdr.Len << fieldIdxBitCount) | fieldIdx),
		capOrVal:       int64(hdr.Cap),
	}
}

func (v *OneOf) FieldIndex() int {
	return int(v.lenAndFieldIdx & fieldIdxMask)
}

// StringVal returns the stored string value.
func (v *OneOf) StringVal() (s string) {
	dest := (*reflect.StringHeader)(unsafe.Pointer(&s))
	dest.Data = uintptr(v.ptr)
	dest.Len = int(v.lenAndFieldIdx >> fieldIdxBitCount)
	return s
}

// BytesVal returns the stored byte slice.
func (v *OneOf) BytesVal() (b []byte) {
	dest := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	dest.Data = uintptr(v.ptr)
	dest.Len = int(v.lenAndFieldIdx >> fieldIdxBitCount)
	dest.Cap = int(v.capOrVal)
	return b
}
