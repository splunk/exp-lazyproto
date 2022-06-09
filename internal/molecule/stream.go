package molecule

import (
	"math"

	"github.com/tigrannajaryan/exp-lazyproto/internal/molecule/src/protowire"
)

const (
	// The defaultBufferSize is the default size for buffers used for embedded
	// values, which must first be written to a buffer to determine their
	// length.  This is not used if BufferFactory is set.
	defaultBufferSize int = 1024 * 64
)

// A ProtoStream supports writing protobuf data in a streaming fashion.  Its methods
// will write their output to the wrapped `io.Writer`.  Zero values are not included.
//
// ProtoStream instances are *not* threadsafe and *not* re-entrant.
type ProtoStream struct {
	// The outputBuffer is the writer to which the protobuf-encoded bytes are
	// written.
	outputBuffer []byte

	// The scratchBuffer is a buffer used and re-used for generating output.
	// Each method should begin by resetting this buffer.
	scratchBuffer []byte

	// The scratchArray is a second, very small array used for packed
	// encodings.  It is large enough to fit two max-size varints (10 bytes
	// each) without reallocation
	scratchArray [20]byte
}

// NewProtoStream creates a new ProtoStream writing to the given Writer.  If the
// writer is nil, the stream cannot be used until it has been set with `Reset`.
func NewProtoStream() *ProtoStream {
	ps := &ProtoStream{
		scratchBuffer: make([]byte, 0, defaultBufferSize),
		outputBuffer:  make([]byte, 0, defaultBufferSize),
	}
	return ps
}

func (ps *ProtoStream) BufferBytes() ([]byte, error) {
	return ps.outputBuffer, nil
}

// Reset sets the Writer to which this ProtoStream streams.  If the writer is nil,
// then the protostream cannot be used until Reset is called with a non-nil value.
func (ps *ProtoStream) Reset() {
	ps.outputBuffer = ps.outputBuffer[:0]
	ps.scratchBuffer = ps.scratchBuffer[:0]
}

// Double writes a value of proto type double to the stream.
func (ps *ProtoStream) Double(fieldNumber int, value float64) error {
	if value == 0.0 {
		return nil
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.encodeKeyToScratch(fieldNumber, protowire.Fixed64Type)
	ps.scratchBuffer = protowire.AppendFixed64(ps.scratchBuffer, math.Float64bits(value))
	ps.writeScratch()
	return nil
}

func (ps *ProtoStream) DoublePrepared(fieldKey PreparedKey, value float64) {
	if value == 0 {
		return
	}

	ps.outputBuffer = append(ps.outputBuffer, byte(fieldKey))
	ps.outputBuffer = protowire.AppendFixed64(ps.outputBuffer, math.Float64bits(value))
}

// DoublePacked writes a slice of values of proto type double to the stream,
// in packed form.
func (ps *ProtoStream) DoublePacked(fieldNumber int, values []float64) {
	if len(values) == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	for _, value := range values {
		ps.scratchBuffer = protowire.AppendFixed64(
			ps.scratchBuffer, math.Float64bits(value),
		)
	}

	ps.writeScratchAsPacked(fieldNumber)
}

// Float writes a value of proto type double to the stream.
func (ps *ProtoStream) Float(fieldNumber int, value float32) error {
	if value == 0.0 {
		return nil
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.encodeKeyToScratch(fieldNumber, protowire.Fixed32Type)
	ps.scratchBuffer = protowire.AppendFixed32(ps.scratchBuffer, math.Float32bits(value))
	ps.writeScratch()
	return nil
}

// FloatPacked writes a slice of values of proto type float to the stream,
// in packed form.
func (ps *ProtoStream) FloatPacked(fieldNumber int, values []float32) {
	if len(values) == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	for _, value := range values {
		ps.scratchBuffer = protowire.AppendFixed32(
			ps.scratchBuffer, math.Float32bits(value),
		)
	}
	ps.writeScratchAsPacked(fieldNumber)
}

// Int32 writes a value of proto type int32 to the stream.
func (ps *ProtoStream) Int32(fieldNumber int, value int32) {
	if value == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.encodeKeyToScratch(fieldNumber, protowire.VarintType)
	ps.scratchBuffer = protowire.AppendVarint(ps.scratchBuffer, uint64(value))
	ps.writeScratch()
	return
}

// Int32Packed writes a slice of values of proto type int32 to the stream,
// in packed form.
func (ps *ProtoStream) Int32Packed(fieldNumber int, values []int32) {
	if len(values) == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	for _, value := range values {
		ps.scratchBuffer = protowire.AppendVarint(ps.scratchBuffer, uint64(value))
	}
	ps.writeScratchAsPacked(fieldNumber)
}

// Int64 writes a value of proto type int64 to the stream.
func (ps *ProtoStream) Int64(fieldNumber int, value int64) {
	if value == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.encodeKeyToScratch(fieldNumber, protowire.VarintType)
	ps.scratchBuffer = protowire.AppendVarint(ps.scratchBuffer, uint64(value))
	ps.writeScratch()
	return
}

func (ps *ProtoStream) Int64Prepared(fieldKey PreparedKey, value int64) {
	if value == 0 {
		return
	}

	ps.outputBuffer = append(ps.outputBuffer, byte(fieldKey))
	ps.outputBuffer = protowire.AppendVarint(ps.outputBuffer, uint64(value))
}

// Int64Packed writes a slice of values of proto type int64 to the stream,
// in packed form.
func (ps *ProtoStream) Int64Packed(fieldNumber int, values []int64) {
	if len(values) == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	for _, value := range values {
		ps.scratchBuffer = protowire.AppendVarint(ps.scratchBuffer, uint64(value))
	}
	ps.writeScratchAsPacked(fieldNumber)
}

// Uint32 writes a value of proto type uint32 to the stream.
func (ps *ProtoStream) Uint32(fieldNumber int, value uint32) {
	if value == 0 {
		return
	}
	ps.encodeKeyToOutput(fieldNumber, protowire.VarintType)
	ps.outputBuffer = protowire.AppendVarint(ps.outputBuffer, uint64(value))
	return
}

// Uint32 writes a value of proto type uint32 to the stream.
func (ps *ProtoStream) Uint32Prepared(fieldKey PreparedKey, value uint32) {
	if value == 0 {
		return
	}
	ps.outputBuffer = append(ps.outputBuffer, byte(fieldKey))
	ps.outputBuffer = protowire.AppendVarint(ps.outputBuffer, uint64(value))
	return
}

//// Uint32 writes a value of proto type uint32 to the stream.
//func (ps *ProtoStream) Uint32LongPrepared(fieldKey PreparedLongKey, value uint32) {
//	if value == 0 {
//		return
//	}
//	ps.outputBuffer = append(ps.outputBuffer, fieldKey...)
//	ps.outputBuffer = protowire.AppendVarint(ps.outputBuffer, uint64(value))
//	return
//}

// Uint32Packed writes a slice of values of proto type uint32 to the stream,
// in packed form.
func (ps *ProtoStream) Uint32Packed(fieldNumber int, values []uint32) {
	if len(values) == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	for _, value := range values {
		ps.scratchBuffer = protowire.AppendVarint(ps.scratchBuffer, uint64(value))
	}
	ps.writeScratchAsPacked(fieldNumber)
}

// Uint64 writes a value of proto type uint64 to the stream.
func (ps *ProtoStream) Uint64(fieldNumber int, value uint64) {
	if value == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.encodeKeyToScratch(fieldNumber, protowire.VarintType)
	ps.scratchBuffer = protowire.AppendVarint(ps.scratchBuffer, value)
	ps.writeScratch()
}

func (ps *ProtoStream) Uint64Prepared(fieldKey PreparedKey, value uint64) {
	if value == 0 {
		return
	}
	ps.outputBuffer = append(ps.outputBuffer, byte(fieldKey))
	ps.outputBuffer = protowire.AppendVarint(ps.outputBuffer, value)
}

// Uint64Packed writes a slice of values of proto type uint64 to the stream,
// in packed form.
func (ps *ProtoStream) Uint64Packed(fieldNumber int, values []uint64) {
	if len(values) == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	for _, value := range values {
		ps.scratchBuffer = protowire.AppendVarint(ps.scratchBuffer, value)
	}
	ps.writeScratchAsPacked(fieldNumber)
}

// Sint32 writes a value of proto type sint32 to the stream.
func (ps *ProtoStream) Sint32(fieldNumber int, value int32) {
	if value == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.encodeKeyToScratch(fieldNumber, protowire.VarintType)
	ps.scratchBuffer = protowire.AppendVarint(
		ps.scratchBuffer, protowire.EncodeZigZag(int64(value)),
	)
	ps.writeScratch()
}

func (ps *ProtoStream) Sint32Prepared(fieldKey PreparedKey, value int32) {
	if value == 0 {
		return
	}
	ps.outputBuffer = append(ps.outputBuffer, byte(fieldKey))
	ps.outputBuffer = protowire.AppendVarint(
		ps.outputBuffer, protowire.EncodeZigZag(int64(value)),
	)
}

// Sint32Packed writes a slice of values of proto type sint32 to the stream,
// in packed form.
func (ps *ProtoStream) Sint32Packed(fieldNumber int, values []int32) {
	if len(values) == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	for _, value := range values {
		ps.scratchBuffer = protowire.AppendVarint(
			ps.scratchBuffer, protowire.EncodeZigZag(int64(value)),
		)
	}
	ps.writeScratchAsPacked(fieldNumber)
}

// Sint64 writes a value of proto type sint64 to the stream.
func (ps *ProtoStream) Sint64(fieldNumber int, value int64) {
	if value == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.encodeKeyToScratch(fieldNumber, protowire.VarintType)
	ps.scratchBuffer = protowire.AppendVarint(ps.scratchBuffer, zigzag64(uint64(value)))
	ps.writeScratch()
}

// Sint64Packed writes a slice of values of proto type sint64 to the stream,
// in packed form.
func (ps *ProtoStream) Sint64Packed(fieldNumber int, values []int64) {
	if len(values) == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	for _, value := range values {
		ps.scratchBuffer = protowire.AppendVarint(
			ps.scratchBuffer, zigzag64(uint64(value)),
		)
	}
	ps.writeScratchAsPacked(fieldNumber)
}

// Fixed32 writes a value of proto type fixed32 to the stream.
func (ps *ProtoStream) Fixed32(fieldNumber int, value uint32) {
	if value == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.encodeKeyToScratch(fieldNumber, protowire.Fixed32Type)
	ps.scratchBuffer = protowire.AppendFixed32(ps.scratchBuffer, value)
	ps.writeScratch()
}

func (ps *ProtoStream) Fixed32Prepared(fieldKey PreparedKey, value uint32) {
	if value == 0 {
		return
	}
	ps.outputBuffer = append(ps.outputBuffer, byte(fieldKey))
	ps.outputBuffer = protowire.AppendFixed32(ps.outputBuffer, value)
}

// Fixed32Packed writes a slice of values of proto type fixed32 to the stream,
// in packed form.
func (ps *ProtoStream) Fixed32Packed(fieldNumber int, values []uint32) {
	if len(values) == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	for _, value := range values {
		ps.scratchBuffer = protowire.AppendFixed32(ps.scratchBuffer, value)
	}
	ps.writeScratchAsPacked(fieldNumber)
}

// Fixed64 writes a value of proto type fixed64 to the stream.
func (ps *ProtoStream) Fixed64(fieldNumber int, value uint64) {
	if value == 0 {
		return
	}
	ps.encodeKeyToOutput(fieldNumber, protowire.Fixed64Type)
	ps.outputBuffer = protowire.AppendFixed64(ps.outputBuffer, value)
}

// Fixed64 writes a value of proto type fixed64 to the stream.
func (ps *ProtoStream) Fixed64Prepared(fieldKey PreparedKey, value uint64) {
	if value == 0 {
		return
	}
	ps.outputBuffer = append(ps.outputBuffer, byte(fieldKey))
	ps.outputBuffer = protowire.AppendFixed64(ps.outputBuffer, value)
}

// Fixed64Packed writes a slice of values of proto type fixed64 to the stream,
// in packed form.
func (ps *ProtoStream) Fixed64Packed(fieldNumber int, values []uint64) {
	if len(values) == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	for _, value := range values {
		ps.scratchBuffer = protowire.AppendFixed64(ps.scratchBuffer, value)
	}
	ps.writeScratchAsPacked(fieldNumber)
}

// Sfixed32 writes a value of proto type sfixed32 to the stream.
func (ps *ProtoStream) Sfixed32(fieldNumber int, value int32) {
	if value == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.encodeKeyToScratch(fieldNumber, protowire.Fixed32Type)
	ps.scratchBuffer = protowire.AppendFixed32(ps.scratchBuffer, uint32(value))
	ps.writeScratch()
}

// Sfixed32Packed writes a slice of values of proto type sfixed32 to the stream,
// in packed form.
func (ps *ProtoStream) Sfixed32Packed(fieldNumber int, values []int32) {
	if len(values) == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	for _, value := range values {
		ps.scratchBuffer = protowire.AppendFixed32(ps.scratchBuffer, uint32(value))
	}
	ps.writeScratchAsPacked(fieldNumber)
}

// Sfixed64 writes a value of proto type sfixed64 to the stream.
func (ps *ProtoStream) Sfixed64(fieldNumber int, value int64) {
	if value == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.encodeKeyToScratch(fieldNumber, protowire.Fixed64Type)
	ps.scratchBuffer = protowire.AppendFixed64(ps.scratchBuffer, uint64(value))
	ps.writeScratch()
}

// Sfixed64 writes a value of proto type sfixed64 to the stream.
func (ps *ProtoStream) SFixed64Prepared(fieldKey PreparedKey, value int64) {
	if value == 0 {
		return
	}

	ps.outputBuffer = append(ps.outputBuffer, byte(fieldKey))
	ps.outputBuffer = protowire.AppendFixed64(ps.outputBuffer, uint64(value))
}

// Sfixed64Packed writes a slice of values of proto type sfixed64 to the stream,
// in packed form.
func (ps *ProtoStream) Sfixed64Packed(fieldNumber int, values []int64) {
	if len(values) == 0 {
		return
	}
	ps.scratchBuffer = ps.scratchBuffer[:0]
	for _, value := range values {
		ps.scratchBuffer = protowire.AppendFixed64(ps.scratchBuffer, uint64(value))
	}
	ps.writeScratchAsPacked(fieldNumber)
}

// Bool writes a value of proto type bool to the stream.
func (ps *ProtoStream) Bool(fieldNumber int, value bool) {
	if value == false {
		return
	}
	ps.encodeKeyToOutput(fieldNumber, protowire.VarintType)
	var bit uint64
	if value {
		bit = 1
	}
	ps.outputBuffer = protowire.AppendVarint(ps.outputBuffer, bit)
}

func (ps *ProtoStream) BoolPrepared(fieldKey PreparedKey, value bool) {
	if value == false {
		return
	}

	var bit uint64
	if value {
		bit = 1
	}

	ps.outputBuffer = append(ps.outputBuffer, byte(fieldKey))
	ps.outputBuffer = protowire.AppendVarint(ps.outputBuffer, bit)
}

// String writes a string to the stream.
func (ps *ProtoStream) String(fieldNumber int, value string) {
	if len(value) == 0 {
		return
	}
	ps.encodeKeyToOutput(fieldNumber, protowire.BytesType)
	ps.outputBuffer = protowire.AppendVarint(ps.outputBuffer, uint64(len(value)))

	ps.outputBuffer = append(ps.outputBuffer, value...)
}

// PreparedKey of string,bytes or embedded wire type.
type PreparedKey byte

//type PreparedKey byte
//type PreparedKey byte
//type PreparedKey byte
//type PreparedKey byte
//type PreparedKey byte
//type PreparedKey byte

//type PreparedLongKey []byte

func PrepareField(fieldNumber int, wireType protowire.Type) PreparedKey {
	b := protowire.AppendVarint([]byte{}, uint64(fieldNumber)<<3+uint64(wireType))
	if len(b) > 1 {
		panic("fieldNumber too high, max field number 15 is supported by PrepareField")
	}
	return PreparedKey(b[0])
}

//func PrepareLongField(fieldNumber int, wireType protowire.Type) PreparedLongKey {
//	b := protowire.AppendVarint([]byte{}, uint64(fieldNumber)<<3+uint64(wireType))
//	return PreparedLongKey(b)
//}

func PrepareStringField(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.BytesType)
}

func PrepareBytesField(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.BytesType)
}

func PrepareEmbeddedField(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.BytesType)
}

func PrepareBoolField(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.VarintType)
}

func PrepareUint32Field(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.VarintType)
}

func PrepareSint32Field(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.VarintType)
}

func PrepareFixed64Field(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.Fixed64Type)
}

func PrepareInt64Field(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.VarintType)
}

func PrepareUint64Field(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.VarintType)
}

func PrepareFixed32Field(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.Fixed32Type)
}

func PrepareDoubleField(fieldNumber int) PreparedKey {
	return PrepareField(fieldNumber, protowire.Fixed64Type)
}

// String writes a string to the stream.
func (ps *ProtoStream) StringPrepared(key PreparedKey, value string) {
	vlen := len(value)
	if vlen == 0 {
		return
	}

	if vlen < 1<<7 {
		ps.outputBuffer = append(ps.outputBuffer, byte(key), byte(vlen))
	} else {
		ps.outputBuffer = append(ps.outputBuffer, byte(key))
		ps.outputBuffer = protowire.AppendVarint(ps.outputBuffer, uint64(vlen))
	}
	ps.outputBuffer = append(ps.outputBuffer, value...)
}

// Bytes writes the given bytes to the stream.
func (ps *ProtoStream) Bytes(fieldNumber int, value []byte) {
	if len(value) == 0 {
		return
	}
	ps.encodeKeyToOutput(fieldNumber, protowire.BytesType)
	ps.outputBuffer = protowire.AppendVarint(ps.outputBuffer, uint64(len(value)))

	ps.writeAll(value)

}

// String writes a string to the stream.
func (ps *ProtoStream) BytesPrepared(key PreparedKey, value []byte) {
	vlen := len(value)
	if vlen == 0 {
		return
	}

	if vlen < 1<<7 {
		ps.outputBuffer = append(ps.outputBuffer, byte(key), byte(vlen))
	} else {
		ps.outputBuffer = append(ps.outputBuffer, byte(key))
		ps.outputBuffer = protowire.AppendVarint(ps.outputBuffer, uint64(vlen))
	}
	ps.outputBuffer = append(ps.outputBuffer, value...)
}

func (ps *ProtoStream) Raw(value []byte) {
	if len(value) == 0 {
		return
	}
	ps.writeAll(value)
}

type EmbeddedToken int

func (ps *ProtoStream) ReserveCapacity(capacity int) {
	if cap(ps.outputBuffer) < capacity {
		b := make([]byte, len(ps.outputBuffer), capacity)
		copy(b, ps.outputBuffer)
		ps.outputBuffer = b
	}
}

func (ps *ProtoStream) BeginEmbedded() EmbeddedToken {
	curPos := len(ps.outputBuffer)

	// Reserve 2 bytes since it is the minimal space we need to encode the key
	// (and likely also most common space size, which is good for performance).
	ps.outputBuffer = append(ps.outputBuffer, 0, 0)

	return EmbeddedToken(curPos)
}

//func AppendVarint(b []byte, v uint64) []byte {
//	for v >= 1<<7 {
//		b = append(b, uint8(v&0x7f|0x80))
//		v >>= 7
//		//offset++
//	}
//	b = append(b, uint8(v))
//	return b
//}

func (ps *ProtoStream) EndEmbedded(startPos EmbeddedToken, fieldNumber int) {
	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.encodeKeyToScratch(fieldNumber, protowire.BytesType)

	childSize := len(ps.outputBuffer) - int(startPos) - 2

	ps.scratchBuffer = protowire.AppendVarint(
		ps.scratchBuffer, uint64(childSize),
	)

	if len(ps.scratchBuffer) == 2 {
		// Reserved space was just right.
		ps.outputBuffer[startPos] = ps.scratchBuffer[0]
		ps.outputBuffer[startPos+1] = ps.scratchBuffer[1]
	} else {
		ps.outputBuffer = append(
			ps.outputBuffer[:startPos],
			append(ps.scratchBuffer, ps.outputBuffer[startPos+2:]...)...,
		)
	}

}

func (ps *ProtoStream) EndEmbeddedPrepared(
	startPos EmbeddedToken, fieldKey PreparedKey,
) {
	childSize := len(ps.outputBuffer) - int(startPos) - 2
	if uint64(childSize) < 1<<7 {
		// Reserved space was just right.
		ps.outputBuffer[startPos] = byte(fieldKey)
		ps.outputBuffer[startPos+1] = byte(childSize)
		return
	}

	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.scratchBuffer = append(ps.scratchBuffer, byte(fieldKey))
	ps.scratchBuffer = protowire.AppendVarint(ps.scratchBuffer, uint64(childSize))

	ps.outputBuffer = append(
		ps.outputBuffer[:startPos],
		append(ps.scratchBuffer, ps.outputBuffer[startPos+2:]...)...,
	)
}

// Embedded is used for constructing embedded messages.  It calls the given
// function with a new ProtoStream, then embeds the result in the current
// stream.
//
// NOTE: if the inner function creates an empty message (such as for a struct
// at its zero value), that empty message will still be added to the stream.
func (ps *ProtoStream) Embedded(fieldNumber int, inner func() error) error {
	curPos := len(ps.outputBuffer)

	// Reserve 2 bytes since it is the minimal space we need to encode the key
	// (and likely also most common space size, which is good for performance).
	ps.outputBuffer = append(ps.outputBuffer, 0, 0)

	err := inner()
	if err != nil {
		return err
	}

	ps.scratchBuffer = ps.scratchBuffer[:0]
	ps.encodeKeyToScratch(fieldNumber, protowire.BytesType)

	childSize := len(ps.outputBuffer) - curPos - 2

	ps.scratchBuffer = protowire.AppendVarint(
		ps.scratchBuffer, uint64(childSize),
	)

	if len(ps.scratchBuffer) == 2 {
		// Reserved space was just right.
		copy(ps.outputBuffer[curPos:], ps.scratchBuffer)
	} else {
		ps.outputBuffer = append(
			ps.outputBuffer[:curPos],
			append(ps.scratchBuffer, ps.outputBuffer[curPos+2:]...)...,
		)
	}

	return nil
}

// writeScratch flushes the scratch buffer to output.
func (ps *ProtoStream) writeScratch() {
	ps.writeAll(ps.scratchBuffer)
}

// writeScratchAsPacked writes the scratch buffer to outputBuffer, prefixed with
// the given key and the length of the scratch buffer.  This is used for packed
// encodings.
func (ps *ProtoStream) writeScratchAsPacked(fieldNumber int) {
	// The scratch buffer is full of the packed data, but we need to write
	// the key and size, so we use scratchArray.  We could use a stack allocation
	// here, but as of writing the go compiler is not smart enough to figure out
	// that the value does not escape.
	keysize := ps.scratchArray[:0]
	keysize = protowire.AppendVarint(
		keysize, uint64((fieldNumber<<3)|int(protowire.BytesType)),
	)
	keysize = protowire.AppendVarint(keysize, uint64(len(ps.scratchBuffer)))

	// Write the key and length prefix.
	ps.writeAll(keysize)

	// Write out the embedded message.
	ps.writeScratch()
}

// writeAll writes an entire buffer to output.
func (ps *ProtoStream) writeAll(buf []byte) {
	ps.outputBuffer = append(ps.outputBuffer, buf...)
}

// encodeKeyToScratch encodes a protobuf key into ps.scratch.
func (ps *ProtoStream) encodeKeyToScratch(fieldNumber int, wireType protowire.Type) {
	ps.scratchBuffer = protowire.AppendVarint(
		ps.scratchBuffer, uint64(fieldNumber)<<3+uint64(wireType),
	)
}
func (ps *ProtoStream) encodeKeyToOutput(fieldNumber int, wireType protowire.Type) {
	ps.outputBuffer = protowire.AppendVarint(
		ps.outputBuffer, uint64(fieldNumber)<<3+uint64(wireType),
	)
}

func zigzag32(v uint64) uint64 {
	return uint64((uint32(v) << 1) ^ uint32((int32(v) >> 31)))
}

func zigzag64(v uint64) uint64 {
	return (v << 1) ^ uint64((int64(v) >> 63))
}
