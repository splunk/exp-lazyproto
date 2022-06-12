package generator

import (
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	_ "github.com/jhump/protoreflect/desc/protoparse"
	"github.com/tigrannajaryan/exp-lazyproto/internal/molecule/src/codec"
)

func (g *generator) oUnmarshalFunc() error {
	g.o(
		`
func Unmarshal$MessageName(bytes []byte, opts lazyproto.UnmarshalOpts) (*$MessageName, error) {
	if opts.WithValidate {
		if err := validate$MessageName(bytes); err != nil {
			return nil, err		
		}
	}

	m := $messagePool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}
`,
	)
	return g.lastErr
}

func (g *generator) oDecodeMethod() error {

	g.o(
		`
func (m *$MessageName) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))
`,
	)

	g.i(1)

	if g.msg.FlagsBitCount > 0 {
		g.o(`// Reset all "decoded" and "presence" flags.`)
		g.o(`m._flags = 0`)
		g.o(``)
	}

	g.oCalcRepeatedFieldCounts()

	g.oMsgDecodeLoop(decodeFull)

	g.i(-1)

	g.o(
		`
	return nil
}
`,
	)

	return g.lastErr
}

func (g *generator) oValidateFunc() error {

	g.o(
		`
func validate$MessageName(b []byte) error {
	buf := codec.NewBuffer(b)
`,
	)

	g.i(1)
	g.oMsgDecodeLoop(decodeValidate)
	g.i(-1)

	g.o(
		`
	return nil
}
	`,
	)

	return g.lastErr
}

var protoTypeToWireType = map[descriptor.FieldDescriptorProto_Type]codec.WireType{
	descriptor.FieldDescriptorProto_TYPE_DOUBLE:   codec.WireFixed64,
	descriptor.FieldDescriptorProto_TYPE_FLOAT:    codec.WireFixed32,
	descriptor.FieldDescriptorProto_TYPE_INT64:    codec.WireVarint,
	descriptor.FieldDescriptorProto_TYPE_UINT64:   codec.WireVarint,
	descriptor.FieldDescriptorProto_TYPE_INT32:    codec.WireVarint,
	descriptor.FieldDescriptorProto_TYPE_FIXED64:  codec.WireFixed64,
	descriptor.FieldDescriptorProto_TYPE_FIXED32:  codec.WireFixed32,
	descriptor.FieldDescriptorProto_TYPE_BOOL:     codec.WireVarint,
	descriptor.FieldDescriptorProto_TYPE_STRING:   codec.WireBytes,
	descriptor.FieldDescriptorProto_TYPE_MESSAGE:  codec.WireBytes,
	descriptor.FieldDescriptorProto_TYPE_BYTES:    codec.WireBytes,
	descriptor.FieldDescriptorProto_TYPE_UINT32:   codec.WireVarint,
	descriptor.FieldDescriptorProto_TYPE_ENUM:     codec.WireVarint,
	descriptor.FieldDescriptorProto_TYPE_SFIXED32: codec.WireFixed32,
	descriptor.FieldDescriptorProto_TYPE_SFIXED64: codec.WireFixed64,
	descriptor.FieldDescriptorProto_TYPE_SINT32:   codec.WireVarint,
	descriptor.FieldDescriptorProto_TYPE_SINT64:   codec.WireVarint,
}

func (g *generator) oMsgDecodeLoop(mode decodeMode) error {
	g.o(
		`
for !buf.EOF() {
	// We need to read a varint that represents the key that encodes the field number
	// and wire type. Speculate that the varint is one byte length and switch on it.
	// This is the hot path that is most common when field number is <= 15. 
	b := buf.PeekByteUnsafe()
	switch b {`,
	)

	g.i(1)
	var slowFields []*Field
	for _, field := range g.msg.Fields {
		if field.GetNumber() <= 15 {
			wireType, ok := protoTypeToWireType[field.GetType()]
			if ok {
				g.o(
					"case 0b0_%04b_%03b: // field number %d (%s), wire type %d (%s)",
					field.GetNumber(), wireType,
					field.GetNumber(), field.GetName(), wireType,
					wireTypeToString[wireType],
				)
				g.o(`	buf.SkipByteUnsafe()`)
				g.setField(field)
				g.oDecodeField(mode, false)
				continue
			}
		}
		slowFields = append(slowFields, field)
	}

	g.o(
		`
default:
	// Our speculation was wrong. Do the full (slow) decoding.`,
	)
	g.i(1)

	g.o(
		`
v, err := buf.DecodeVarint()
if err != nil {
	return err
}
fieldNum, wireType, err := codec.AsTagAndWireType(v)
if err != nil {
	return err
}
`,
	)

	g.o(`switch fieldNum {`)
	for _, field := range slowFields {
		g.setField(field)
		g.o(`case %d:`, field.GetNumber())
		g.i(1)
		g.o(`// Field %q`, field.GetName())
		g.oDecodeField(mode, true)
		g.i(-1)
	}
	g.o(`default:`)
	g.o(`	// Unknown field number.`)
	g.o(`	if err := buf.SkipFieldByWireType(wireType); err != nil {`)
	g.o(`		return err`)
	g.o(`	}`)

	g.o(`}`) // switch

	g.i(-1)

	g.o(`}`) // switch

	g.i(-1)

	g.o(`}`) // func

	return g.lastErr
}

func (g *generator) oDecodeField(mode decodeMode, checkWireType bool) {
	switch mode {
	case decodeValidate:
		fallthrough
	case decodeFull:
		g.oDecodeFieldValidateOrFull(mode, checkWireType)
	case decodeCountRepeat:
		g.oRepeatFieldCount()
	}
}

func (g *generator) oDecodeFieldValidateOrFull(mode decodeMode, checkWireType bool) {
	decode, ok := primitiveTypeDecode[g.field.GetType()]
	if ok {
		decode.mode = mode
		g.oDecodeFieldPrimitive(decode, checkWireType)
	} else if g.field.GetType() == descriptor.FieldDescriptorProto_TYPE_ENUM {
		g.oDecodeFieldEnum(
			g.enumDescrToEnum[g.field.GetEnumType()].GetName(), mode, checkWireType,
		)
	} else if g.field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		g.oDecodeFieldEmbeddedMessage(mode, checkWireType)
	} else {
		g.lastErr = fmt.Errorf("unsupported field type %v", g.field.GetType())
	}
}

func (g *generator) oRepeatFieldCount() {
	if g.field.IsRepeated() {
		counterName := g.field.GetName() + "Count"
		g.o(`	%s++`, counterName)
		g.o(`	buf.SkipRawBytes()`)
	} else {
		g.oSkipFieldByWireType()
	}
}

func (g *generator) oSkipFieldByWireType() {
	wireType := protoTypeToWireType[g.field.GetType()]
	switch wireType {
	case codec.WireVarint:
		g.o(`buf.SkipVarint()`)
	case codec.WireFixed64:
		g.o(`buf.SkipFixed64()`)
	case codec.WireBytes:
		g.o(`buf.SkipRawBytes()`)
	case codec.WireFixed32:
		g.o(`buf.SkipFixed32()`)
	}
}

type decodePrimitive struct {
	asProtoType      string
	oneOfType        string
	expectedWireType codec.WireType
	mode             decodeMode
}

var wireTypeToString = map[codec.WireType]string{
	codec.WireVarint:     "Varint",
	codec.WireFixed64:    "Fixed64",
	codec.WireBytes:      "Bytes",
	codec.WireStartGroup: "StartGroup",
	codec.WireEndGroup:   "EndGroup",
	codec.WireFixed32:    "Fixed32",
}

func (g *generator) oDecodeFieldPrimitive(task decodePrimitive, checkWireType bool) {
	if checkWireType {
		g.o(
			`
if wireType != codec.Wire%s {	
	return fmt.Errorf("invalid wire type %%d for field number %d ($MessageName.$fieldName)", wireType)
}`, wireTypeToString[task.expectedWireType], g.field.GetNumber(),
		)
	}

	if task.mode == decodeValidate {
		switch g.field.GetType() {
		case descriptor.FieldDescriptorProto_TYPE_STRING:
			fallthrough
		case descriptor.FieldDescriptorProto_TYPE_BYTES:
			g.o(`err := buf.SkipRawBytes()`)
		default:
			g.o(`_, err := buf.As%s()`, task.asProtoType)
		}
		g.o(
			`
if err != nil {
	return err
}`,
		)
		return
	} else {
		g.o(
			`
v, err := buf.As%s()
if err != nil {
	return err
}`, task.asProtoType,
		)
	}

	if g.field.GetOneOf() != nil {
		choiceName := composeOneOfChoiceName(g.msg, g.field)
		g.o(
			"m.%s = oneof.New%s(v, int(%s))", g.field.GetOneOf().GetName(),
			task.oneOfType, choiceName,
		)
	} else if g.field.IsRepeated() {
		counterName := g.field.GetName() + "Count"
		g.o(
			`
// The slice is pre-allocated, assign to the appropriate index.
m.$fieldName[%[1]s] = v
%[1]s++`, counterName,
		)
	} else {
		g.o(`m.$fieldName = v`)
		if g.options.WithPresence {
			g.o(`m._flags |= %s`, g.msg.PresenceFlagName[g.field])
		}
	}
}

func (g *generator) oDecodeFieldEnum(
	enumTypeName string, mode decodeMode, checkWireType bool,
) {

	if checkWireType {
		g.o(
			`
if wireType != codec.WireVarint {	
	return fmt.Errorf("invalid wire type %%d for field number %d ($MessageName.$fieldName)", wireType)
}`, g.field.GetNumber(),
		)
	}

	g.o(
		`
v, err := buf.AsUint32()
if err != nil {
	return err
}`,
	)

	if mode == decodeValidate {
		g.o(`_ = v`)
		return
	}

	g.o(`m.$fieldName = %s(v)`, enumTypeName)
	if g.options.WithPresence {
		g.o(`m._flags |= %s`, g.msg.PresenceFlagName[g.field])
	}
}

type decodeMode int

const (
	decodeValidate    decodeMode = 0
	decodeFull        decodeMode = 1
	decodeCountRepeat decodeMode = 2
)

var primitiveTypeDecode = map[descriptor.FieldDescriptorProto_Type]decodePrimitive{
	descriptor.FieldDescriptorProto_TYPE_BOOL: {
		asProtoType:      "Bool",
		oneOfType:        "Bool",
		expectedWireType: codec.WireVarint,
	},

	descriptor.FieldDescriptorProto_TYPE_FIXED64: {
		asProtoType:      "Fixed64",
		oneOfType:        "Int64",
		expectedWireType: codec.WireFixed64,
	},

	descriptor.FieldDescriptorProto_TYPE_UINT64: {
		asProtoType:      "Uint64",
		oneOfType:        "Uint32",
		expectedWireType: codec.WireVarint,
	},

	descriptor.FieldDescriptorProto_TYPE_SFIXED64: {
		asProtoType:      "SFixed64",
		oneOfType:        "Int64",
		expectedWireType: codec.WireFixed64,
	},

	descriptor.FieldDescriptorProto_TYPE_INT64: {
		asProtoType:      "Int64",
		oneOfType:        "Int64",
		expectedWireType: codec.WireVarint,
	},

	descriptor.FieldDescriptorProto_TYPE_FIXED32: {
		asProtoType:      "Fixed32",
		oneOfType:        "Int32",
		expectedWireType: codec.WireFixed32,
	},

	descriptor.FieldDescriptorProto_TYPE_SINT32: {
		asProtoType:      "Sint32",
		oneOfType:        "Sint32",
		expectedWireType: codec.WireVarint,
	},

	descriptor.FieldDescriptorProto_TYPE_UINT32: {
		asProtoType:      "Uint32",
		oneOfType:        "Uint32",
		expectedWireType: codec.WireVarint,
	},

	descriptor.FieldDescriptorProto_TYPE_DOUBLE: {
		asProtoType:      "Double",
		oneOfType:        "Double",
		expectedWireType: codec.WireFixed64,
	},

	descriptor.FieldDescriptorProto_TYPE_STRING: {
		asProtoType:      "StringUnsafe",
		oneOfType:        "String",
		expectedWireType: codec.WireBytes,
	},

	descriptor.FieldDescriptorProto_TYPE_BYTES: {
		asProtoType:      "BytesUnsafe",
		oneOfType:        "Bytes",
		expectedWireType: codec.WireBytes,
	},
}

func (g *generator) oDecodeFieldEmbeddedMessage(mode decodeMode, checkWireType bool) {

	if checkWireType {
		g.o(
			`
if wireType != codec.WireBytes {
	return fmt.Errorf("invalid wire type %%d for field number %d ($MessageName.$fieldName)", wireType)
}`, g.field.GetNumber(),
		)
	}

	if mode == decodeValidate {
		g.o(
			`
v, err := buf.DecodeRawBytes()
if err != nil {
	return err
}
err = validate$FieldMessageTypeName(v)
if err != nil {
	return err
}`,
		)
	} else {

		g.o(
			`
v, err := buf.AsBytesUnsafe()
if err != nil {
	return err
}
`,
		)

		if g.field.IsRepeated() {
			counterName := g.field.GetName() + "Count"
			g.o(
				`
// The slice is pre-allocated, assign to the appropriate index.
elem := m.$fieldName[%[1]s]
%[1]s++
elem._protoMessage.Parent = &m._protoMessage
elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)`, counterName,
			)
		} else if g.field.GetOneOf() != nil {
			choiceName := composeOneOfChoiceName(g.msg, g.field)
			g.o(
				`
elem := $fieldTypeMessagePool.Get()
elem._protoMessage.Parent = &m._protoMessage
elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
m.%s = oneof.NewPtr(unsafe.Pointer(elem), int(%s))`,
				g.field.GetOneOf().GetName(), choiceName,
			)
		} else {
			g.o(
				`
m.$fieldName = $fieldTypeMessagePool.Get()
m.$fieldName._protoMessage.Parent = &m._protoMessage
m.$fieldName._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)`,
			)
		}
	}
}

func (g *generator) oCalcRepeatedFieldCounts() {
	fields := g.getRepeatedFields()
	if len(fields) == 0 {
		return
	}

	g.o(``)
	g.o(`// Count all repeated fields. We need one counter per field.`)

	for _, field := range fields {
		g.setField(field)
		counterName := field.GetName() + "Count"
		g.o(`%s := 0`, counterName)
	}

	g.oMsgDecodeLoop(decodeCountRepeat)

	g.o(``)
	g.o(`// Pre-allocate slices for repeated fields.`)

	for _, field := range fields {
		g.setField(field)
		counterName := field.GetName() + "Count"

		if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			g.o(`if cap(m.$fieldName) < %s {`, counterName)
			g.o(`	// Need new space.`)
			g.o(`	m.$fieldName = make(%s, %s)`, g.convertTypeToGo(field), counterName)
			g.o(`} else {`)
			g.o(`	// Existing capacity is enough.`)
			g.o(`	m.$fieldName = m.$fieldName[0:%s]`, counterName)
			g.o(`}`)
			g.o(`$fieldTypeMessagePool.GetSlice(m.$fieldName)`)
		} else {
			//g.o(`if cap(m.$fieldName) < %s {`, counterName)
			g.o(`m.$fieldName = make(%s, %s)`, g.convertTypeToGo(field), counterName)
			//g.o(`} else {`)
			//g.o(`	m.$fieldName = m.$fieldName[0:%s]`, counterName)
			//g.o(`}`)
		}
	}
	g.o(``)
	g.o(`// Reset the buffer to start iterating over the fields again`)
	g.o(`buf.Reset(protomessage.BytesFromBytesView(m._protoMessage.Bytes))`)
	g.o(``)
	g.o(`// Set slice indexes to 0 to begin iterating over repeated fields.`)
	for _, field := range fields {
		g.setField(field)
		counterName := field.GetName() + "Count"
		g.o(`%s = 0`, counterName)
	}
}

func (g *generator) getRepeatedFields() []*Field {
	var r []*Field
	for _, field := range g.msg.Fields {
		if field.IsRepeated() {
			r = append(r, field)
		}
	}
	return r
}
