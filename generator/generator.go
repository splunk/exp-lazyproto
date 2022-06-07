package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	_ "github.com/jhump/protoreflect/desc/protoparse"
)

func Generate(protoPath string, inputProtoFiles []string, outputDir string) error {
	g := generator{
		protoPath:    protoPath,
		outputDir:    outputDir,
		templateData: map[string]string{},
	}

	for _, f := range inputProtoFiles {
		if err := g.processFile(f); err != nil {
			return err
		}
	}
	return g.lastErr
}

type generator struct {
	protoPath string
	outputDir string
	outBuf    *bytes.Buffer
	lastErr   error

	file         *File
	msg          *Message
	field        *Field
	templateData map[string]string
	spaces       int

	useSizedMarshaler bool
}

func (g *generator) processFile(inputFilePath string) error {
	p := protoparse.Parser{
		Accessor: func(filename string) (io.ReadCloser, error) {
			return os.Open(path.Join(g.protoPath, filename))
		},
		IncludeSourceCodeInfo: true,
	}
	fdescrs, err := p.ParseFiles(inputFilePath)
	if err != nil {
		return err
	}

	for _, fdescr := range fdescrs {
		if err := g.oStartFile(fdescr); err != nil {
			return err
		}
		for _, enum := range fdescr.GetEnumTypes() {
			if err := g.oEnum(enum); err != nil {
				return err
			}
		}
		for _, descr := range fdescr.GetMessageTypes() {
			msg := NewMessage(descr)
			g.setMessage(msg)
			if err := g.oMessage(msg); err != nil {
				return err
			}
		}
		if err := g.formatAndWriteToFile(fdescr); err != nil {
			return err
		}
	}

	return g.lastErr
}

func (g *generator) formatAndWriteToFile(fdescr *desc.FileDescriptor) error {
	fname := path.Base(fdescr.GetName()) + ".lz.go"
	fname = path.Join(g.outputDir, fname)
	fdir := path.Dir(fname)
	if err := os.MkdirAll(fdir, 0700); err != nil {
		return err
	}

	var err error
	f, err := os.Create(fname)
	if err != nil {
		return err
	}

	srcCode, err := format.Source(g.outBuf.Bytes())
	if err != nil {
		return err
	}

	_, err = f.Write(srcCode)
	return err
}

func (g *generator) oStartFile(fdescr *desc.FileDescriptor) error {
	g.outBuf = bytes.NewBuffer(nil)

	g.o("package %s", fdescr.GetPackage())
	g.o("")
	g.o(
		`
import (
	"sync"
	"unsafe"

	"github.com/tigrannajaryan/lazyproto/internal/protomessage"
	"github.com/tigrannajaryan/lazyproto/internal/oneof"

	"github.com/tigrannajaryan/molecule"
	"github.com/tigrannajaryan/molecule/src/codec"
)
`,
	)

	if g.useSizedMarshaler {
		g.o(`import "github.com/tigrannajaryan/lazyproto/internal/streams/sizedstream"`)
	}

	return g.lastErr
}

func (g *generator) setMessage(msg *Message) {
	g.msg = msg
	g.templateData["$MessageName"] = msg.GetName()
	g.templateData["$messagePool"] = getPoolName(msg.GetName())
}

func (g *generator) setField(field *Field) {
	g.field = field
	g.templateData["$fieldName"] = field.GetName()
	g.templateData["$FieldName"] = field.GetCapitalName()

	if field.GetMessageType() != nil {
		g.templateData["$fieldTypeMessagePool"] = getPoolName(field.GetMessageType().GetName())
		g.templateData["$FieldMessageTypeName"] = field.GetMessageType().GetName()
	} else {
		g.templateData["$fieldTypeMessagePool"] = "$fieldTypeMessagePool not defined for " + field.GetName()
		g.templateData["$FieldMessageTypeName"] = "$FieldMessageTypeName not defined for " + field.GetName()
	}
}

func (g *generator) o(str string, a ...any) {

	str = strings.TrimLeft(str, "\n")

	for k, v := range g.templateData {
		str = strings.ReplaceAll(str, k, v)
	}

	str = fmt.Sprintf(str, a...)

	strs := strings.Split(str, "\n")
	for i := range strs {
		if strings.TrimSpace(strs[i]) != "" {
			strs[i] = strings.Repeat("\t", g.spaces) + strs[i]
		}
	}

	str = strings.Join(strs, "\n")

	_, err := io.WriteString(g.outBuf, str+"\n")
	if err != nil {
		g.lastErr = err
	}
}

func (g *generator) i(ofs int) {
	g.spaces += ofs
}

func (g *generator) convertTypeToGo(field *Field) string {
	var s string

	if field.IsRepeated() {
		s = "[]"
	}

	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		s += "bool"
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		s += "uint64"
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		s += "int64"
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		s += "uint32"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		s += "uint32"
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		s += "float64"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		s += "string"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		s += "[]byte"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		s += "*" + field.GetMessageType().GetName()
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		s += field.GetEnumType().GetName()
	default:
		g.lastErr = fmt.Errorf("unsupported field type %v", field.GetType())
	}
	return s
}

func getLeadingComment(si *descriptor.SourceCodeInfo_Location) string {
	if si != nil {
		return strings.TrimSpace(si.GetLeadingComments())
	}
	return ""
}

func (g *generator) oMessage(msg *Message) error {
	g.o("// ====================== $MessageName message implementation ======================")
	g.o("")

	if err := g.oMsgStruct(); err != nil {
		return err
	}

	if err := g.oUnmarshalFree(); err != nil {
		return err
	}

	if err := g.oOneOf(); err != nil {
		return err
	}

	if err := g.oFieldsAccessors(msg); err != nil {
		return err
	}

	if err := g.oMsgDecodeFunc(msg); err != nil {
		return err
	}

	if err := g.oMarshalFunc(msg); err != nil {
		return err
	}

	if err := g.oPool(msg); err != nil {
		return err
	}

	g.o("")

	return nil
}

func (g *generator) oMsgStruct() error {
	c := getLeadingComment(g.msg.GetSourceInfo())
	if c != "" {
		g.o("// %s", c)
	}

	g.o("type $MessageName struct {")
	g.i(1)
	g.o("_protoMessage protomessage.ProtoMessage")

	if g.msg.NeedBitFlag {
		var bitCountFieldType string
		switch {
		case g.msg.FlagBitCount <= 8:
			bitCountFieldType = "uint8"
		case g.msg.FlagBitCount <= 16:
			bitCountFieldType = "uint16"
		case g.msg.FlagBitCount <= 32:
			bitCountFieldType = "uint32"
		case g.msg.FlagBitCount <= 64:
			bitCountFieldType = "uint64"
		default:
			return fmt.Errorf("more than 64 bits flags not supported")
		}

		g.o("_flags %s", bitCountFieldType)
	}
	g.o("")

	for _, field := range g.msg.Fields {
		if field.GetOneOf() != nil {
			// Skip oneof fields for now. They will be generated separately.
			continue
		}
		g.setField(field)
		c := getLeadingComment(field.GetSourceInfo())
		if c != "" {
			g.o("// %s", c)
		}
		g.o("$fieldName %s", g.convertTypeToGo(field))
	}

	// Generate oneof fields.
	for _, oneof := range g.msg.GetOneOfs() {
		g.o("%s oneof.OneOf", oneof.GetName())
	}
	g.i(-1)
	g.o("}")
	g.o("")

	return g.lastErr
}

func composeOneOfTypeName(msg *Message, oneof *desc.OneOfDescriptor) string {
	return msg.GetName() + capitalCamelCase(oneof.GetName())
}

func composeOneOfChoiceName(msg *Message, choice *Field) string {
	return fmt.Sprintf(
		"%s%s", msg.GetName(), choice.GetCapitalName(),
	)
}

func composeOneOfNoneChoiceName(msg *Message, oneof *desc.OneOfDescriptor) string {
	return msg.GetName() + capitalCamelCase(oneof.GetName()+"None")
}

func (g *generator) oOneOf() error {
	for _, oneof := range g.msg.GetOneOfs() {
		g.oOneOfTypeConsts(oneof)
		g.oOneOfTypeFunc(oneof)
	}

	return g.lastErr
}

func (g *generator) oOneOfTypeConsts(oneof *desc.OneOfDescriptor) {
	typeName := composeOneOfTypeName(g.msg, oneof)
	g.o(
		"// %s defines the possible types for oneof field %q.", typeName,
		oneof.GetName(),
	)
	g.o("type %s int\n", typeName)
	g.o("const (")
	g.i(1)

	noneChoiceName := composeOneOfNoneChoiceName(g.msg, oneof)
	g.o("// %s indicates that none of the oneof choices is set.", noneChoiceName)
	g.o("%s %s = 0", noneChoiceName, typeName)

	for i, choice := range oneof.GetChoices() {
		choiceField := g.msg.FieldsMap[choice.GetName()]
		choiceName := composeOneOfChoiceName(g.msg, choiceField)
		g.o(
			"// %s indicates that oneof field %q is set.", choiceName,
			choiceField.GetName(),
		)
		g.o("%s %s = %d", choiceName, typeName, i+1)
	}

	g.i(-1)
	g.o(")\n")
}

func (g *generator) oOneOfTypeFunc(oneof *desc.OneOfDescriptor) {
	typeName := composeOneOfTypeName(g.msg, oneof)
	funcName := fmt.Sprintf("%sType", capitalCamelCase(oneof.GetName()))
	g.o(
		"// %s returns the type of the current stored oneof %q.", funcName,
		oneof.GetName(),
	)
	g.o("// To set the type use one of the setters.")
	g.o("func (m *$MessageName) %s() %s {", funcName, typeName)
	g.o("	return %s(m.%s.FieldIndex())", typeName, oneof.GetName())
	g.o("}\n")

	funcName = fmt.Sprintf("%sUnset", capitalCamelCase(oneof.GetName()))
	g.o(
		"// %s unsets the oneof field %q, so that it contains none of the choices.",
		funcName, oneof.GetName(),
	)
	g.o("func (m *$MessageName) %s() {", funcName)
	g.o("	m.%s = oneof.NewOneOfNone()", oneof.GetName())
	g.o("}\n")
}

func (g *generator) oUnmarshalFree() error {
	g.o(
		`
func Unmarshal$MessageName(bytes []byte) (*$MessageName, error) {
	m := $messagePool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *$MessageName) Free() {
	$messagePool.Release(m)
}
`,
	)
	return g.lastErr
}

func (g *generator) oMsgDecodeFunc(msg *Message) error {
	g.o(
		`
func (m *$MessageName) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))`,
	)

	g.i(1)
	g.oRepeatedFieldCounts(msg)

	g.o("")
	g.o(
		`
// Iterate and decode the fields.
err2 := molecule.MessageEach(
	buf, func(fieldNum int32, value molecule.Value) (bool, error) {
		switch fieldNum {`,
	)

	g.i(2)
	g.oFieldDecode(msg.Fields)
	g.i(-2)

	g.i(-1)

	g.o(
		`
			}
			return true, nil
		},
	)
	if err2 != nil {
		return err2
	}
	return nil
}
`,
	)

	return g.lastErr
}

func (g *generator) oFieldDecodePrimitive(asProtoType string, oneOfType string) {
	g.o(
		`
v, err := value.As%s()
if err != nil {
	return false, err
}`, asProtoType,
	)

	if g.field.GetOneOf() != nil {
		choiceName := composeOneOfChoiceName(g.msg, g.field)
		g.o(
			"m.%s = oneof.NewOneOf%s(v, int(%s))", g.field.GetOneOf().GetName(),
			oneOfType, choiceName,
		)
	} else {
		g.o("m.$fieldName = v")
	}
}

func (g *generator) oFieldDecodeEnum(enumTypeName string) {
	g.o(
		`
v, err := value.AsUint32()
if err != nil {
	return false, err
}
m.$fieldName = %s(v)`, enumTypeName,
	)
}

func (g *generator) oFieldDecode(fields []*Field) string {
	for _, field := range fields {
		g.setField(field)
		g.o("case %d:", field.GetNumber())
		g.i(1)
		g.o(`// Decode "$fieldName".`)
		switch field.GetType() {
		case descriptor.FieldDescriptorProto_TYPE_BOOL:
			g.oFieldDecodePrimitive("Bool", "Bool")

		case descriptor.FieldDescriptorProto_TYPE_FIXED64:
			g.oFieldDecodePrimitive("Fixed64", "Int64")

		case descriptor.FieldDescriptorProto_TYPE_INT64:
			g.oFieldDecodePrimitive("Int64", "Int64")

		case descriptor.FieldDescriptorProto_TYPE_FIXED32:
			g.oFieldDecodePrimitive("Fixed32", "Int32")

		case descriptor.FieldDescriptorProto_TYPE_UINT32:
			g.oFieldDecodePrimitive("Uint32", "Uint32")

		case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
			g.oFieldDecodePrimitive("Double", "Double")

		case descriptor.FieldDescriptorProto_TYPE_ENUM:
			g.oFieldDecodeEnum(field.GetEnumType().GetName())

		case descriptor.FieldDescriptorProto_TYPE_STRING:
			g.oFieldDecodePrimitive("StringUnsafe", "String")

		case descriptor.FieldDescriptorProto_TYPE_BYTES:
			g.oFieldDecodePrimitive("BytesUnsafe", "Bytes")

		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			g.o(
				`
v, err := value.AsBytesUnsafe()
if err != nil {
	return false, err
}`,
			)

			if field.IsRepeated() {
				counterName := field.GetName() + "Count"
				g.o(
					`
// The slice is pre-allocated, assign to the appropriate index.
elem := m.$fieldName[%[1]s]
%[1]s++
elem._protoMessage.Parent = &m._protoMessage
elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)`, counterName,
				)
			} else if field.GetOneOf() != nil {
				choiceName := composeOneOfChoiceName(g.msg, field)
				g.o(
					`
elem := $fieldTypeMessagePool.Get()
elem._protoMessage.Parent = &m._protoMessage
elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
m.%s = oneof.NewOneOfPtr(unsafe.Pointer(elem), int(%s))`,
					field.GetOneOf().GetName(), choiceName,
				)
			} else {
				g.o(
					`
m.$fieldName = $fieldTypeMessagePool.Get()
m.$fieldName._protoMessage.Parent = &m._protoMessage
m.$fieldName._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)`,
				)
			}

		default:
			g.lastErr = fmt.Errorf("unsupported field type %v", field.GetType())
		}
		g.i(-1)
	}

	return ""
}

func (g *generator) oRepeatedFieldCounts(msg *Message) {
	fields := g.getRepeatedFields(msg)
	if len(fields) == 0 {
		return
	}

	g.o("")
	g.o("// Count all repeated fields. We need one counter per field.")

	for _, field := range fields {
		g.setField(field)
		counterName := field.GetName() + "Count"
		g.o("%s := 0", counterName)
	}

	g.o("err := molecule.MessageFieldNums(")
	g.o("	buf, func(fieldNum int32) {")
	for _, field := range fields {
		g.setField(field)
		counterName := field.GetName() + "Count"
		g.i(2)
		g.o("if fieldNum == %d {", field.GetNumber())
		g.o("	%s++", counterName)
		g.o("}")
		g.i(-2)
	}
	g.o("	},")
	g.o(")")
	g.o("if err != nil {")
	g.o("	return err")
	g.o("}")

	g.o("")
	g.o("// Pre-allocate slices for repeated fields.")

	for _, field := range fields {
		g.setField(field)
		counterName := field.GetName() + "Count"
		g.o("m.$fieldName = $fieldTypeMessagePool.GetSlice(%s)", counterName)
	}
	g.o("")
	g.o("// Reset the buffer to start iterating over the fields again")
	g.o("buf.Reset(protomessage.BytesFromBytesView(m._protoMessage.Bytes))")
	g.o("")
	g.o("// Set slice indexes to 0 to begin iterating over repeated fields.")
	for _, field := range fields {
		g.setField(field)
		counterName := field.GetName() + "Count"
		g.o("%s = 0", counterName)
	}
}

func (g *generator) getRepeatedFields(msg *Message) []*Field {
	var r []*Field
	for _, field := range msg.Fields {
		if field.IsRepeated() {
			r = append(r, field)
		}
	}
	return r
}

func (g *generator) oFieldsAccessors(msg *Message) error {
	if msg.NeedBitFlag {
		// Generate decode bit flags
		g.o("// Bitmasks that indicate that the particular nested message is decoded.")
		bitMask := uint64(1)
		for _, field := range msg.Fields {
			g.setField(field)
			if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {

				g.o("const %s = 0x%016X", g.fieldFlagName(), bitMask)
				bitMask *= 2
			}
		}
		g.o("")
	}

	for _, field := range msg.Fields {
		g.setField(field)
		if err := g.oFieldGetter(); err != nil {
			return err
		}
		if err := g.oFieldSetter(); err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) fieldFlagName() string {
	return fmt.Sprintf("flag%s%sDecoded", g.msg.GetName(), g.field.GetCapitalName())
}

func (g *generator) oFieldGetter() error {
	g.o("// $FieldName returns the value of the $fieldName.")

	if g.field.GetOneOf() != nil {
		g.o(
			`// If the field "%s" is not set to "$fieldName" then the returned value is undefined.`,
			g.field.GetOneOf().GetName(),
		)
	}

	goType := g.convertTypeToGo(g.field)

	g.o("func (m *$MessageName) $FieldName() %s {", goType)

	if g.field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		g.i(1)
		g.o("if m._flags&%s == 0 {", g.fieldFlagName())
		g.i(1)
		g.o("// Decode nested message(s).")
		if g.field.IsRepeated() {
			g.o("for i := range m.$fieldName {")
			g.o("	// TODO: decide how to handle decoding errors.")
			g.o("	_ = m.$fieldName[i].decode()")
			g.o("}")
		} else {
			if g.field.GetOneOf() != nil {
				choiceName := composeOneOfChoiceName(g.msg, g.field)
				g.o(
					"if m.%s.FieldIndex() == int(%s) {", g.field.GetOneOf().GetName(),
					choiceName,
				)
				g.i(1)
				g.o(
					"$fieldName := (*$FieldMessageTypeName)(m.%s.PtrVal())",
					g.field.GetOneOf().GetName(),
				)
			} else {
				g.o("$fieldName := m.$fieldName")
			}

			g.o("if $fieldName != nil {")
			g.o("	// TODO: decide how to handle decoding errors.")
			g.o("	_ = $fieldName.decode()")
			g.o("}")

			if g.field.GetOneOf() != nil {
				g.i(-1)
				g.o("}")
			}
		}
		g.i(-1)

		g.o("	m._flags |= %s", g.fieldFlagName())
		g.o("}")
		g.i(-1)
	}

	if g.field.GetOneOf() != nil {
		switch g.field.GetType() {
		case descriptor.FieldDescriptorProto_TYPE_BOOL:
			g.o("	return m.%s.BoolVal()", g.field.GetOneOf().GetName())
		case descriptor.FieldDescriptorProto_TYPE_INT64:
			g.o("	return m.%s.Int64Val()", g.field.GetOneOf().GetName())
		case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
			g.o("	return m.%s.DoubleVal()", g.field.GetOneOf().GetName())
		case descriptor.FieldDescriptorProto_TYPE_STRING:
			g.o("	return m.%s.StringVal()", g.field.GetOneOf().GetName())
		case descriptor.FieldDescriptorProto_TYPE_BYTES:
			g.o("	return m.%s.BytesVal()", g.field.GetOneOf().GetName())
		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			g.o("	return (%s)(m.%s.PtrVal())", goType, g.field.GetOneOf().GetName())
		default:
			return fmt.Errorf("unsupported oneof field type %v", g.field.GetType())
		}
	} else {
		g.o(`	return m.$fieldName`)
	}
	g.o("}\n")

	return g.lastErr
}

func (g *generator) calcOneOfFieldIndex() int {
	fieldIdx := -1
	for i, field := range g.field.GetOneOf().GetChoices() {
		if g.field.GetNumber() == field.GetNumber() {
			fieldIdx = i
			break
		}
	}
	if fieldIdx == -1 {
		g.lastErr = fmt.Errorf("cannot find index of oneof field %s", g.field.GetName())
	}
	return fieldIdx
}

func (g *generator) oFieldSetter() error {
	g.o("// Set$FieldName sets the value of the $fieldName.")

	if g.field.GetOneOf() != nil {
		g.o(
			`// The oneof field "%s" will be set to "$fieldName".`,
			g.field.GetOneOf().GetName(),
		)
	}

	g.o("func (m *$MessageName) Set$FieldName(v %s) {", g.convertTypeToGo(g.field))

	if g.field.GetOneOf() != nil {
		fieldIdx := g.calcOneOfFieldIndex()
		if fieldIdx == -1 {
			return g.lastErr
		}

		choiceName := composeOneOfChoiceName(g.msg, g.field)

		switch g.field.GetType() {
		case descriptor.FieldDescriptorProto_TYPE_BOOL:
			g.o(
				"	m.%s = oneof.NewOneOfBool(v, int(%s))",
				g.field.GetOneOf().GetName(), choiceName,
			)
		case descriptor.FieldDescriptorProto_TYPE_INT64:
			g.o(
				"	m.%s = oneof.NewOneOfInt64(v, int(%s))",
				g.field.GetOneOf().GetName(), choiceName,
			)
		case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
			g.o(
				"	m.%s = oneof.NewOneOfDouble(v, int(%s))",
				g.field.GetOneOf().GetName(), choiceName,
			)
		case descriptor.FieldDescriptorProto_TYPE_STRING:
			g.o(
				"	m.%s = oneof.NewOneOfString(v, int(%s))",
				g.field.GetOneOf().GetName(), choiceName,
			)
		case descriptor.FieldDescriptorProto_TYPE_BYTES:
			g.o(
				"	m.%s = oneof.NewOneOfBytes(v, int(%s))",
				g.field.GetOneOf().GetName(), choiceName,
			)
		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			g.o(
				"	m.%s = oneof.NewOneOfPtr(unsafe.Pointer(v), int(%s))",
				g.field.GetOneOf().GetName(), choiceName,
			)
		default:
			return fmt.Errorf("unsupported oneof field type %v", g.field.GetType())
		}
	} else {
		g.o("	m.$fieldName = v")
	}

	if g.field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		g.o("")
		g.o("	// Make sure the field's Parent points to this message.")
		if g.field.IsRepeated() {
			g.o("	for _, elem := range m.$fieldName {")
			g.o("		elem._protoMessage.Parent = &m._protoMessage")
			g.o("	}")
		} else {
			g.o("	v._protoMessage.Parent = &m._protoMessage")
		}
	}
	g.o("")
	g.o("	// Mark this message modified, if not already.")
	g.o("	m._protoMessage.MarkModified()")
	g.o("}\n")

	if g.field.IsRepeated() {
		if err := g.oFieldSliceMethods(); err != nil {
			return err
		}
	}

	return g.lastErr
}

func (g *generator) oFieldSliceMethods() error {
	g.o(
		`
func (m *$MessageName) $FieldNameRemoveIf(f func(*$FieldMessageTypeName) bool) {
	// Call getter to load the field.
	m.$FieldName()

	newLen := 0
	for i := 0; i < len(m.$fieldName); i++ {
		if f(m.$fieldName[i]) {
			continue
		}
		if newLen == i {
			// Nothing to move, element is at the right place.
			newLen++
			continue
		}
		m.$fieldName[newLen] = m.$fieldName[i]
		newLen++
	}
	if newLen != len(m.$fieldName) {
		m.$fieldName = m.$fieldName[:newLen]
		// Mark this message modified, if not already.
		m._protoMessage.MarkModified()
	}
}
`,
	)
	return g.lastErr
}

func (g *generator) oMarshalFunc(msg *Message) error {
	for _, field := range msg.Fields {
		g.setField(field)
		g.oPrepareMarshalField(msg, field)
	}

	g.o("")
	if g.useSizedMarshaler {
		g.o("func (m *$MessageName) Marshal(ps *sizedstream.ProtoStream) error {")
	} else {
		g.o("func (m *$MessageName) Marshal(ps *molecule.ProtoStream) error {")
	}
	g.i(1)
	g.o("if m._protoMessage.IsModified() {")
	g.i(1)

	// Order the fields by their number to ensure marshaling is done in an
	// order that we can rely on in the tests (same order as other as Protobuf
	// libs so that we can compare the results).
	fields := make([]*Field, len(msg.Fields))
	copy(fields, msg.Fields)
	sort.Slice(
		fields, func(i, j int) bool {
			return fields[i].GetNumber() < fields[j].GetNumber()
		},
	)

	for _, field := range fields {
		g.setField(field)

		if field.GetOneOf() != nil {
			fieldIndex := g.calcOneOfFieldIndex()
			if fieldIndex != 0 {
				// Already generated this oneof, nothing else to do.
				continue
			}
			g.o("// Marshal %q.", g.field.GetOneOf().GetName())

			typeName := composeOneOfTypeName(g.msg, g.field.GetOneOf())
			g.o("switch %s(m.%s.FieldIndex()) {", typeName, g.field.GetOneOf().GetName())

			// Add the "none" case.
			noneChoiceName := composeOneOfNoneChoiceName(g.msg, g.field.GetOneOf())
			g.o("case %s:", noneChoiceName)
			g.o("	// Nothing to do, oneof is unset.")

			for _, choice := range field.GetOneOf().GetChoices() {
				oneofField := g.msg.FieldsMap[choice.GetName()]
				g.setField(oneofField)
				typeName := composeOneOfChoiceName(g.msg, g.field)
				g.o("case %s:", typeName)
				g.i(1)
				g.oMarshalField()
				g.i(-1)
			}
			g.o("}")
		} else {
			g.oMarshalField()
		}
	}
	g.i(-1)
	g.o("} else {")
	g.o("	// Message is unchanged. Used original bytes.")
	g.o("	ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))")
	g.o("}")
	g.o("return nil")
	g.i(-1)
	g.o("}")
	return g.lastErr
}

func embeddedFieldName(msg *Message, field *Field) string {
	return fmt.Sprintf("prepared%s%s", msg.GetName(), field.GetCapitalName())
}

func (g *generator) oMarshalPreparedField(protoTypeName string) {
	if g.field.GetOneOf() != nil {
		g.o(
			"ps.%[1]sPrepared(prepared$MessageName$FieldName, m.%[2]s.%[1]sVal())",
			protoTypeName,
			g.field.GetOneOf().GetName(),
		)
	} else {
		g.o("ps.%sPrepared(prepared$MessageName$FieldName, m.$fieldName)", protoTypeName)
	}
}

func (g *generator) oMarshalField() {
	g.o(`// Marshal "$fieldName".`)
	switch g.field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		g.oMarshalPreparedField("Bool")

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		g.oMarshalPreparedField("String")

	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		g.oMarshalPreparedField("Bytes")

	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		g.oMarshalPreparedField("Fixed64")

	case descriptor.FieldDescriptorProto_TYPE_INT64:
		g.oMarshalPreparedField("Int64")

	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		g.oMarshalPreparedField("Fixed32")

	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		g.oMarshalPreparedField("Uint32")

	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		g.oMarshalPreparedField("Double")

	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		g.o("ps.Uint32Prepared(prepared$MessageName$FieldName, uint32(m.$fieldName))")

	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		g.oMarshalMessageTypeField()

	default:
		g.lastErr = fmt.Errorf("unsupported field type %v", g.field.GetType())
	}
}

func (g *generator) oMarshalMessageTypeField() {
	if g.field.IsRepeated() {
		g.o("for _, elem := range m.$fieldName {")
		g.o("	token := ps.BeginEmbedded()")
		g.o("	if err := elem.Marshal(ps); err != nil {")
		g.o("		return err")
		g.o("	}")
		g.o(
			"	ps.EndEmbeddedPrepared(token, %s)",
			embeddedFieldName(g.msg, g.field),
		)
	} else {
		if g.field.GetOneOf() != nil {
			g.o(
				"elem := (*$FieldMessageTypeName)(m.%s.PtrVal())",
				g.field.GetOneOf().GetName(),
			)
		} else {
			g.o("elem := m.$fieldName")
		}

		g.o("if elem != nil {")
		g.o("	token := ps.BeginEmbedded()")
		g.o("	if err := elem.Marshal(ps); err != nil {")
		g.o("		return err")
		g.o("	}")
		g.o(
			"	ps.EndEmbeddedPrepared(token, %s)",
			embeddedFieldName(g.msg, g.field),
		)
	}
	g.o("}")
}

func (g *generator) oPrepareMarshalField(msg *Message, field *Field) {
	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		g.o(g.preparedFieldDecl(msg, field, "Bool"))

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		g.o(g.preparedFieldDecl(msg, field, "String"))

	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		g.o(g.preparedFieldDecl(msg, field, "Bytes"))

	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		g.o(g.preparedFieldDecl(msg, field, "Fixed64"))

	case descriptor.FieldDescriptorProto_TYPE_INT64:
		g.o(g.preparedFieldDecl(msg, field, "Int64"))

	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		g.o(g.preparedFieldDecl(msg, field, "Fixed32"))

	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		g.o(g.preparedFieldDecl(msg, field, "Uint32"))

	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		g.o(g.preparedFieldDecl(msg, field, "Double"))

	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		g.o(g.preparedFieldDecl(msg, field, "Uint32"))

	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		g.o(g.preparedFieldDecl(msg, field, "Embedded"))

	default:
		g.lastErr = fmt.Errorf("unsupported field type %v", field.GetType())
	}
}

func unexportedName(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

func (g *generator) oPool(msg *Message) error {
	g.o("")
	g.o(
		`
// Pool of $MessageName structs.
type $messagePoolType struct {
	pool []*$MessageName
	mux  sync.Mutex
}

var $messagePool = $messagePoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *$messagePoolType) Get() *$MessageName {
	p.mux.Lock()
	defer p.mux.Unlock()

	// Have elements in the pool?
	if len(p.pool) >= 1 {
		// Get the last element.
		r := p.pool[len(p.pool)-1]
		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-1]
		return r
	}

	// Pool is empty, create a new element.
	return &$MessageName{}
}

func (p *$messagePoolType) GetSlice(count int) []*$MessageName {
	// Create a new slice.
	r := make([]*$MessageName, count)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return r
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]$MessageName, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}

	return r
}`,
	)

	g.oPoolReleaseSlice(msg)
	g.oPoolRelease(msg)

	return g.lastErr
}

func getPoolName(msgName string) string {
	return unexportedName(msgName) + "Pool"
}

func (g *generator) oPoolReleaseElem(msg *Message) {
	for _, field := range msg.Fields {
		g.setField(field)

		if field.GetOneOf() != nil {
			fieldIndex := g.calcOneOfFieldIndex()
			if fieldIndex != 0 {
				// We generate all oneof cases when we see the first field. Skip for the rest.
				continue
			}

			typeName := composeOneOfTypeName(g.msg, field.GetOneOf())
			g.o("switch %s(elem.%s.FieldIndex()) {", typeName, field.GetOneOf().GetName())
			for _, choice := range field.GetOneOf().GetChoices() {
				choiceField := g.msg.FieldsMap[choice.GetName()]
				g.setField(choiceField)
				if choiceField.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
					choiceName := composeOneOfChoiceName(g.msg, choiceField)
					g.o("case %s:", choiceName)
					g.i(1)
					g.o(
						"ptr := (*$FieldMessageTypeName)(elem.%s.PtrVal())",
						field.GetOneOf().GetName(),
					)
					g.o("if ptr != nil {")
					g.o("	$fieldTypeMessagePool.Release(ptr)")
					g.o("}")
					g.i(-1)
				}
			}
			g.o("}")
			continue
		}

		if field.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			// Only embedded messages need to be freed.
			continue
		}

		// Not a oneof field.
		g.o("// Release nested $fieldName recursively to their pool.")
		if field.IsRepeated() {
			g.o("$fieldTypeMessagePool.ReleaseSlice(elem.$fieldName)")
		} else {
			g.o("if elem.$fieldName != nil {")
			g.o("	$fieldTypeMessagePool.Release(elem.$fieldName)")
			g.o("}")
		}
	}

	g.o("")
	g.o(
		`
// Zero-initialize the released element.
*elem = $MessageName{}`,
	)
}

func (g *generator) oPoolReleaseSlice(msg *Message) {
	g.o("")
	g.o(
		`
// ReleaseSlice releases a slice of elements back to the pool.
func (p *$messagePoolType) ReleaseSlice(slice []*$MessageName) {
	for _, elem := range slice {`,
	)

	g.i(2)
	g.oPoolReleaseElem(msg)
	g.i(-2)

	g.o(
		`	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}`,
	)
}

func (g *generator) oPoolRelease(msg *Message) {
	g.o("")
	g.o(
		`
// Release an element back to the pool.
func (p *$messagePoolType) Release(elem *$MessageName) {`,
	)

	g.i(1)
	g.oPoolReleaseElem(msg)
	g.i(-1)

	g.o("")
	g.o(
		`
	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}`,
	)
}

func (g *generator) preparedFieldDecl(
	msg *Message, field *Field, protoTypeName string,
) string {
	if g.useSizedMarshaler {
		return fmt.Sprintf(
			"var prepared%s%s = sizedstream.Prepare%sField(%d)", msg.GetName(),
			field.GetCapitalName(), protoTypeName, field.GetNumber(),
		)
	} else {
		return fmt.Sprintf(
			"var prepared%s%s = molecule.Prepare%sField(%d)", msg.GetName(),
			field.GetCapitalName(), protoTypeName, field.GetNumber(),
		)
	}
}

func (g *generator) oEnum(enum *desc.EnumDescriptor) error {
	c := getLeadingComment(enum.GetSourceInfo())
	if c != "" {
		g.o("// %s", c)
	}
	g.o(
		`
type %[1]s uint32

const (`, enum.GetName(),
	)

	g.i(1)
	for _, value := range enum.GetValues() {
		c := getLeadingComment(value.GetSourceInfo())
		if c != "" {
			g.o("// %s", c)
		}
		g.o(
			"%[1]s_%[2]s %[1]s = %[3]d", enum.GetName(), value.GetName(),
			value.GetNumber(),
		)
	}
	g.i(-1)

	g.o(")")
	g.o("")

	return g.lastErr
}
