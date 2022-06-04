package stream

import (
	"io"
)

const firstBufSize = 4096

type BackwardMemStream struct {
	curBuf      []byte
	lastWritten int
	nextBufSize int

	readyBufs    [][]byte
	readyBufsLen int
}

func NewBackwardMemStream() *BackwardMemStream {
	s := &BackwardMemStream{}
	s.curBuf = make([]byte, firstBufSize)
	s.lastWritten = firstBufSize
	s.nextBufSize = firstBufSize * 2
	return s
}

func (s *BackwardMemStream) WriteTo(out io.Writer) error {
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

func (s *BackwardMemStream) Len() int {
	return s.readyBufsLen + len(s.curBuf) - s.lastWritten
}

type PreparedFixed64Key byte

// Fixed64Prepared writes a value of fixed64 to the stream.
func (s *BackwardMemStream) Fixed64Prepared(fieldKey PreparedFixed64Key, value uint64) {
	if value == 0 {
		return
	}

	// Write in backward order.
	s.writeFixed64(value)
	s.writeByte(byte(fieldKey))
}

// PreparedKey of string,bytes or embedded wire type.
type PreparedKey byte

// StringPrepared writes a string to the stream.
func (s *BackwardMemStream) StringPrepared(key PreparedKey, value string) {
	vlen := len(value)
	if vlen == 0 {
		return
	}

	// Write in backward order.
	s.writeBytes([]byte(value))
	s.writeVarint(uint64(vlen))
	s.writeByte(byte(key))
}

type EmbeddedToken int

func (s *BackwardMemStream) BeginEmbedded() EmbeddedToken {
	return EmbeddedToken(s.Len())
}

func (s *BackwardMemStream) EndEmbedded(beginToken EmbeddedToken, fieldKey PreparedKey) {
	embeddedSize := s.Len() - int(beginToken)

	// Write in backward order.
	s.writeVarint(uint64(embeddedSize))
	s.writeByte(byte(fieldKey))
}

func (s *BackwardMemStream) writeByte(b byte) {
	if s.lastWritten == 0 {
		s.allocNewBuf()
	}
	s.lastWritten--
	s.curBuf[s.lastWritten] = b
}

func (s *BackwardMemStream) appendReadyBuf(b []byte) {
	s.readyBufs = append(s.readyBufs, b)
	s.readyBufsLen += len(b)
}

func (s *BackwardMemStream) writeBytes(b []byte) {
	if len(b) >= firstBufSize {
		// Lots of bytes, just make it a ready buffer itself.
		// First we need to push current accumulated buffer:
		s.appendReadyBuf(s.curBuf[s.lastWritten:])

		if s.lastWritten > 0 {
			// Free space remains in the current buffer. Make it available
			// for future writes.
			s.curBuf = s.curBuf[:s.lastWritten:s.lastWritten]
		} else {
			// No more free same in the current buffer, just allocate a new one.
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

func (s *BackwardMemStream) allocNewBuf() {
	s.appendReadyBuf(s.curBuf[s.lastWritten:])

	// Double the buf size on every new allocation to have constant amortized cost.
	s.curBuf = make([]byte, s.nextBufSize)
	s.lastWritten = s.nextBufSize
	s.nextBufSize *= 2
}

func (s *BackwardMemStream) writeFixed64(v uint64) {
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

func (s *BackwardMemStream) writeVarint(v uint64) {
	if s.lastWritten < 10 {
		s.allocNewBuf()
	}

	switch {
	case v < 1<<7:
		s.curBuf[s.lastWritten-1] = byte(v)
		s.lastWritten -= 1

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
