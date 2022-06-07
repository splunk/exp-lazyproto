package sizedstream

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"
)

func TestProtoStreamByte(t *testing.T) {
	// Check some interesting byte counts that are likely to hit edge cases.
	byteCounts := []int{0, 1, 2, 3, 4094, 4095, 4096, 4097, 8191, 8192, 8193, 100000}

	for _, byteCount := range byteCounts {
		t.Run(
			strconv.Itoa(byteCount), func(t *testing.T) {
				s := NewProtoStream()
				for i := 0; i < byteCount; i++ {
					s.writeByte(byte(i % 256))
				}

				destBytes, err := s.BufferBytes()
				assert.NoError(t, err)
				assert.Len(t, destBytes, byteCount)
				for i := 0; i < byteCount; i++ {
					assert.Equal(t, byte(i%256), destBytes[i])
				}
			},
		)
	}
}

func TestProtoStreamBytes(t *testing.T) {
	// Check some interesting byte counts that are likely to hit edge cases.
	byteCounts := []int{0, 1, 2, 3, 4094, 4095, 4096, 4097, 8191, 8192, 8193, 100000}

	for _, byteCount := range byteCounts {
		t.Run(
			strconv.Itoa(byteCount), func(t *testing.T) {

				s := NewProtoStream()
				for i := 0; i < byteCount; {
					maxChunkSize := byteCount + 1
					if maxChunkSize > 10000 {
						// We don't want chunks too large. Use smaller chunks to
						// trigger more edge cases.
						maxChunkSize = 10000
					}
					sz := rand.Int() % maxChunkSize
					if sz > (byteCount - i) {
						sz = byteCount - i
					}

					b := make([]byte, sz)
					for j := 0; j < len(b); j++ {
						b[j] = byte(i % 256)
						i++
					}
					s.Raw(b)
				}

				destBytes, err := s.BufferBytes()
				assert.NoError(t, err)
				require.Len(
					t, destBytes, byteCount, "Incorrect length of resulting bytes",
				)
				for i := 0; i < byteCount; i++ {
					require.Equal(t, byte(i%256), destBytes[i])
				}
			},
		)
	}
}

func TestProtoStreamVarint(t *testing.T) {
	vals := []uint64{
		0, 1, 2,
		127, 128, 129,
		1<<14 - 1, 1 << 14, 1<<14 + 1,
		1<<21 - 1, 1 << 21, 1<<21 + 1,
		1<<28 - 1, 1 << 28, 1<<28 + 1,
		1<<35 - 1, 1 << 35, 1<<35 + 1,
		1<<42 - 1, 1 << 42, 1<<42 + 1,
		1<<49 - 1, 1 << 49, 1<<49 + 1,
		1<<56 - 1, 1 << 56, 1<<56 + 1,
		1<<63 - 1, 1 << 63, 1<<63 + 1,
	}

	for _, val := range vals {
		t.Run(
			strconv.FormatUint(val, 10), func(t *testing.T) {
				s := NewProtoStream()
				s.writeVarint(val)

				destBytes, err := s.BufferBytes()
				assert.NoError(t, err)

				var checkBuf []byte
				checkBuf = protowire.AppendVarint(checkBuf, val)

				assert.EqualValues(t, checkBuf, destBytes)
			},
		)
	}
}

//const varintTestRange = 20000

//func BenchmarkProtoStreamVarint(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		s := NewProtoStream()
//		for v := uint64(0); v < varintTestRange; v++ {
//			s.writeVarint(v)
//		}
//	}
//}
//
//func BenchmarkGoogleProtobufVarint(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		var s []byte
//		for v := uint64(0); v < varintTestRange; v++ {
//			s = protowire.AppendVarint(s, v)
//		}
//	}
//}
