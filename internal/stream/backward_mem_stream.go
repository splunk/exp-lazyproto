package stream

import "io"

const bufSize = 4096

type BackwardMemStream struct {
	curBuf      []byte
	lastWritten int

	readyBufs [][]byte
}

func NewBackwardMemStream() *BackwardMemStream {
	s := &BackwardMemStream{}
	s.curBuf = make([]byte, bufSize)
	s.lastWritten = bufSize
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

func (s *BackwardMemStream) Byte(b byte) {
	if s.lastWritten == 0 {
		s.allocNewBuf()
	}
	s.lastWritten--
	s.curBuf[s.lastWritten] = b
}

func (s *BackwardMemStream) Bytes(b []byte) {
	if len(b) >= bufSize {
		// Lots of bytes, just make it a ready buffer itself.
		// First we need to push current accumulated buffer:
		s.readyBufs = append(s.readyBufs, s.curBuf[s.lastWritten:])

		if s.lastWritten > 0 {
			// Free space remains in the current buffer. Make it available
			// for future writes.
			s.curBuf = s.curBuf[:s.lastWritten:s.lastWritten]
		} else {
			// No more free same in the current buffer, just allocate a new one.
			s.curBuf = make([]byte, bufSize)
			s.lastWritten = bufSize
		}
		// Now add the new bytes.
		s.readyBufs = append(s.readyBufs, b)
		return
	}

	bytesFit := s.lastWritten
	if bytesFit > len(b) {
		bytesFit = len(b)
	}
	s.lastWritten -= bytesFit
	copy(s.curBuf[s.lastWritten:], b[len(b)-bytesFit:])

	bytesRemaining := len(b) - bytesFit

	if s.lastWritten == 0 {
		s.allocNewBuf()
	}

	if bytesRemaining > 0 {
		s.lastWritten -= bytesRemaining
		copy(s.curBuf[s.lastWritten:], b[:bytesRemaining])
	}
}

func (s *BackwardMemStream) allocNewBuf() {
	s.readyBufs = append(s.readyBufs, s.curBuf[s.lastWritten:])
	s.curBuf = make([]byte, bufSize)
	s.lastWritten = bufSize
}
