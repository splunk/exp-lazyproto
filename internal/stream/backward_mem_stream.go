package stream

import (
	"bytes"
	"io"

	"google.golang.org/protobuf/encoding/protowire"
)

const firstBufSize = 40 * 4096

type ProtoStream struct {
	curBuf      []byte
	lastWritten int
	nextBufSize int

	readyBufs    [][]byte
	readyBufsLen int
}

func NewProtoStream() *ProtoStream {
	s := &ProtoStream{}
	s.curBuf = make([]byte, firstBufSize)
	s.lastWritten = firstBufSize
	s.nextBufSize = firstBufSize * 2
	s.Reset()
	return s
}

func (s *ProtoStream) Reset() {
	s.lastWritten = len(s.curBuf)
	s.readyBufs = nil
	s.readyBufsLen = 0
}

func (s *ProtoStream) WriteTo(out io.Writer) error {
	_, err := out.Write(s.curBuf[s.lastWritten:])
	if err != nil {
		return err
	}
	for i := len(s.readyBufs) - 1; i >= 0; i-- {
		_, err := out.Write(s.readyBufs[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ProtoStream) BufferBytes() ([]byte, error) {
	if s.readyBufs == nil {
		return s.curBuf[s.lastWritten:], nil
	}

	destBytes := make([]byte, 0, s.Len())
	destBuf := bytes.NewBuffer(destBytes)
	err := s.WriteTo(destBuf)
	if err != nil {
		return nil, err
	}
	return destBuf.Bytes(), nil
}

func (s *ProtoStream) Len() int {
	return s.readyBufsLen + len(s.curBuf) - s.lastWritten
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
	// Write in backward order.
	s.writeVarint(uint64(value))
	s.writeByte(byte(fieldKey))
}

// Fixed64Prepared writes a value of fixed64 to the stream.
func (s *ProtoStream) Fixed64Prepared(fieldKey PreparedFixed64Key, value uint64) {
	if value == 0 {
		return
	}

	// Write in backward order.
	s.writeFixed64(value)
	s.writeByte(byte(fieldKey))
}

// StringPrepared writes a string to the stream.
func (s *ProtoStream) StringPrepared(key PreparedKey, value string) {
	vlen := len(value)
	if vlen == 0 {
		return
	}

	// Write in backward order.
	//s.rawString(value)

	if vlen < s.lastWritten {
		s.lastWritten -= vlen
		copy(s.curBuf[s.lastWritten:], value)
	} else {
		bytesFit := s.lastWritten
		s.lastWritten = 0
		bytesRemaining := vlen - bytesFit
		copy(s.curBuf, value[bytesRemaining:])

		s.allocNewBuf()

		if bytesRemaining > 0 {
			s.lastWritten -= bytesRemaining
			copy(s.curBuf[s.lastWritten:], value[:bytesRemaining])
		}
	}

	if s.lastWritten < 11 {
		s.allocNewBuf()
	}

	if vlen < 1<<7 {
		s.curBuf[s.lastWritten-1] = byte(vlen)
		s.curBuf[s.lastWritten-2] = byte(key)
		s.lastWritten -= 2
		return
	}

	//s.writeVarint(uint64(vlen))
	//if vlen < 1<<7 {
	//	s.lastWritten--
	//	s.curBuf[s.lastWritten] = byte(vlen)
	//} else {
	s._writeVarintAbove127(uint64(vlen))
	//}

	//s.writeByte(byte(key))
	s.lastWritten--
	s.curBuf[s.lastWritten] = byte(key)
}

type EmbeddedToken int

func (s *ProtoStream) BeginEmbedded() EmbeddedToken {
	return EmbeddedToken(s.Len())
}

func (s *ProtoStream) EndEmbeddedPrepared(
	beginToken EmbeddedToken, fieldKey PreparedKey,
) {
	embeddedSize := s.Len() - int(beginToken)

	// Write in backward order.
	//s.writeVarint(uint64(embeddedSize))
	//s.writeByte(byte(fieldKey))

	if s.lastWritten < 11 {
		s.allocNewBuf()
	}

	//s.writeVarint(uint64(vlen))
	if embeddedSize < 1<<7 {
		s.curBuf[s.lastWritten-1] = byte(embeddedSize)
		s.curBuf[s.lastWritten-2] = byte(fieldKey)
		s.lastWritten -= 2
		return
	}

	s._writeVarintAbove127(uint64(embeddedSize))

	//s.writeByte(byte(key))
	s.lastWritten--
	s.curBuf[s.lastWritten] = byte(fieldKey)
}

func (s *ProtoStream) writeByte(b byte) {
	if s.lastWritten == 0 {
		s.allocNewBuf()
	}
	s.lastWritten--
	s.curBuf[s.lastWritten] = b
}

func (s *ProtoStream) appendReadyBuf(b []byte) {
	s.readyBufs = append(s.readyBufs, b)
	s.readyBufsLen += len(b)
}

// Raw writes the byte sequence as is. Used to write prepared embedded byte sequences.
func (s *ProtoStream) Raw(b []byte) {
	if len(b) >= firstBufSize {
		// Lots of bytes, just make it a ready buffer itself.
		// First we need to push current accumulated buffer:
		s.appendReadyBuf(s.curBuf[s.lastWritten:])

		if s.lastWritten > 0 {
			// Free space remains in the current buffer. Make it available
			// for future writes.
			s.curBuf = s.curBuf[:s.lastWritten:s.lastWritten]
		} else {
			// No more free space in the current buffer, just allocate a new one.
			s.curBuf = make([]byte, firstBufSize)
			s.lastWritten = firstBufSize
		}
		// Now add the new bytes.
		s.appendReadyBuf(b)
		return
	}

	bytesFit := s.lastWritten
	if bytesFit > len(b) {
		bytesFit = len(b)
	}
	s.lastWritten -= bytesFit
	copy(s.curBuf[s.lastWritten:], b[len(b)-bytesFit:])

	bytesRemaining := len(b) - bytesFit
	if bytesFit < 0 {
		panic("bad")
	}

	if s.lastWritten == 0 {
		s.allocNewBuf()
	}

	if bytesRemaining > 0 {
		s.lastWritten -= bytesRemaining
		copy(s.curBuf[s.lastWritten:], b[:bytesRemaining])
	}
}

// Raw writes the byte sequence as is. Used to write prepared embedded byte sequences.
func (s *ProtoStream) rawString(b string) {
}

func (s *ProtoStream) allocNewBuf() {
	s.appendReadyBuf(s.curBuf[s.lastWritten:])

	// Double the buf size on every new allocation to have constant amortized cost.
	s.curBuf = make([]byte, s.nextBufSize)
	s.lastWritten = s.nextBufSize
	s.nextBufSize *= 2
}

func (s *ProtoStream) writeFixed64(v uint64) {
	if s.lastWritten < 8 {
		s.allocNewBuf()
	}
	s.curBuf[s.lastWritten-8] = byte(v >> 0)
	s.curBuf[s.lastWritten-7] = byte(v >> 8)
	s.curBuf[s.lastWritten-6] = byte(v >> 16)
	s.curBuf[s.lastWritten-5] = byte(v >> 24)
	s.curBuf[s.lastWritten-4] = byte(v >> 32)
	s.curBuf[s.lastWritten-3] = byte(v >> 40)
	s.curBuf[s.lastWritten-2] = byte(v >> 48)
	s.curBuf[s.lastWritten-1] = byte(v >> 56)
	s.lastWritten -= 8
}

func (s *ProtoStream) writeVarint(v uint64) {
	if s.lastWritten > 0 && v < 1<<7 {
		s.lastWritten -= 1
		s.curBuf[s.lastWritten] = byte(v)
		return
	}
	if s.lastWritten < 10 {
		s.allocNewBuf()
	}
	s._writeVarintAbove127(v)
}

func (s *ProtoStream) writeVarintNoBufSizeCheck(v uint64) {
	if v < 1<<7 {
		s.lastWritten--
		s.curBuf[s.lastWritten] = byte(v)
		return
	}
	s._writeVarintAbove127(v)
}

func (s *ProtoStream) _writeVarintAbove127(v uint64) {
	//if s.lastWritten < 10 {
	//	s.allocNewBuf()
	//}

	switch {
	//case v < 1<<7:
	//	s.curBuf[s.lastWritten-1] = byte(v)
	//	s.lastWritten -= 1

	case v < 1<<14:
		s.curBuf[s.lastWritten-2] = byte((v>>0)&0x7f | 0x80)
		s.curBuf[s.lastWritten-1] = byte(v >> 7)
		s.lastWritten -= 2

	case v < 1<<21:
		s.curBuf[s.lastWritten-3] = byte((v>>0)&0x7f | 0x80)
		s.curBuf[s.lastWritten-2] = byte((v>>7)&0x7f | 0x80)
		s.curBuf[s.lastWritten-1] = byte(v >> 14)
		s.lastWritten -= 3

	case v < 1<<28:
		s.curBuf[s.lastWritten-4] = byte((v>>0)&0x7f | 0x80)
		s.curBuf[s.lastWritten-3] = byte((v>>7)&0x7f | 0x80)
		s.curBuf[s.lastWritten-2] = byte((v>>14)&0x7f | 0x80)
		s.curBuf[s.lastWritten-1] = byte(v >> 21)
		s.lastWritten -= 4

	case v < 1<<35:
		s.curBuf[s.lastWritten-5] = byte((v>>0)&0x7f | 0x80)
		s.curBuf[s.lastWritten-4] = byte((v>>7)&0x7f | 0x80)
		s.curBuf[s.lastWritten-3] = byte((v>>14)&0x7f | 0x80)
		s.curBuf[s.lastWritten-2] = byte((v>>21)&0x7f | 0x80)
		s.curBuf[s.lastWritten-1] = byte(v >> 28)
		s.lastWritten -= 5

	case v < 1<<42:
		s.curBuf[s.lastWritten-6] = byte((v>>0)&0x7f | 0x80)
		s.curBuf[s.lastWritten-5] = byte((v>>7)&0x7f | 0x80)
		s.curBuf[s.lastWritten-4] = byte((v>>14)&0x7f | 0x80)
		s.curBuf[s.lastWritten-3] = byte((v>>21)&0x7f | 0x80)
		s.curBuf[s.lastWritten-2] = byte((v>>28)&0x7f | 0x80)
		s.curBuf[s.lastWritten-1] = byte(v >> 35)
		s.lastWritten -= 6

	case v < 1<<49:
		s.curBuf[s.lastWritten-7] = byte((v>>0)&0x7f | 0x80)
		s.curBuf[s.lastWritten-6] = byte((v>>7)&0x7f | 0x80)
		s.curBuf[s.lastWritten-5] = byte((v>>14)&0x7f | 0x80)
		s.curBuf[s.lastWritten-4] = byte((v>>21)&0x7f | 0x80)
		s.curBuf[s.lastWritten-3] = byte((v>>28)&0x7f | 0x80)
		s.curBuf[s.lastWritten-2] = byte((v>>35)&0x7f | 0x80)
		s.curBuf[s.lastWritten-1] = byte(v >> 42)
		s.lastWritten -= 7

	case v < 1<<56:
		s.curBuf[s.lastWritten-8] = byte((v>>0)&0x7f | 0x80)
		s.curBuf[s.lastWritten-7] = byte((v>>7)&0x7f | 0x80)
		s.curBuf[s.lastWritten-6] = byte((v>>14)&0x7f | 0x80)
		s.curBuf[s.lastWritten-5] = byte((v>>21)&0x7f | 0x80)
		s.curBuf[s.lastWritten-4] = byte((v>>28)&0x7f | 0x80)
		s.curBuf[s.lastWritten-3] = byte((v>>35)&0x7f | 0x80)
		s.curBuf[s.lastWritten-2] = byte((v>>42)&0x7f | 0x80)
		s.curBuf[s.lastWritten-1] = byte(v >> 49)
		s.lastWritten -= 8

	case v < 1<<63:
		s.curBuf[s.lastWritten-9] = byte((v>>0)&0x7f | 0x80)
		s.curBuf[s.lastWritten-8] = byte((v>>7)&0x7f | 0x80)
		s.curBuf[s.lastWritten-7] = byte((v>>14)&0x7f | 0x80)
		s.curBuf[s.lastWritten-6] = byte((v>>21)&0x7f | 0x80)
		s.curBuf[s.lastWritten-5] = byte((v>>28)&0x7f | 0x80)
		s.curBuf[s.lastWritten-4] = byte((v>>35)&0x7f | 0x80)
		s.curBuf[s.lastWritten-3] = byte((v>>42)&0x7f | 0x80)
		s.curBuf[s.lastWritten-2] = byte((v>>49)&0x7f | 0x80)
		s.curBuf[s.lastWritten-1] = byte(v >> 56)
		s.lastWritten -= 9

	default:
		s.curBuf[s.lastWritten-10] = byte((v>>0)&0x7f | 0x80)
		s.curBuf[s.lastWritten-9] = byte((v>>7)&0x7f | 0x80)
		s.curBuf[s.lastWritten-8] = byte((v>>14)&0x7f | 0x80)
		s.curBuf[s.lastWritten-7] = byte((v>>21)&0x7f | 0x80)
		s.curBuf[s.lastWritten-6] = byte((v>>28)&0x7f | 0x80)
		s.curBuf[s.lastWritten-5] = byte((v>>35)&0x7f | 0x80)
		s.curBuf[s.lastWritten-4] = byte((v>>42)&0x7f | 0x80)
		s.curBuf[s.lastWritten-3] = byte((v>>49)&0x7f | 0x80)
		s.curBuf[s.lastWritten-2] = byte((v>>56)&0x7f | 0x80)
		s.curBuf[s.lastWritten-1] = 1
		s.lastWritten -= 10
	}
}
