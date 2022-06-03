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
	"github.com/yoheimuta/go-protoparser/v4/parser"

	_ "github.com/jhump/protoreflect/desc/protoparse"
)

func Generate(inputProtoFiles []string, outputDir string) error {
	g := generator{outputDir: outputDir}

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

	file   *File
	msg    *Message
	field  *parser.Field
	spaces int
}

func (g *generator) VisitComment(comment *parser.Comment) {
}

func (g *generator) VisitEmptyStatement(statement *parser.EmptyStatement) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) VisitEnum(enum *parser.Enum) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) VisitEnumField(field *parser.EnumField) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) VisitExtend(extend *parser.Extend) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) VisitExtensions(extensions *parser.Extensions) (next bool) {
	//TODO implement me
	panic("implement me")
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

	//asts, err := p.ParseToAST(inputFilePath)
	//if err != nil {
	//	return err
	//}
	//for _, ast := range asts {
	//	ast.GetSyntax()
	//}

	for _, fdescr := range fdescrs {
		if err := g.oFile(fdescr); err != nil {
			return err
		}
		for _, descr := range fdescr.GetMessageTypes() {
			msg := NewMessage(descr)
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

	g.o("package %s\n", fdescr.GetPackage())
	g.o(
		`
import (
	"sync"

	lazyproto "github.com/tigrannajaryan/exp-lazyproto"
	"github.com/tigrannajaryan/molecule"
	"github.com/tigrannajaryan/molecule/src/codec"
)

`,
	)

	return nil
}

func (g *generator) o(str string, a ...any) bool {

	str = fmt.Sprintf(str, a...)

	strs := strings.Split(str, "\n")
	for i := range strs {
		if strings.TrimSpace(strs[i]) != "" {
			strs[i] = strings.Repeat("\t", g.spaces) + strs[i]
		}
	}

	str = strings.Join(strs, "\n")

	_, err := g.outF.WriteString(str)
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
	g.o(
		"// ====================== Generated for message %s ======================\n\n",
		msg.GetName(),
	)
	si := msg.GetSourceInfo()
	if si != nil {
		if si.GetLeadingComments() != "" {
			g.o("//%s", si.GetLeadingComments())
		}
	}

	g.o("type %s struct {\n", msg.GetName())
	g.i(1)
	g.o("protoMessage lazyproto.ProtoMessage\n")
	for _, field := range msg.Fields {
		si := field.GetSourceInfo()
		if si != nil {
			if si.GetLeadingComments() != "" {
				g.o("//%s", si.GetLeadingComments())
			}
		}
		g.o("%s %s\n", field.GetName(), g.convertType(field))
	}
	g.i(-1)
	g.o("}\n")

	g.o(
		`
func New%s(bytes []byte) *%s {
	m := &%s{protoMessage: lazyproto.ProtoMessage{Bytes: bytes}}
	m.decode()
	return m
}
`, msg.GetName(), msg.GetName(), msg.GetName(),
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

	g.o("\n")

	return nil
}

func (g *generator) oMsgDecodeFunc(msg *Message) error {
	g.o(
		`
func (m *%s) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)
`, msg.GetName(),
	)

	g.i(1)
	g.oRepeatedFieldCounts(msg)

	g.o(
		`
// Iterate and decode the fields.
molecule.MessageEach(
	buf, func(fieldNum int32, value molecule.Value) (bool, error) {
		switch fieldNum {
`,
	)

	g.i(2)
	g.oFieldDecode(msg.Fields)
	g.i(-2)

	g.i(-1)

	g.o(
		`			}
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
m.%s = v
`, asType, field.GetName(),
	)
}

func (g *generator) oFieldDecode(fields []*Field) string {
	for _, field := range fields {
		g.o("case %d:\n", field.GetNumber())
		g.i(1)
		g.o("// Decode %s.", field.GetName())
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
m.%s[%s] = &%s{
	protoMessage: lazyproto.ProtoMessage{
		Parent: &m.protoMessage, Bytes: v,
	},
}
%s++
`, field.GetName(), counterName, field.GetMessageType().GetName(), counterName,
				)
			} else {
				g.o(
					`
m.%s = &%s{
	protoMessage: lazyproto.ProtoMessage{
		Parent: &m.protoMessage, Bytes: v,
	},
}
`, field.GetName(), field.GetMessageType().GetName(),
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

	g.o("\n// Count all repeated fields. We need one counter per field.\n")

	for _, field := range fields {
		counterName := field.GetName() + "Count"
		g.o("%s := 0", counterName)
	}

	g.o(
		`
molecule.MessageFieldNums(
	buf, func(fieldNum int32) {`,
	)
	for _, field := range fields {
		counterName := field.GetName() + "Count"
		g.i(2)
		g.o(
			`
if fieldNum == %d {
	%s++
}`, field.GetNumber(), counterName,
		)
		g.i(-2)
	}
	g.o(
		`
	},
)
`,
	)

	g.o("\n// Pre-allocate slices for repeated fields.\n")

	for _, field := range fields {
		counterName := field.GetName() + "Count"
		g.o(
			"m.%s = %s.GetSlice(%s)\n", field.GetName(),
			getPoolName(field.GetMessageType().GetName()),
			counterName,
		)
	}
	g.o("\n// Reset the buffer to start iterating over the fields again")
	g.o("\nbuf.Reset(m.protoMessage.Bytes)\n")
	g.o("\n// Set slice indexes to 0 to begin iterating over repeated fields.\n")
	for _, field := range fields {
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
		if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			if firstFlag {
				g.o("\n// Bitmasks that indicate that the particular nested message is decoded.\n")
			}

			g.o("const %s = 0x%016X\n", g.fieldFlagName(msg, field), bitMask)
			bitMask *= 2
			firstFlag = false
		}
	}

	for _, field := range msg.Fields {
		if err := g.FieldAccessors(msg, field); err != nil {
			return err
		}
	}
	return nil
}

func (g *generator) fieldFlagName(msg *Message, field *Field) string {
	return fmt.Sprintf("flag%s%sDecoded", msg.GetName(), field.GetCapitalName())
}

func (g *generator) FieldAccessors(msg *Message, field *Field) error {
	g.o(
		"\nfunc (m *%s) Get%s() %s {\n", msg.GetName(), field.GetCapitalName(),
		g.convertType(field),
	)

	if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		g.i(1)
		g.o("if m.protoMessage.Flags&%s == 0 {\n", g.fieldFlagName(msg, field))
		g.i(1)
		g.o("// Decode nested message(s).\n")
		if field.IsRepeated() {
			g.o("for i := range m.%s {\n", field.GetName())
			g.o("	m.%s[i].decode()\n", field.GetName())
			g.o("}\n")

		} else {
			g.o("m.%s.decode()\n", field.GetName())
		}
		g.i(-1)

		g.o("	m.protoMessage.Flags |= %s\n", g.fieldFlagName(msg, field))
		g.o("}\n")
		g.i(-1)
	}

	g.o(
		`	return m.%s
}
`, field.GetName(),
	)
	return g.lastErr
}

func (g *generator) oMarshalFunc(msg *Message) error {
	for _, field := range msg.Fields {
		g.oPrepareMarshalField(msg, field)
	}

	g.o("\nfunc (m *%s) Marshal(ps *molecule.ProtoStream) error {\n", msg.GetName())
	g.i(1)
	g.o("if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {\n")
	g.i(1)
	for _, field := range msg.Fields {
		g.oMarshalField(msg, field)
	}
	g.i(-1)
	g.o("} else {\n")
	g.o("	// Message is unchanged. Used original bytes.\n")
	g.o("	ps.Raw(m.protoMessage.Bytes)\n")
	g.o("}\n")
	g.o("return nil\n")
	g.i(-1)
	g.o("}\n")
	return g.lastErr
}

func (g *generator) oMarshalPreparedField(msg *Message, field *Field, typeName string) {
	g.o(
		"ps.%sPrepared(prepared%s%s, m.%s)\n", typeName, msg.GetName(),
		field.GetCapitalName(), field.GetName(),
	)
}

func (g *generator) oMarshalField(msg *Message, field *Field) {
	g.o("// Marshal %s\n", field.camelName)
	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		g.oMarshalPreparedField(msg, field, "String")

	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		g.oMarshalPreparedField(msg, field, "Fixed64")

	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		g.oMarshalPreparedField(msg, field, "Uint32")

	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		if field.IsRepeated() {
			g.o("for _, elem := range m.%s {\n", field.camelName)
			g.i(1)
			g.o("token := ps.BeginEmbedded()\n")
			g.o("elem.Marshal(ps)\n")
			g.o("ps.EndEmbedded(token, %d)\n", field.GetNumber())
			g.i(-1)
		} else {
			g.o("if m.%s != nil {\n", field.camelName)
			g.i(1)
			g.o("token := ps.BeginEmbedded()\n")
			g.o("m.%s.Marshal(ps)\n", field.camelName)
			g.o("ps.EndEmbedded(token, %d)\n", field.GetNumber())
			g.i(-1)
		}

		g.o("}\n")
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
	}
}

func unexportedName(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

func (g *generator) oPool(msg *Message) error {
	poolName := unexportedName(msg.GetName()) + "Pool"

	g.o(
		`
// Pool of %s structs.
type %sType struct {
	pool []*%s
	mux  sync.Mutex
}

var %s = %sType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *%sType) Get() *%s {
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
	return &%s{}
}

func (p *%sType) GetSlice(count int) []*%s {
	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Cut the required slice from the end of the pool.
		r := p.pool[len(p.pool)-count:]
		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]
		return r
	}

	// Create a new slice.
	r := make([]*%s, count)

	// Initialize with what remains in the pool.
	i := 0
	for ; i < len(p.pool); i++ {
		r[i] = p.pool[i]
	}
	p.pool = nil

	if i < count {
		// Create remaining elements.
		storage := make([]%s, count-i)
		j := 0
		for ; i < count; i++ {
			r[i] = &storage[j]
			j++
		}
	}

	return r
}
`,
		msg.GetName(), poolName, msg.GetName(), poolName, poolName, poolName,
		msg.GetName(), msg.GetName(), poolName, msg.GetName(), msg.GetName(),
		msg.GetName(),
	)

	g.oPoolReleaseSlice(msg)
	g.oPoolRelease(msg)

	return g.lastErr
}

func getPoolName(msgName string) string {
	return unexportedName(msgName) + "Pool"
}

func (g *generator) oPoolReleaseElem(msg *Message) {
	//poolName := getPoolName(msg.GetName())

	for _, field := range msg.Fields {
		if field.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			continue
		}

		fieldType := getPoolName(field.GetMessageType().GetName())
		g.o("// Release nested %s recursively to their pool.\n", field.GetName())
		if field.IsRepeated() {
			g.o("%s.ReleaseSlice(elem.%s)\n", fieldType, field.GetName())
		} else {
			g.o("if elem.%s != nil {\n", field.GetName())
			g.o("	%s.Release(elem.%s)\n", fieldType, field.GetName())
			g.o("}\n")
		}
	}

	g.o(
		`
// Zero-initialize the released element.
*elem = %s{}
`, msg.GetName(),
	)
}

func (g *generator) oPoolReleaseSlice(msg *Message) {
	poolName := getPoolName(msg.GetName())

	g.o(
		`
// ReleaseSlice releases a slice of elements back to the pool.
func (p *%sType) ReleaseSlice(slice []*%s) {
	for _, elem := range slice {
`, poolName, msg.GetName(),
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
}
`,
	)
}

func (g *generator) oPoolRelease(msg *Message) {
	poolName := getPoolName(msg.GetName())

	g.o(
		`
// Release an element back to the pool.
func (p *%sType) Release(elem *%s) {
`, poolName, msg.GetName(),
	)

	g.i(1)
	g.oPoolReleaseElem(msg)
	g.i(-1)

	g.o(
		`
	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}
`,
	)
}

func preparedFieldDecl(msg *Message, field *Field, typeName string) string {
	return fmt.Sprintf(
		"var prepared%s%s = molecule.Prepare%sField(%d)\n", msg.GetName(),
		field.GetCapitalName(), typeName, field.GetNumber(),
	)
}
