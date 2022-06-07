package sizedstream

import (
	"google.golang.org/protobuf/encoding/protowire"
)

type ProtoStream struct {
	buf          []byte
	ofs          int
	scratchArray [20]byte
}

func NewProtoStream() *ProtoStream {
	s := &ProtoStream{}
	s.buf = make([]byte, 1024*200)
	return s
}

func (s *ProtoStream) Reset() {
	s.ofs = 0
}

func (s *ProtoStream) BufferBytes() ([]byte, error) {
	return s.buf[:s.ofs], nil
}

func (s *ProtoStream) Len() int {
	return s.ofs
}

type PreparedUint32Key byte
type PreparedFixed64Key byte

// PreparedKey of string,bytes or embedded wire type.
type PreparedKey byte

func PrepareField(fieldNumber int, wireType protowire.Type) PreparedKey {
	b := protowire.AppendVarint([]byte{}, uint64(fieldNumber)<<3+uint64(wireType))
	if len(b) > 1 {
		panic("fieldNumber too high, max field number 15 is supported by PrepareField")
	}
	return PreparedKey(b[0])
}

func PrepareStringField(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.BytesType)
}

func PrepareEmbeddedField(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.BytesType)
}

func PrepareUint32Field(fieldNumber int) PreparedUint32Key {
	return PreparedUint32Key(PrepareField(fieldNumber, protowire.VarintType))
}
func PrepareFixed64Field(fieldNumber int) PreparedFixed64Key {
	return PreparedFixed64Key(PrepareField(fieldNumber, protowire.Fixed64Type))
}

// Uint32Prepared writes a value of uint32 to the stream.
func (s *ProtoStream) Uint32Prepared(fieldKey PreparedUint32Key, value uint32) {
	if value == 0 {
		return
	}
	s.writeByte(byte(fieldKey))
	s.writeVarint(uint64(value))
	return
}

// Fixed64Prepared writes a value of fixed64 to the stream.
func (s *ProtoStream) Fixed64Prepared(fieldKey PreparedFixed64Key, value uint64) {
	if value == 0 {
		return
	}

	s.writeByte(byte(fieldKey))
	s.writeFixed64(value)
}

// StringPrepared writes a string to the stream.
func (s *ProtoStream) StringPrepared(key PreparedKey, value string) {
	vlen := len(value)
	if vlen == 0 {
		return
	}

	if vlen < 1<<7 {
		s.buf[s.ofs] = byte(key)
		s.buf[s.ofs+1] = byte(vlen)
		s.ofs += 2
	} else {
		s.writeByte(byte(key))
		s.writeVarint(uint64(vlen))
	}
	copy(s.buf[s.ofs:], value)
	s.ofs += len(value)
}

type EmbeddedToken int

func (s *ProtoStream) BeginEmbedded() EmbeddedToken {
	ofs := s.ofs
	s.ofs += 2
	return EmbeddedToken(ofs)
}

func (s *ProtoStream) EndEmbeddedPrepared(
	startPos EmbeddedToken, fieldKey PreparedKey,
) {
	embeddedSize := s.Len() - int(startPos) - 2

	//s.writeByte(byte(fieldKey))
	//s.writeVarint(uint64(embeddedSize))

	if uint64(embeddedSize) < 1<<7 {
		s.buf[startPos] = byte(fieldKey)
		s.buf[startPos+1] = byte(embeddedSize)
		return
	}

	scratchBuffer := s.scratchArray[:0]
	scratchBuffer = append(scratchBuffer, byte(fieldKey))
	scratchBuffer = protowire.AppendVarint(scratchBuffer, uint64(embeddedSize))

	shift := len(scratchBuffer) - 2
	copy(s.buf[int(startPos)+2+shift:], s.buf[startPos+2:s.ofs])
	copy(s.buf[startPos:], scratchBuffer)

	//s.buf = append(
	//	s.buf[:startPos],
	//	append(scratchBuffer, s.buf[startPos+2:]...)...,
	//)
	s.ofs += shift
}

func (s *ProtoStream) writeByte(b byte) {
	s.buf[s.ofs] = b
	s.ofs++
}

// Raw writes the byte sequence as is. Used to write prepared embedded byte sequences.
func (s *ProtoStream) Raw(b []byte) {
	copy(s.buf[s.ofs:], b)
	s.ofs += len(b)
}

func (s *ProtoStream) writeFixed64(v uint64) {
	s.buf[s.ofs] = byte(v >> 0)
	s.buf[s.ofs+1] = byte(v >> 8)
	s.buf[s.ofs+2] = byte(v >> 16)
	s.buf[s.ofs+3] = byte(v >> 24)
	s.buf[s.ofs+4] = byte(v >> 32)
	s.buf[s.ofs+5] = byte(v >> 40)
	s.buf[s.ofs+6] = byte(v >> 48)
	s.buf[s.ofs+7] = byte(v >> 56)
	s.ofs += 8
}

func (s *ProtoStream) writeVarint(v uint64) {
	switch {
	case v < 1<<7:
		s.buf[s.ofs] = byte(v)
		s.ofs += 1

	case v < 1<<14:
		s.buf[s.ofs] = byte((v>>0)&0x7f | 0x80)
		s.buf[s.ofs+1] = byte(v >> 7)
		s.ofs += 2

	case v < 1<<21:
		s.buf[s.ofs] = byte((v>>0)&0x7f | 0x80)
		s.buf[s.ofs+1] = byte((v>>7)&0x7f | 0x80)
		s.buf[s.ofs+2] = byte(v >> 14)
		s.ofs += 3

	case v < 1<<28:
		s.buf[s.ofs] = byte((v>>0)&0x7f | 0x80)
		s.buf[s.ofs+1] = byte((v>>7)&0x7f | 0x80)
		s.buf[s.ofs+2] = byte((v>>14)&0x7f | 0x80)
		s.buf[s.ofs+3] = byte(v >> 21)
		s.ofs += 4

	case v < 1<<35:
		s.buf[s.ofs] = byte((v>>0)&0x7f | 0x80)
		s.buf[s.ofs+1] = byte((v>>7)&0x7f | 0x80)
		s.buf[s.ofs+2] = byte((v>>14)&0x7f | 0x80)
		s.buf[s.ofs+3] = byte((v>>21)&0x7f | 0x80)
		s.buf[s.ofs+4] = byte(v >> 28)
		s.ofs += 5

	case v < 1<<42:
		s.buf[s.ofs] = byte((v>>0)&0x7f | 0x80)
		s.buf[s.ofs+1] = byte((v>>7)&0x7f | 0x80)
		s.buf[s.ofs+2] = byte((v>>14)&0x7f | 0x80)
		s.buf[s.ofs+3] = byte((v>>21)&0x7f | 0x80)
		s.buf[s.ofs+4] = byte((v>>28)&0x7f | 0x80)
		s.buf[s.ofs+5] = byte(v >> 35)
		s.ofs += 6

	case v < 1<<49:
		s.buf[s.ofs] = byte((v>>0)&0x7f | 0x80)
		s.buf[s.ofs+1] = byte((v>>7)&0x7f | 0x80)
		s.buf[s.ofs+2] = byte((v>>14)&0x7f | 0x80)
		s.buf[s.ofs+3] = byte((v>>21)&0x7f | 0x80)
		s.buf[s.ofs+4] = byte((v>>28)&0x7f | 0x80)
		s.buf[s.ofs+5] = byte((v>>35)&0x7f | 0x80)
		s.buf[s.ofs+6] = byte(v >> 42)
		s.ofs += 7

	case v < 1<<56:
		s.buf[s.ofs] = byte((v>>0)&0x7f | 0x80)
		s.buf[s.ofs+1] = byte((v>>7)&0x7f | 0x80)
		s.buf[s.ofs+2] = byte((v>>14)&0x7f | 0x80)
		s.buf[s.ofs+3] = byte((v>>21)&0x7f | 0x80)
		s.buf[s.ofs+4] = byte((v>>28)&0x7f | 0x80)
		s.buf[s.ofs+5] = byte((v>>35)&0x7f | 0x80)
		s.buf[s.ofs+6] = byte((v>>42)&0x7f | 0x80)
		s.buf[s.ofs+7] = byte(v >> 49)
		s.ofs += 8

	case v < 1<<63:
		s.buf[s.ofs] = byte((v>>0)&0x7f | 0x80)
		s.buf[s.ofs+1] = byte((v>>7)&0x7f | 0x80)
		s.buf[s.ofs+2] = byte((v>>14)&0x7f | 0x80)
		s.buf[s.ofs+3] = byte((v>>21)&0x7f | 0x80)
		s.buf[s.ofs+4] = byte((v>>28)&0x7f | 0x80)
		s.buf[s.ofs+5] = byte((v>>35)&0x7f | 0x80)
		s.buf[s.ofs+6] = byte((v>>42)&0x7f | 0x80)
		s.buf[s.ofs+7] = byte((v>>49)&0x7f | 0x80)
		s.buf[s.ofs+8] = byte(v >> 56)
		s.ofs += 9

	default:
		s.buf[s.ofs] = byte((v>>0)&0x7f | 0x80)
		s.buf[s.ofs+1] = byte((v>>7)&0x7f | 0x80)
		s.buf[s.ofs+2] = byte((v>>14)&0x7f | 0x80)
		s.buf[s.ofs+3] = byte((v>>21)&0x7f | 0x80)
		s.buf[s.ofs+4] = byte((v>>28)&0x7f | 0x80)
		s.buf[s.ofs+5] = byte((v>>35)&0x7f | 0x80)
		s.buf[s.ofs+6] = byte((v>>42)&0x7f | 0x80)
		s.buf[s.ofs+7] = byte((v>>49)&0x7f | 0x80)
		s.buf[s.ofs+8] = byte((v>>56)&0x7f | 0x80)
		s.buf[s.ofs+9] = 1
		s.ofs += 10
	}
}
