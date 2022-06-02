package generator

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/yoheimuta/go-protoparser/v4"
	"github.com/yoheimuta/go-protoparser/v4/parser"
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
	inF, err := os.Open(inputFilePath)
	if err != nil {
		return err
	}

	ast, err := protoparser.Parse(inF)
	if err != nil {
		return err
	}

	g.file = &File{
		InputFilePath: inputFilePath,
		Messages:      map[string]*Message{},
	}
	for _, elem := range ast.ProtoBody {
		g.lastErr = nil
		elem.Accept(g)
		if g.lastErr != nil {
			return g.lastErr
		}
	}

	return g.oParsed()
}

func (g *generator) oParsed() error {
	fname := path.Base(g.file.InputFilePath) + ".lz.go"
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

	g.o("package %s\n", g.file.PackageName)
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

	for _, msg := range g.file.Messages {
		g.o("type %s struct {\n", msg.Name)
		g.i(1)
		g.o("ProtoMessage\n")
		for _, field := range msg.Fields {
			g.o("%s %s\n", field.FieldName, g.convertType(field))
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
`, msg.Name, msg.Name, msg.Name, msg.Name,
		)

		g.oRepeatedFieldCounts()

		g.o(
			`
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
`,
		)

		g.i(3)
		g.oFieldDecode(msg.Fields)
		g.i(-3)

		g.o(
			`
			}
			return true, nil
		},
	)
}

`,
		)
	}

	return nil
}

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

func (g *generator) convertType(field *parser.Field) string {
	var s string

	if field.IsRepeated {
		s = "[]"
	}

	switch field.Type {
	case "fixed64":
		s += "uint64"
	case "uint32":
		s += "uint32"
	case "string":
		s += "string"
	default:
		s += "*" + field.Type
	}
	return s
}

func (g *generator) oFieldDecode(fields map[string]*parser.Field) string {
	for _, field := range fields {
		g.o("case %s:", field.FieldNumber)
		g.i(1)
		switch field.Type {
		case "string":
			g.o(
				`
v, err := value.AsStringUnsafe()
if err != nil {
	return false, err
}
m.%s = v
`, field.FieldName,
			)
		}
		g.i(-1)
	}

	return ""
}

func (g *generator) oRepeatedFieldCounts() {

}
