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

func (g *generator) VisitField(field *parser.Field) (next bool) {
	g.field = field
	g.msg.Fields[field.FieldName] = g.field
	return true
}

func (g *generator) VisitGroupField(field *parser.GroupField) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) VisitImport(i *parser.Import) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) VisitMapField(field *parser.MapField) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) VisitMessage(message *parser.Message) (next bool) {
	g.msg = &Message{
		Name:   message.MessageName,
		Fields: map[string]*parser.Field{},
	}
	g.file.Messages[message.MessageName] = g.msg
	return true
}

func (g *generator) VisitOneof(oneof *parser.Oneof) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) VisitOneofField(field *parser.OneofField) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) VisitOption(option *parser.Option) (next bool) {
	//TODO implement me
	return true
}

func (g *generator) VisitPackage(p *parser.Package) (next bool) {
	g.file.PackageName = p.Name
	return true
}

func (g *generator) VisitReserved(reserved *parser.Reserved) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) VisitRPC(rpc *parser.RPC) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) VisitService(service *parser.Service) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) VisitSyntax(syntax *parser.Syntax) (next bool) {
	//TODO implement me
	panic("implement me")
}

func (g *generator) processFile(inputFilePath string) error {
	p := protoparse.Parser{
		Accessor: func(filename string) (io.ReadCloser, error) {
			return os.Open(inputFilePath)
		},
	}
	fdescrs, err := p.ParseFiles(inputFilePath)
	if err != nil {
		return err
	}

	for _, fdescr := range fdescrs {
		if err := g.oFile(fdescr); err != nil {
			return err
		}
		for _, msg := range fdescr.GetMessageTypes() {
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
	"github.com/richardartoul/molecule"
	"github.com/richardartoul/molecule/src/codec"
)
`,
	)

	g.o(
		`
const flagsMessageModified = 1

type ProtoMessage struct {
	flags  uint64
	parent *ProtoMessage
	bytes []byte
}

`,
	)

	return nil
}

//func (g *generator) oParsed() error {
//	fname := path.Base(g.file.InputFilePath) + ".lz.go"
//	fname = path.Join(g.outputDir, fname)
//	fdir := path.Dir(fname)
//	if err := os.MkdirAll(fdir, 0700); err != nil {
//		return err
//	}
//
//	var err error
//	g.outF, err = os.Create(fname)
//	if err != nil {
//		return err
//	}
//
//	g.o("package %s\n", g.file.PackageName)
//	g.o(
//		`
//import (
//	"github.com/richardartoul/molecule"
//	"github.com/richardartoul/molecule/src/codec"
//)
//`,
//	)
//
//	g.o(
//		`
//const flagsMessageModified = 1
//type ProtoMessage struct {
//	flags  uint64
//	parent *ProtoMessage
//	bytes []byte
//}
//`,
//	)
//
//	for _, msg := range g.file.Messages {
//		g.o("type %s struct {\n", msg.Name)
//		g.i(1)
//		g.o("ProtoMessage\n")
//		//for _, field := range msg.Fields {
//		//	g.o("%s %s\n", field.FieldName, g.convertType(field))
//		//}
//		g.i(-1)
//		g.o("}\n")
//
//		g.o(
//			`
//func New%s(bytes []byte) *%s {
//	m := &%s{ProtoMessage: ProtoMessage{bytes: bytes}}
//	m.decode()
//	return m
//}
//
//func (m *%s) decode() {
//	buf := codec.NewBuffer(m.bytes)
//`, msg.Name, msg.Name, msg.Name, msg.Name,
//		)
//
//		g.oRepeatedFieldCounts()
//
//		g.o(
//			`
//	molecule.MessageEach(
//		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
//			switch fieldNum {
//`,
//		)
//
//		g.i(3)
//		//g.oFieldDecode(msg.Fields)
//		g.i(-3)
//
//		g.o(
//			`
//			}
//			return true, nil
//		},
//	)
//}
//
//`,
//		)
//	}
//
//	return nil
//}

func (g *generator) o(str string, a ...any) bool {

	str = fmt.Sprintf(str, a...)

	strs := strings.Split(str, "\n")
	for i := range strs {
		if strings.TrimSpace(strs[i]) != "" {
			strs[i] = strings.Repeat(" ", g.spaces) + strs[i]
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
	g.spaces += 4 * ofs
}

func (g *generator) convertType(field *desc.FieldDescriptor) string {
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

func (g *generator) oMessage(msg *desc.MessageDescriptor) error {
	g.o("type %s struct {\n", msg.GetName())
	g.i(1)
	g.o("ProtoMessage\n")
	for _, field := range msg.GetFields() {
		g.o("%s %s\n", field.GetName(), g.convertType(field))
	}
	g.i(-1)
	g.o("}\n")

	g.o(
		`
func New%s(bytes []byte) *%s {
	m := &%s{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *%s) decode() {
	buf := codec.NewBuffer(m.bytes)
`, msg.GetName(), msg.GetName(), msg.GetName(), msg.GetName(),
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
	g.oFieldDecode(msg.GetFields())
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

func (g *generator) oFieldDecodePrimitive(field *desc.FieldDescriptor, asType string) {
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

func (g *generator) oFieldDecode(fields []*desc.FieldDescriptor) string {
	for _, field := range fields {
		g.o("case %d:\n", field.GetNumber())
		g.i(1)
		g.o("// Decode %s", field.GetName())
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
	ProtoMessage: ProtoMessage{bytes: v, parent: &m.ProtoMessage},
}
%s++
`, field.GetName(), counterName, field.GetMessageType().GetName(), counterName,
				)
			} else {
				g.o(
					`
m.%s = &%s{	
	ProtoMessage: ProtoMessage{bytes: v, parent: &m.ProtoMessage},
}
`, field.GetName(), field.GetMessageType().GetName(),
				)
			}
		}
		g.i(-1)
	}

	return ""
}

func (g *generator) oRepeatedFieldCounts(msg *desc.MessageDescriptor) {
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
			"m.%s = make([]*%s, 0, %s)\n", field.GetName(),
			field.GetMessageType().GetName(),
			counterName,
		)
	}
	g.o("\n// Reset the buffer to start iterating over the fields again")
	g.o("\nbuf.Reset(m.bytes)\n")
	g.o("\n// Set slice indexes to 0 to begin iterating over repeated fields.\n")
	for _, field := range fields {
		counterName := field.GetName() + "Count"
		g.o("%s = 0", counterName)
	}
}

func (g *generator) getRepeatedFields(msg *desc.MessageDescriptor) []*desc.FieldDescriptor {
	var r []*desc.FieldDescriptor
	for _, field := range msg.GetFields() {
		if field.IsRepeated() {
			r = append(r, field)
		}
	}
	return r
}
