package generator

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	_ "github.com/jhump/protoreflect/desc/protoparse"
)

func Generate(inputProtoFiles []string, outputDir string) error {
	g := generator{outputDir: outputDir, templateData: map[string]string{}}

	for _, f := range inputProtoFiles {
		if err := g.processFile(f); err != nil {
			return err
		}
	}
	return nil
}

type generator struct {
	outputDir string
	outF      *os.File
	lastErr   error

	file         *File
	msg          *Message
	field        *Field
	templateData map[string]string
	spaces       int
}

func (g *generator) processFile(inputFilePath string) error {
	p := protoparse.Parser{
		Accessor: func(filename string) (io.ReadCloser, error) {
			return os.Open(inputFilePath)
		},
		IncludeSourceCodeInfo: true,
	}
	fdescrs, err := p.ParseFiles(inputFilePath)
	if err != nil {
		return err
	}

	for _, fdescr := range fdescrs {
		if err := g.oFile(fdescr); err != nil {
			return err
		}
		for _, descr := range fdescr.GetMessageTypes() {
			msg := NewMessage(descr)
			g.setMessage(msg)
			if err := g.oMessage(msg); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *generator) oFile(fdescr *desc.FileDescriptor) error {
	fname := path.Base(fdescr.GetName()) + ".lz.go"
	fname = path.Join(g.outputDir, fname)
	fdir := path.Dir(fname)
	if err := os.MkdirAll(fdir, 0700); err != nil {
		return err
	}

	var err error
	g.outF, err = os.Create(fname)
	if err != nil {
		return err
	}

	g.o("package %s", fdescr.GetPackage())
	g.o("")
	g.o(
		`
import (
	"sync"

	lazyproto "github.com/tigrannajaryan/exp-lazyproto"
	"github.com/tigrannajaryan/molecule"
	"github.com/tigrannajaryan/molecule/src/codec"
	protostream "github.com/tigrannajaryan/exp-lazyproto/internal/stream"
)
`,
	)

	return nil
}

func (g *generator) setMessage(msg *Message) {
	g.msg = msg
	g.templateData["MessageName"] = msg.GetName()
	g.templateData["messagePool"] = getPoolName(msg.GetName())
}

func (g *generator) setField(field *Field) {
	g.field = field
	g.templateData["fieldName"] = field.GetName()
	g.templateData["FieldName"] = field.GetCapitalName()

	if field.GetMessageType() != nil {
		g.templateData["FieldMessageTypeName"] = field.GetMessageType().GetName()
		g.templateData["fieldTypeMessagePool"] = getPoolName(field.GetMessageType().GetName())
	} else {
		g.templateData["FieldMessageTypeName"] = "FieldMessageTypeName not defined for " + field.GetName()
		g.templateData["fieldTypeMessagePool"] = "fieldTypeMessagePool not defined for " + field.GetName()
	}
}

func (g *generator) o(str string, a ...any) bool {

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

	_, err := g.outF.WriteString(str + "\n")
	if err != nil {
		g.lastErr = err
		return false
	}
	return true
}

func (g *generator) i(ofs int) {
	g.spaces += ofs
}

func (g *generator) convertType(field *Field) string {
	var s string

	if field.IsRepeated() {
		s = "[]"
	}

	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		s += "uint64"
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		s += "uint32"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		s += "string"
	default:
		s += "*" + field.GetMessageType().GetName()
	}
	return s
}

func (g *generator) oMessage(msg *Message) error {
	g.o("// ====================== Generated for message MessageName ======================")
	g.o("")
	si := msg.GetSourceInfo()
	if si != nil {
		c := strings.TrimSpace(si.GetLeadingComments())
		if c != "" {
			g.o("// %s", c)
		}
	}

	g.o("type MessageName struct {")
	g.i(1)
	g.o("protoMessage lazyproto.ProtoMessage")
	for _, field := range msg.Fields {
		g.setField(field)
		si := field.GetSourceInfo()
		if si != nil {
			c := strings.TrimSpace(si.GetLeadingComments())
			if c != "" {
				g.o("// %s", c)
			}
		}
		g.o("fieldName %s", g.convertType(field))
	}
	g.i(-1)
	g.o("}")
	g.o("")

	g.o(
		`
func NewMessageName(bytes []byte) *MessageName {
	m := messagePool.Get()
	m.protoMessage.Bytes = bytes
	m.decode()
	return m
}

func (m *MessageName) Free() {
	messagePool.Release(m)
}`,
	)

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

func (g *generator) oMsgDecodeFunc(msg *Message) error {
	g.o("")
	g.o(
		`
func (m *MessageName) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)`,
	)

	g.i(1)
	g.oRepeatedFieldCounts(msg)

	g.o("")
	g.o(
		`
// Iterate and decode the fields.
molecule.MessageEach(
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
}
`,
	)

	return g.lastErr
}

func (g *generator) oFieldDecodePrimitive(field *Field, asType string) {
	g.o(
		`
v, err := value.As%s()
if err != nil {
	return false, err
}
m.fieldName = v`, asType,
	)
}

func (g *generator) oFieldDecode(fields []*Field) string {
	for _, field := range fields {
		g.setField(field)
		g.o("case %d:", field.GetNumber())
		g.i(1)
		g.o("// Decode fieldName.")
		switch field.GetType() {
		case descriptor.FieldDescriptorProto_TYPE_FIXED64:
			g.oFieldDecodePrimitive(field, "Fixed64")

		case descriptor.FieldDescriptorProto_TYPE_UINT32:
			g.oFieldDecodePrimitive(field, "Uint32")

		case descriptor.FieldDescriptorProto_TYPE_STRING:
			g.oFieldDecodePrimitive(field, "StringUnsafe")

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
elem := m.fieldName[%[1]s]
%[1]s++
elem.protoMessage.Parent = &m.protoMessage
elem.protoMessage.Bytes = v`, counterName,
				)
			} else {
				g.o(
					`
m.fieldName = fieldTypeMessagePool.Get()
m.fieldName.protoMessage.Parent = &m.protoMessage
m.fieldName.protoMessage.Bytes = v`,
				)
			}
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

	g.o("molecule.MessageFieldNums(")
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

	g.o("")
	g.o("// Pre-allocate slices for repeated fields.")

	for _, field := range fields {
		g.setField(field)
		counterName := field.GetName() + "Count"
		g.o("m.fieldName = fieldTypeMessagePool.GetSlice(%s)", counterName)
	}
	g.o("")
	g.o("// Reset the buffer to start iterating over the fields again")
	g.o("buf.Reset(m.protoMessage.Bytes)")
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
	// Generate decode bit flags
	bitMask := uint64(2) // Start from 2 since bit 1 is used for "flagsMessageModified"
	firstFlag := true
	for _, field := range msg.Fields {
		g.setField(field)
		if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			if firstFlag {
				g.o("")
				g.o("// Bitmasks that indicate that the particular nested message is decoded.")
			}

			g.o("const %s = 0x%016X", g.fieldFlagName(msg, field), bitMask)
			bitMask *= 2
			firstFlag = false
		}
	}

	for _, field := range msg.Fields {
		g.setField(field)
		if err := g.oFieldGetter(msg, field); err != nil {
			return err
		}
		if err := g.oFieldSetter(msg, field); err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) fieldFlagName(msg *Message, field *Field) string {
	return fmt.Sprintf("flag%s%sDecoded", msg.GetName(), field.GetCapitalName())
}

func (g *generator) oFieldGetter(msg *Message, field *Field) error {
	g.o("")
	g.o("func (m *MessageName) FieldName() %s {", g.convertType(field))

	if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		g.i(1)
		g.o("if m.protoMessage.Flags&%s == 0 {", g.fieldFlagName(msg, field))
		g.i(1)
		g.o("// Decode nested message(s).")
		if field.IsRepeated() {
			g.o("for i := range m.fieldName {")
			g.o("	m.fieldName[i].decode()")
			g.o("}")
		} else {
			g.o("m.fieldName.decode()")
		}
		g.i(-1)

		g.o("	m.protoMessage.Flags |= %s", g.fieldFlagName(msg, field))
		g.o("}")
		g.i(-1)
	}

	g.o(
		`	return m.fieldName
}`,
	)

	return g.lastErr
}

func (g *generator) oFieldSetter(msg *Message, field *Field) error {
	g.o("")
	g.o("func (m *MessageName) SetFieldName(v %s) {", g.convertType(field))
	g.o("	m.fieldName = v")
	g.o("	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {")
	g.o("		m.protoMessage.MarkModified()")
	g.o("	}")
	g.o("}")

	return g.lastErr
}

func (g *generator) oMarshalFunc(msg *Message) error {
	for _, field := range msg.Fields {
		g.setField(field)
		g.oPrepareMarshalField(msg, field)
	}

	g.o("")
	g.o("func (m *MessageName) Marshal(ps *protostream.ProtoStream) error {")
	g.i(1)
	g.o("if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {")
	g.i(1)
	//for _, field := range msg.Fields {
	for i := len(msg.Fields) - 1; i >= 0; i-- {
		field := msg.Fields[i]
		g.setField(field)
		g.oMarshalField(msg, field)
	}
	g.i(-1)
	g.o("} else {")
	g.o("	// Message is unchanged. Used original bytes.")
	g.o("	ps.Raw(m.protoMessage.Bytes)")
	g.o("}")
	g.o("return nil")
	g.i(-1)
	g.o("}")
	return g.lastErr
}

func embeddedFieldName(msg *Message, field *Field) string {
	return fmt.Sprintf("prepared%s%s", msg.GetName(), field.GetCapitalName())
}

func (g *generator) oMarshalPreparedField(msg *Message, field *Field, typeName string) {
	g.o("ps.%sPrepared(preparedMessageNameFieldName, m.fieldName)", typeName)
}

func (g *generator) oMarshalField(msg *Message, field *Field) {
	g.o("// Marshal fieldName")
	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		g.oMarshalPreparedField(msg, field, "String")

	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		g.oMarshalPreparedField(msg, field, "Fixed64")

	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		g.oMarshalPreparedField(msg, field, "Uint32")

	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		if field.IsRepeated() {
			g.o("for i:=len(m.fieldName)-1; i>=0; i-- {")
			g.o("   elem := m.fieldName[i]")
			g.o("	token := ps.BeginEmbedded()")
			g.o("	elem.Marshal(ps)")
			g.o("	ps.EndEmbeddedPrepared(token, %s)", embeddedFieldName(msg, field))
		} else {
			g.o("if m.fieldName != nil {")
			g.o("	token := ps.BeginEmbedded()")
			g.o("	m.fieldName.Marshal(ps)")
			g.o("	ps.EndEmbeddedPrepared(token, %s)", embeddedFieldName(msg, field))
		}
		g.o("}")
	}
}

func (g *generator) oPrepareMarshalField(msg *Message, field *Field) {
	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		g.o(preparedFieldDecl(msg, field, "String"))

	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		g.o(preparedFieldDecl(msg, field, "Fixed64"))

	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		g.o(preparedFieldDecl(msg, field, "Uint32"))

	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		g.o(preparedFieldDecl(msg, field, "Embedded"))
	}
}

func unexportedName(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

func (g *generator) oPool(msg *Message) error {
	g.o("")
	g.o(
		`
// Pool of MessageName structs.
type messagePoolType struct {
	pool []*MessageName
	mux  sync.Mutex
}

var messagePool = messagePoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *messagePoolType) Get() *MessageName {
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
	return &MessageName{}
}

func (p *messagePoolType) GetSlice(count int) []*MessageName {
	// Create a new slice.
	r := make([]*MessageName, count)

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
		storage := make([]MessageName, count-copied)
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

		if field.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			continue
		}

		g.o("// Release nested fieldName recursively to their pool.")
		if field.IsRepeated() {
			g.o("fieldTypeMessagePool.ReleaseSlice(elem.fieldName)")
		} else {
			g.o("if elem.fieldName != nil {")
			g.o("	fieldTypeMessagePool.Release(elem.fieldName)")
			g.o("}")
		}
	}

	g.o("")
	g.o(
		`
// Zero-initialize the released element.
*elem = MessageName{}`,
	)
}

func (g *generator) oPoolReleaseSlice(msg *Message) {
	g.o("")
	g.o(
		`
// ReleaseSlice releases a slice of elements back to the pool.
func (p *messagePoolType) ReleaseSlice(slice []*MessageName) {
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
func (p *messagePoolType) Release(elem *MessageName) {`,
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

func preparedFieldDecl(msg *Message, field *Field, typeName string) string {
	return fmt.Sprintf(
		"var prepared%s%s = protostream.Prepare%sField(%d)", msg.GetName(),
		field.GetCapitalName(), typeName, field.GetNumber(),
	)
}
