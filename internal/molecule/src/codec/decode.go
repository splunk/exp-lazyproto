// This file contains modifications from the original source code found in: https://github.com/jhump/protoreflect

package codec

import (
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"unsafe"

	"github.com/jhump/protoreflect/codec"
)

// ErrOverflow is returned when an integer is too large to be represented.
var ErrOverflow = errors.New("proto: integer overflow")

// ErrBadWireType is returned when decoding a wire-type from a buffer that
// is not valid.
var ErrBadWireType = errors.New("proto: bad wiretype")

var varintTypes = map[FieldType]bool{}
var fixed32Types = map[FieldType]bool{}
var fixed64Types = map[FieldType]bool{}

func init() {
	varintTypes[FieldType_BOOL] = true
	varintTypes[FieldType_INT32] = true
	varintTypes[FieldType_INT64] = true
	varintTypes[FieldType_UINT32] = true
	varintTypes[FieldType_UINT64] = true
	varintTypes[FieldType_SINT32] = true
	varintTypes[FieldType_SINT64] = true
	varintTypes[FieldType_ENUM] = true

	fixed32Types[FieldType_FIXED32] = true
	fixed32Types[FieldType_SFIXED32] = true
	fixed32Types[FieldType_FLOAT] = true

	fixed64Types[FieldType_FIXED64] = true
	fixed64Types[FieldType_SFIXED64] = true
	fixed64Types[FieldType_DOUBLE] = true
}

// DecodeVarint reads a varint-encoded integer from the Buffer.
// This is the format for the
// int32, int64, uint32, uint64, bool, and enum
// protocol buffer types.
//
// This implementation is inlined from https://github.com/dennwc/varint to avoid the call-site overhead
func (cb *Buffer) DecodeVarint() (uint64, error) {
	if cb.Len() > 0 && cb.buf[cb.index] < 1<<7 {
		// i == 0
		b := cb.buf[cb.index]
		cb.index++
		return uint64(b), nil
	}
	return cb.decodeVarintUniform()
}

func (cb *Buffer) decodeVarintUniform() (uint64, error) {
	const (
		step = 7
		bit  = 1 << 7
		mask = bit - 1
	)
	if cb.Len() >= 10 {
		// i == 0
		b := cb.buf[cb.index]
		if b < bit {
			cb.index++
			return uint64(b), nil
		}
		x := uint64(b & mask)
		var s uint = step

		// i == 1
		b = cb.buf[cb.index+1]
		if b < bit {
			cb.index += 2
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&mask) << s
		s += step

		// i == 2
		b = cb.buf[cb.index+2]
		if b < bit {
			cb.index += 3
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&mask) << s
		s += step

		// i == 3
		b = cb.buf[cb.index+3]
		if b < bit {
			cb.index += 4
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&mask) << s
		s += step

		// i == 4
		b = cb.buf[cb.index+4]
		if b < bit {
			cb.index += 5
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&mask) << s
		s += step

		// i == 5
		b = cb.buf[cb.index+5]
		if b < bit {
			cb.index += 6
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&mask) << s
		s += step

		// i == 6
		b = cb.buf[cb.index+6]
		if b < bit {
			cb.index += 7
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&mask) << s
		s += step

		// i == 7
		b = cb.buf[cb.index+7]
		if b < bit {
			cb.index += 8
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&mask) << s
		s += step

		// i == 8
		b = cb.buf[cb.index+8]
		if b < bit {
			cb.index += 9
			return x | uint64(b)<<s, nil
		}
		x |= uint64(b&mask) << s
		s += step

		// i == 9
		b = cb.buf[cb.index+9]
		if b < bit {
			if b > 1 {
				return 0, ErrOverflow
			}
			cb.index += 10
			return x | uint64(b)<<s, nil
		} else if cb.Len() == 10 {
			return 0, io.ErrUnexpectedEOF
		}
		for _, b := range cb.buf[cb.index+10:] {
			if b < bit {
				return 0, ErrOverflow
			}
		}
		return 0, io.ErrUnexpectedEOF
	}

	if cb.Len() == 0 {
		return 0, io.ErrUnexpectedEOF
	}

	// i == 0
	b := cb.buf[cb.index]
	if b < bit {
		cb.index++
		return uint64(b), nil
	} else if cb.Len() == 1 {
		return 0, io.ErrUnexpectedEOF
	}
	x := uint64(b & mask)
	var s uint = step

	// i == 1
	b = cb.buf[cb.index+1]
	if b < bit {
		cb.index += 2
		return x | uint64(b)<<s, nil
	} else if cb.Len() == 2 {
		return 0, io.ErrUnexpectedEOF
	}
	x |= uint64(b&mask) << s
	s += step

	// i == 2
	b = cb.buf[cb.index+2]
	if b < bit {
		cb.index += 3
		return x | uint64(b)<<s, nil
	} else if cb.Len() == 3 {
		return 0, io.ErrUnexpectedEOF
	}
	x |= uint64(b&mask) << s
	s += step

	// i == 3
	b = cb.buf[cb.index+3]
	if b < bit {
		cb.index += 4
		return x | uint64(b)<<s, nil
	} else if cb.Len() == 4 {
		return 0, io.ErrUnexpectedEOF
	}
	x |= uint64(b&mask) << s
	s += step

	// i == 4
	b = cb.buf[cb.index+4]
	if b < bit {
		cb.index += 5
		return x | uint64(b)<<s, nil
	} else if cb.Len() == 5 {
		return 0, io.ErrUnexpectedEOF
	}
	x |= uint64(b&mask) << s
	s += step

	// i == 5
	b = cb.buf[cb.index+5]
	if b < bit {
		cb.index += 6
		return x | uint64(b)<<s, nil
	} else if cb.Len() == 6 {
		return 0, io.ErrUnexpectedEOF
	}
	x |= uint64(b&mask) << s
	s += step

	// i == 6
	b = cb.buf[cb.index+6]
	if b < bit {
		cb.index += 7
		return x | uint64(b)<<s, nil
	} else if cb.Len() == 7 {
		return 0, io.ErrUnexpectedEOF
	}
	x |= uint64(b&mask) << s
	s += step

	// i == 7
	b = cb.buf[cb.index+7]
	if b < bit {
		cb.index += 8
		return x | uint64(b)<<s, nil
	} else if cb.Len() == 8 {
		return 0, io.ErrUnexpectedEOF
	}
	x |= uint64(b&mask) << s
	s += step

	// i == 8
	b = cb.buf[cb.index+8]
	if b < bit {
		cb.index += 9
		return x | uint64(b)<<s, nil
	} else if cb.Len() == 9 {
		return 0, io.ErrUnexpectedEOF
	}
	x |= uint64(b&mask) << s
	s += step

	// i == 9
	b = cb.buf[cb.index+9]
	if b < bit {
		if b > 1 {
			return 0, ErrOverflow
		}
		cb.index += 10
		return x | uint64(b)<<s, nil
	} else if cb.Len() == 10 {
		return 0, io.ErrUnexpectedEOF
	}
	for _, b := range cb.buf[cb.index+10:] {
		if b < bit {
			return 0, ErrOverflow
		}
	}
	return 0, io.ErrUnexpectedEOF
}

func (cb *Buffer) SkipVarint() error {
	const (
		step = 7
		bit  = 1 << 7
		mask = bit - 1
	)
	if cb.Len() >= 10 {
		// i == 0
		b := cb.buf[cb.index]
		if b < bit {
			cb.index++
			return nil
		}

		// i == 1
		b = cb.buf[cb.index+1]
		if b < bit {
			cb.index += 2
			return nil
		}

		// i == 2
		b = cb.buf[cb.index+2]
		if b < bit {
			cb.index += 3
			return nil
		}

		// i == 3
		b = cb.buf[cb.index+3]
		if b < bit {
			cb.index += 4
			return nil
		}

		// i == 4
		b = cb.buf[cb.index+4]
		if b < bit {
			cb.index += 5
			return nil
		}

		// i == 5
		b = cb.buf[cb.index+5]
		if b < bit {
			cb.index += 6
			return nil
		}

		// i == 6
		b = cb.buf[cb.index+6]
		if b < bit {
			cb.index += 7
			return nil
		}

		// i == 7
		b = cb.buf[cb.index+7]
		if b < bit {
			cb.index += 8
			return nil
		}

		// i == 8
		b = cb.buf[cb.index+8]
		if b < bit {
			cb.index += 9
			return nil
		}

		// i == 9
		b = cb.buf[cb.index+9]
		if b < bit {
			if b > 1 {
				return ErrOverflow
			}
			cb.index += 10
			return nil
		} else if cb.Len() == 10 {
			return io.ErrUnexpectedEOF
		}
		for _, b := range cb.buf[cb.index+10:] {
			if b < bit {
				return ErrOverflow
			}
		}
		return io.ErrUnexpectedEOF
	}

	if cb.Len() == 0 {
		return io.ErrUnexpectedEOF
	}

	// i == 0
	b := cb.buf[cb.index]
	if b < bit {
		cb.index++
		return nil
	} else if cb.Len() == 1 {
		return io.ErrUnexpectedEOF
	}

	// i == 1
	b = cb.buf[cb.index+1]
	if b < bit {
		cb.index += 2
		return nil
	} else if cb.Len() == 2 {
		return io.ErrUnexpectedEOF
	}

	// i == 2
	b = cb.buf[cb.index+2]
	if b < bit {
		cb.index += 3
		return nil
	} else if cb.Len() == 3 {
		return io.ErrUnexpectedEOF
	}

	// i == 3
	b = cb.buf[cb.index+3]
	if b < bit {
		cb.index += 4
		return nil
	} else if cb.Len() == 4 {
		return io.ErrUnexpectedEOF
	}

	// i == 4
	b = cb.buf[cb.index+4]
	if b < bit {
		cb.index += 5
		return nil
	} else if cb.Len() == 5 {
		return io.ErrUnexpectedEOF
	}

	// i == 5
	b = cb.buf[cb.index+5]
	if b < bit {
		cb.index += 6
		return nil
	} else if cb.Len() == 6 {
		return io.ErrUnexpectedEOF
	}

	// i == 6
	b = cb.buf[cb.index+6]
	if b < bit {
		cb.index += 7
		return nil
	} else if cb.Len() == 7 {
		return io.ErrUnexpectedEOF
	}

	// i == 7
	b = cb.buf[cb.index+7]
	if b < bit {
		cb.index += 8
		return nil
	} else if cb.Len() == 8 {
		return io.ErrUnexpectedEOF
	}

	// i == 8
	b = cb.buf[cb.index+8]
	if b < bit {
		cb.index += 9
		return nil
	} else if cb.Len() == 9 {
		return io.ErrUnexpectedEOF
	}

	// i == 9
	b = cb.buf[cb.index+9]
	if b < bit {
		if b > 1 {
			return ErrOverflow
		}
		cb.index += 10
		return nil
	} else if cb.Len() == 10 {
		return io.ErrUnexpectedEOF
	}
	for _, b := range cb.buf[cb.index+10:] {
		if b < bit {
			return ErrOverflow
		}
	}
	return io.ErrUnexpectedEOF
}

func (cb *Buffer) SkipFixed64() error {
	i := cb.index + 8
	if i < 0 || i > len(cb.buf) {
		return io.ErrUnexpectedEOF
	}
	cb.index = i
	return nil
}

// DecodeFixed64 reads a 64-bit integer from the Buffer.
// This is the format for the
// fixed64, sfixed64, and double protocol buffer types.
func (cb *Buffer) DecodeFixed64() (x uint64, err error) {
	// x, err already 0
	i := cb.index + 8
	if i < 0 || i > len(cb.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	cb.index = i

	x = uint64(cb.buf[i-8])
	x |= uint64(cb.buf[i-7]) << 8
	x |= uint64(cb.buf[i-6]) << 16
	x |= uint64(cb.buf[i-5]) << 24
	x |= uint64(cb.buf[i-4]) << 32
	x |= uint64(cb.buf[i-3]) << 40
	x |= uint64(cb.buf[i-2]) << 48
	x |= uint64(cb.buf[i-1]) << 56
	return
}

func (cb *Buffer) SkipFixed32() error {
	i := cb.index + 4
	if i < 0 || i > len(cb.buf) {
		return io.ErrUnexpectedEOF
	}
	cb.index = i
	return nil
}

// DecodeFixed32 reads a 32-bit integer from the Buffer.
// This is the format for the
// fixed32, sfixed32, and float protocol buffer types.
func (cb *Buffer) DecodeFixed32() (x uint32, err error) {
	// x, err already 0
	i := cb.index + 4
	if i < 0 || i > len(cb.buf) {
		err = io.ErrUnexpectedEOF
		return
	}
	cb.index = i

	x = uint32(cb.buf[i-4])
	x |= uint32(cb.buf[i-3]) << 8
	x |= uint32(cb.buf[i-2]) << 16
	x |= uint32(cb.buf[i-1]) << 24
	return
}

// DecodeZigZag32 decodes a signed 32-bit integer from the given
// zig-zag encoded value.
func DecodeZigZag32(v uint64) int32 {
	return int32((uint32(v) >> 1) ^ uint32((int32(v&1)<<31)>>31))
}

// DecodeZigZag64 decodes a signed 64-bit integer from the given
// zig-zag encoded value.
func DecodeZigZag64(v uint64) int64 {
	return int64((v >> 1) ^ uint64((int64(v&1)<<63)>>63))
}

// DecodeRawBytes reads a count-delimited byte buffer from the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (cb *Buffer) DecodeRawBytes() (buf []byte, err error) {
	var n uint64
	if cb.Len() > 0 && cb.buf[cb.index] < 1<<7 {
		// Inline the DecodeVarint for the case when it is one byte (hot path).
		n = uint64(cb.buf[cb.index])
		cb.index++
	} else {
		n, err = cb.decodeVarintUniform()
		if err != nil {
			return nil, err
		}
	}

	nb := int(n)
	if nb < 0 {
		return nil, fmt.Errorf("proto: bad byte length %d", nb)
	}
	end := cb.index + nb
	if end < cb.index || end > len(cb.buf) {
		return nil, io.ErrUnexpectedEOF
	}

	// We set a cap on the returned slice equal to the length of the buffer so that it is not possible
	// to read past the end of this slice
	buf = cb.buf[cb.index:end:end]
	cb.index = end
	return
}

// SkipRawBytes skips a count-delimited byte buffer from the Buffer.
// This is the format used for the bytes protocol buffer
// type and for embedded messages.
func (cb *Buffer) SkipRawBytes() (err error) {
	var n uint64
	if cb.Len() > 0 && cb.buf[cb.index] < 1<<7 {
		// Inline the DecodeVarint for the case when it is one byte (hot path).
		n = uint64(cb.buf[cb.index])
		cb.index++
	} else {
		n, err = cb.decodeVarintUniform()
		if err != nil {
			return err
		}
	}

	nb := int(n)
	if nb < 0 {
		return fmt.Errorf("proto: bad byte length %d", nb)
	}
	end := cb.index + nb
	if end < cb.index || end > len(cb.buf) {
		return io.ErrUnexpectedEOF
	}

	cb.index = end

	return nil
}

// ReadGroup reads the input until a "group end" tag is found
// and returns the data up to that point. Subsequent reads from
// the buffer will read data after the group end tag. If alloc
// is true, the data is copied to a new slice before being returned.
// Otherwise, the returned slice is a view into the buffer's
// underlying byte slice.
//
// This function correctly handles nested groups: if a "group start"
// tag is found, then that group's end tag will be included in the
// returned data.
func (cb *Buffer) ReadGroup(alloc bool) ([]byte, error) {
	var groupEnd, dataEnd int
	groupEnd, dataEnd, err := cb.findGroupEnd()
	if err != nil {
		return nil, err
	}
	var results []byte
	if !alloc {
		results = cb.buf[cb.index:dataEnd]
	} else {
		results = make([]byte, dataEnd-cb.index)
		copy(results, cb.buf[cb.index:])
	}
	cb.index = groupEnd
	return results, nil
}

// SkipGroup is like ReadGroup, except that it discards the
// data and just advances the buffer to point to the input
// right *after* the "group end" tag.
func (cb *Buffer) SkipGroup() error {
	groupEnd, _, err := cb.findGroupEnd()
	if err != nil {
		return err
	}
	cb.index = groupEnd
	return nil
}

func (cb *Buffer) findGroupEnd() (groupEnd int, dataEnd int, err error) {
	bs := cb.buf
	start := cb.index
	defer func() {
		cb.index = start
	}()
	for {
		fieldStart := cb.index
		// read a field tag
		var v uint64
		v, err = cb.DecodeVarint()
		if err != nil {
			return 0, 0, err
		}
		_, wireType, err := AsTagAndWireType(v)
		if err != nil {
			return 0, 0, err
		}
		// skip past the field's data
		switch wireType {
		case WireFixed32:
			if err := cb.Skip(4); err != nil {
				return 0, 0, err
			}
		case WireFixed64:
			if err := cb.Skip(8); err != nil {
				return 0, 0, err
			}
		case WireVarint:
			// skip varint by finding last byte (has high bit unset)
			i := cb.index
			limit := i + 10 // varint cannot be >10 bytes
			for {
				if i >= limit {
					return 0, 0, ErrOverflow
				}
				if i >= len(bs) {
					return 0, 0, io.ErrUnexpectedEOF
				}
				if bs[i]&0x80 == 0 {
					break
				}
				i++
			}
			// TODO: This would only overflow if buffer length was MaxInt and we
			// read the last byte. This is not a real/feasible concern on 64-bit
			// systems. Something to worry about for 32-bit systems? Do we care?
			cb.index = i + 1
		case WireBytes:
			l, err := cb.DecodeVarint()
			if err != nil {
				return 0, 0, err
			}
			if err := cb.Skip(int(l)); err != nil {
				return 0, 0, err
			}
		case WireStartGroup:
			if err := cb.SkipGroup(); err != nil {
				return 0, 0, err
			}
		case WireEndGroup:
			return cb.index, fieldStart, nil
		default:
			return 0, 0, ErrBadWireType
		}
	}
}

// AsTagAndWireType converts the given varint in to a field number and wireType
//
// As of now, this function is inlined.  Please double check that any modifications do not modify the
// inline eligibility
func AsTagAndWireType(v uint64) (tag int32, wireType WireType, err error) {
	// rest is int32 tag number
	// low 7 bits is wire type
	wireType = WireType(v & 7)
	tag = int32(v >> 3)
	if tag <= 0 {
		err = ErrBadWireType // We return a constant error here as this allows the function to be inlined
	}

	return
}

// AsDouble interprets the value as a double.
func (cb *Buffer) AsDouble() (float64, error) {
	v, err := cb.DecodeFixed64()
	if err != nil {
		return 0, err
	}
	return math.Float64frombits(v), nil
}

// AsFloat interprets the value as a float.
func (cb *Buffer) AsFloat() (float32, error) {
	v, err := cb.DecodeFixed32()
	if err != nil {
		return 0, err
	}
	return math.Float32frombits(v), nil
}

// AsInt32 interprets the value as an int32.
func (cb *Buffer) AsInt32() (int32, error) {
	v, err := cb.DecodeVarint()
	if err != nil {
		return 0, err
	}
	s := int64(v)
	if s > math.MaxInt32 {
		return 0, fmt.Errorf("AsInt32: %d overflows int32", s)
	}
	if s < math.MinInt32 {
		return 0, fmt.Errorf("AsInt32: %d underflows int32", s)
	}
	return int32(v), nil
}

// AsInt64 interprets the value as an int64.
func (cb *Buffer) AsInt64() (int64, error) {
	v, err := cb.DecodeVarint()
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

// AsUint32 interprets the value as a uint32.
func (cb *Buffer) AsUint32() (uint32, error) {
	v, err := cb.DecodeVarint()
	if err != nil {
		return 0, err
	}
	if v > math.MaxUint32 {
		return 0, fmt.Errorf("AsUInt32: %d overflows uint32", v)
	}
	return uint32(v), nil
}

// AsUint64 interprets the value as a uint64.
func (cb *Buffer) AsUint64() (uint64, error) {
	v, err := cb.DecodeVarint()
	if err != nil {
		return 0, err
	}
	return v, nil
}

// AsSint32 interprets the value as a sint32.
func (cb *Buffer) AsSint32() (int32, error) {
	v, err := cb.DecodeVarint()
	if err != nil {
		return 0, err
	}
	if v > math.MaxUint32 {
		return 0, fmt.Errorf("AsSint32: %d overflows int32", v)
	}
	return codec.DecodeZigZag32(v), nil
}

// AsSint64 interprets the value as a sint64.
func (cb *Buffer) AsSint64() (int64, error) {
	v, err := cb.DecodeVarint()
	if err != nil {
		return 0, err
	}
	return codec.DecodeZigZag64(v), nil
}

// AsFixed32 interprets the value as a fixed32.
func (cb *Buffer) AsFixed32() (uint32, error) {
	v, err := cb.DecodeFixed32()
	if err != nil {
		return 0, err
	}
	return v, nil
}

// AsFixed64 interprets the value as a fixed64.
func (cb *Buffer) AsFixed64() (uint64, error) {
	v, err := cb.DecodeFixed64()
	if err != nil {
		return 0, err
	}
	return v, nil
}

// AsSFixed32 interprets the value as a SFixed32.
func (cb *Buffer) AsSFixed32() (int32, error) {
	v, err := cb.DecodeFixed32()
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

// AsSFixed64 interprets the value as a SFixed64.
func (cb *Buffer) AsSFixed64() (int64, error) {
	v, err := cb.DecodeFixed64()
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

// AsBool interprets the value as a bool.
func (cb *Buffer) AsBool() (bool, error) {
	v, err := cb.DecodeVarint()
	if err != nil {
		return false, err
	}
	return v == 1, nil
}

// AsStringUnsafe interprets the value as a string. The returned string is an unsafe view over
// the underlying bytes. Use AsStringSafe() to obtain a "safe" string that is a copy of the
// underlying data.
func (cb *Buffer) AsStringUnsafe() (s string, err error) {
	n, err := cb.DecodeVarint()
	if err != nil {
		return "", err
	}

	nb := int(n)
	if nb < 0 {
		return "", fmt.Errorf("proto: bad byte length %d", nb)
	}
	end := cb.index + nb
	if end < cb.index || end > len(cb.buf) {
		return "", io.ErrUnexpectedEOF
	}

	vh := (*reflect.SliceHeader)(unsafe.Pointer(&cb.buf))

	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sh.Data = vh.Data + uintptr(cb.index)
	sh.Len = end - cb.index

	cb.index = end

	return
}

// AsStringSafe interprets the value as a string by allocating a safe copy of the underlying data.
func (cb *Buffer) AsStringSafe() (string, error) {
	v, err := cb.DecodeRawBytes()
	if err != nil {
		return "", err
	}
	return string(v), nil
}

// AsBytesUnsafe interprets the value as a byte slice. The returned []byte is an unsafe view over
// the underlying bytes. Use AsBytesSafe() to obtain a "safe" [] that is a copy of the
// underlying data.
func (cb *Buffer) AsBytesUnsafe() ([]byte, error) {
	return cb.DecodeRawBytes()
}

// AsBytesSafe interprets the value as a byte slice by allocating a safe copy of the underlying data.
func (cb *Buffer) AsBytesSafe() ([]byte, error) {
	v, err := cb.DecodeRawBytes()
	if err != nil {
		return nil, err
	}
	return append([]byte(nil), v...), nil
}

func (cb *Buffer) SkipFieldByWireType(wireType WireType) error {
	switch wireType {
	case WireVarint:
		return cb.SkipVarint()
	case WireFixed64:
		return cb.SkipFixed64()
	case WireBytes:
		return cb.SkipRawBytes()
	case WireFixed32:
		return cb.SkipFixed32()
	default:
		return ErrBadWireType
	}
}
