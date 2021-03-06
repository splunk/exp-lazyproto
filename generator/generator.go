package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	_ "github.com/jhump/protoreflect/desc/protoparse"
)

type Options struct {
	WithPresence bool
}

func Generate(
	// Path where to look for imported proto files.
	protoImportPath []string,

	// List of input proto files.
	inputProtoFiles []string,

	// Directory to output generated files to.
	outputDir string,

	// Generation options.
	options Options,
) error {
	g := generator{
		protoImportPath:       protoImportPath,
		outputDir:             outputDir,
		options:               options,
		templateData:          map[string]string{},
		messageDescrToMessage: map[*desc.MessageDescriptor]*Message{},
		enumDescrToEnum:       map[*desc.EnumDescriptor]*Enum{},
	}

	for _, f := range inputProtoFiles {
		if err := g.processFile(f); err != nil {
			return err
		}
	}
	return g.lastErr
}

type generator struct {
	protoImportPath []string
	outputDir       string

	options Options

	// Buffer to accumulated generated file.
	outBuf  *bytes.Buffer
	lastErr error
	// Number of tabs to indent output of the o() method.
	indentTabs int

	// List of messages that need to be generated for the current file.
	messagesToGen []*Message
	// Map of all message declared in the current file or any of imported files.
	messageDescrToMessage map[*desc.MessageDescriptor]*Message

	// List of enums that need to be generated for the current file.
	enumsToGen []*Enum
	// Map of all enums declared in the current file or any of imported files.
	enumDescrToEnum map[*desc.EnumDescriptor]*Enum

	// The current message being generated.
	msg *Message
	// The current field being generated.
	field *Field

	// Fields to replace in templates. Key is the field name, value is the value
	// replace the key by.
	templateData map[string]string

	useSizedMarshaler bool
}

func (g *generator) processFile(inputFilePath string) error {
	// Parse the input file.
	p := protoparse.Parser{
		// Accessor is used when the parser needs to read an input file.
		Accessor: func(filename string) (io.ReadCloser, error) {
			// Try all import paths until we find the requested file.
			for _, includePath := range g.protoImportPath {
				f, err := os.Open(path.Join(includePath, filename))
				if err == nil {
					return f, nil
				}
			}
			return nil, fmt.Errorf(
				"file %s not found in paths %v", inputFilePath, g.protoImportPath,
			)
		},
		IncludeSourceCodeInfo: true,
	}

	fmt.Printf("Reading %s\n", inputFilePath)

	fileDescrs, err := p.ParseFiles(inputFilePath)
	if err != nil {
		return err
	}

	for _, fileDescr := range fileDescrs {
		if err := g.oStartFile(fileDescr); err != nil {
			return err
		}

		// List all enums declared in this file.
		g.listAllEnums(nil, fileDescr.GetEnumTypes(), true)

		// List all messages declared in this file.
		g.listAllMessages(nil, fileDescr.GetMessageTypes(), true)

		// List and remember all messages in imported files, but without generating.
		// We need these messages to be available for lookup if they are used as a
		// field type.
		for _, dep := range fileDescr.GetDependencies() {
			g.listAllMessages(nil, dep.GetMessageTypes(), false)
		}

		if err := g.oEnums(); err != nil {
			return err
		}

		if err := g.oMessages(); err != nil {
			return err
		}

		// Generation is done. Format the generated code nicely and write it to the
		// target file.
		if err := g.formatAndWriteToFile(fileDescr); err != nil {
			return err
		}
	}

	return g.lastErr
}

func (g *generator) oEnums() error {
	for _, enum := range g.enumsToGen {
		if err := g.oEnumDecl(enum); err != nil {
			return err
		}
	}
	return g.lastErr
}

func (g *generator) oMessages() error {
	for _, msg := range g.messagesToGen {
		g.setMessage(msg)

		if err := g.prepareMessage(); err != nil {
			return err
		}

		if err := g.oMessage(); err != nil {
			return err
		}
	}
	return g.lastErr
}

func (g *generator) listAllEnums(
	parent *Message, enums []*desc.EnumDescriptor, toGen bool,
) {
	for _, descr := range enums {
		enum := NewEnum(parent, descr)

		if toGen {
			g.enumsToGen = append(g.enumsToGen, enum)
		}
		g.enumDescrToEnum[descr] = enum
	}
}

func (g *generator) listAllMessages(
	parent *Message, msgs []*desc.MessageDescriptor,
	toGen bool,
) {
	for _, descr := range msgs {
		msg := NewMessage(parent, descr)

		if toGen {
			g.messagesToGen = append(g.messagesToGen, msg)
		}
		g.messageDescrToMessage[descr] = msg

		// Add all enums declared inside this message to the list of all enums.
		g.listAllEnums(msg, msg.GetNestedEnumTypes(), toGen)

		// Add all messages declared inside this message to the list of all messages.
		g.listAllMessages(msg, msg.GetNestedMessageTypes(), toGen)
	}
}

func (g *generator) formatAndWriteToFile(fdescr *desc.FileDescriptor) error {
	destFileName := path.Base(strings.TrimSuffix(fdescr.GetName(), ".proto")) + ".pb.go"
	destFileName = path.Join(g.outputDir, destFileName)
	destDir := path.Dir(destFileName)

	fmt.Printf("Generating %s\n", destFileName)

	if err := os.MkdirAll(destDir, 0700); err != nil {
		return err
	}

	var err error
	f, err := os.Create(destFileName)
	if err != nil {
		return err
	}

	// Nicely format the generated Go code.
	goCode, err := format.Source(g.outBuf.Bytes())
	if err != nil {
		// Write unformatted code to have something to look at.
		f.Write(g.outBuf.Bytes())
		// But still return an error.
		return err
	}

	_, err = f.Write(goCode)
	return err
}

func (g *generator) oStartFile(fdescr *desc.FileDescriptor) error {
	g.outBuf = bytes.NewBuffer(nil)

	g.o(
		`
// Code generated by lazyproto. DO NOT EDIT.
// source: %s
`, fdescr.GetName(),
	)

	packagePath := strings.Split(fdescr.GetPackage(), ".")

	if len(packagePath) > 0 {
		packageName := packagePath[len(packagePath)-1]
		g.o(`package %s`, packageName)
		g.o(``)
	}

	g.o(
		`
import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/tigrannajaryan/exp-lazyproto"
	"github.com/tigrannajaryan/exp-lazyproto/internal/protomessage"
	"github.com/tigrannajaryan/exp-lazyproto/internal/oneof"

	"github.com/tigrannajaryan/exp-lazyproto/internal/molecule"
	"github.com/tigrannajaryan/exp-lazyproto/internal/molecule/src/codec"
)

var _ = oneof.OneOf{} // To avoid unused import warning.
var _ = unsafe.Pointer(nil) // To avoid unused import warning.
var _ = fmt.Errorf // To avoid unused import warning.

`,
	)

	if g.useSizedMarshaler {
		g.o(`import "github.com/tigrannajaryan/exp-lazyproto/internal/streams/sizedstream"`)
	}

	return g.lastErr
}

// prepareMessage prepares data that we need during generation of this message.
func (g *generator) prepareMessage() error {
	// Prepare "decoded" flags.
	for _, field := range g.msg.Fields {
		maskVal := uint64(1) << g.msg.FlagsBitCount
		if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			// We need a flag for this field since it is an embedded message.

			// Create a name for the flag.
			flagName := fmt.Sprintf(
				"flags_%s_%s_Decoded", g.msg.GetName(), field.GetCapitalName(),
			)

			// Remember the flag definition.
			g.msg.DecodedFlags = append(
				g.msg.DecodedFlags, flagBitDef{
					flagName: flagName,
					maskVal:  maskVal,
				},
			)
			g.msg.DecodedFlagName[field] = flagName

			// Count the total number of bits allocated in the flags so far.
			g.msg.FlagsBitCount++
		}
	}

	// Prepare "presence" flags.
	for _, field := range g.msg.Fields {
		g.setField(field)

		if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			// Embedded messages don't need a bit for presence since they use
			// non-nil pointer value to indicate presence instead.
			continue
		}

		if g.isOneOfField() {
			// Oneof fields don't need a bit for presence since they use
			// oneof.OneOf's intrinsic ability to indicate presence instead.
			continue
		}

		maskVal := uint64(1) << g.msg.FlagsBitCount

		if g.options.WithPresence {
			// Create a name for the flag.
			flagName := fmt.Sprintf(
				"flags_%s_%s_Present", g.msg.GetName(), field.GetCapitalName(),
			)

			// Remember the flag definition.
			g.msg.PresenceFlags = append(
				g.msg.PresenceFlags, flagBitDef{
					flagName: flagName,
					maskVal:  maskVal,
				},
			)
			g.msg.PresenceFlagName[field] = flagName

			// Count the total number of bits allocated in the flags so far.
			g.msg.FlagsBitCount++
		}
	}

	// Use the smallest unsigned int type that fits the number of the bits
	// we need in the flag field.
	if g.msg.FlagsBitCount > 0 {
		switch {
		case g.msg.FlagsBitCount <= 8:
			g.msg.FlagsUnderlyingType = "uint8"
		case g.msg.FlagsBitCount <= 16:
			g.msg.FlagsUnderlyingType = "uint16"
		case g.msg.FlagsBitCount <= 32:
			g.msg.FlagsUnderlyingType = "uint32"
		case g.msg.FlagsBitCount <= 64:
			g.msg.FlagsUnderlyingType = "uint64"
		default:
			return fmt.Errorf("more than 64 bits flags not supported")
		}

		g.msg.FlagsTypeAlias = fmt.Sprintf("flags_%s", g.msg.GetName())
	}

	return nil
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

	fieldMessage := g.messageDescrToMessage[field.GetMessageType()]

	if fieldMessage != nil {
		g.templateData["$fieldTypeMessagePool"] = getPoolName(fieldMessage.GetName())
		g.templateData["$FieldMessageTypeName"] = fieldMessage.GetName()
	} else {
		g.templateData["$fieldTypeMessagePool"] = "$fieldTypeMessagePool not defined for " + field.GetName()
		g.templateData["$FieldMessageTypeName"] = "$FieldMessageTypeName not defined for " + field.GetName()
	}
}

func (g *generator) o(str string, a ...any) {

	// Delete leading newlines (we use them to format multiline message nicely in this
	// file, but they are not desirable in the output).
	str = strings.TrimLeft(str, "\n")

	// Replace all templated values.
	for k, v := range g.templateData {
		str = strings.ReplaceAll(str, k, v)
	}

	// Substitute provided parameters.
	str = fmt.Sprintf(str, a...)

	// Split into lines.
	strs := strings.Split(str, "\n")
	for i := range strs {
		if strings.TrimSpace(strs[i]) != "" {
			// Prepend with necessary number of tabs.
			strs[i] = strings.Repeat("\t", g.indentTabs) + strs[i]
		}
	}

	// Join into one string again.
	str = strings.Join(strs, "\n")

	// And write to output buffer.
	_, err := io.WriteString(g.outBuf, str+"\n")
	if err != nil {
		g.lastErr = err
	}
}

func (g *generator) i(ofs int) {
	g.indentTabs += ofs
}

func (g *generator) convertTypeToGo(field *Field) string {
	var s string

	if field.IsRepeated() {
		s = "[]"
	}

	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		s += "bool"

	case descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_UINT64:
		s += "uint64"

	case descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_INT64:
		s += "int64"

	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		s += "uint32"

	case descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SINT32:
		s += "int32"

	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		s += "uint32"

	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		s += "float64"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		s += "string"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		s += "[]byte"
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		s += "*" + g.messageDescrToMessage[field.GetMessageType()].GetName()
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		s += g.enumDescrToEnum[field.GetEnumType()].GetName()
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

func (g *generator) oMessage() error {
	g.o(`// ====================== $MessageName message implementation ======================`)
	g.o(``)

	if err := g.oMessageStruct(); err != nil {
		return err
	}

	if err := g.oUnmarshalFunc(); err != nil {
		return err
	}

	if err := g.oFreeMethod(); err != nil {
		return err
	}

	if err := g.oOneOfFields(); err != nil {
		return err
	}

	if err := g.oFlagsDeclAndConsts(); err != nil {
		return err
	}

	if err := g.oFieldsAccessors(); err != nil {
		return err
	}

	if err := g.oValidateFunc(); err != nil {
		return err
	}

	if err := g.oDecodeMethod(); err != nil {
		return err
	}

	if err := g.oMarshalMethod(); err != nil {
		return err
	}

	if err := g.oPool(); err != nil {
		return err
	}

	g.o(``)

	return nil
}

func (g *generator) oComment(comment string) {
	if comment == "" {
		return
	}

	lines := strings.Split(comment, "\n")
	for _, line := range lines {
		g.o(`// %s`, line)
	}
}

func (g *generator) oMessageStruct() error {
	g.oComment(getLeadingComment(g.msg.GetSourceInfo()))

	g.o(`type $MessageName struct {`)
	g.i(1)
	g.o(`_protoMessage protomessage.ProtoMessage`)

	// Add _flags field if any bits are allocated.
	if g.msg.FlagsBitCount > 0 {
		g.o(`_flags %s`, g.msg.FlagsTypeAlias)
	}
	g.o(``)

	// Add fields to the struct.
	first := true
	for _, field := range g.msg.Fields {
		if field.GetOneOf() != nil {
			// Skip oneof fields for now. They will be generated later.
			continue
		}
		g.setField(field)
		comment := getLeadingComment(field.GetSourceInfo())
		if !first && comment != "" {
			g.o(``)
		}
		g.oComment(comment)
		g.o(`$fieldName %s`, g.convertTypeToGo(field))
		first = false
	}

	// Generate oneof fields now.
	for _, oneof := range g.msg.GetOneOfs() {
		g.o(`%s oneof.OneOf`, oneof.GetName())
	}
	g.i(-1)
	g.o(`}`) // struct
	g.o(``)

	return g.lastErr
}

func composeOneOfAliasTypeName(msg *Message, oneof *desc.OneOfDescriptor) string {
	return msg.GetName() + capitalCamelCase(oneof.GetName())
}

func composeOneOfChoiceName(msg *Message, choice *Field) string {
	return msg.GetName() + choice.GetCapitalName()
}

func composeOneOfNoneChoiceName(msg *Message, oneof *desc.OneOfDescriptor) string {
	return msg.GetName() + capitalCamelCase(oneof.GetName()+"None")
}

func (g *generator) oOneOfFields() error {
	for _, oneof := range g.msg.GetOneOfs() {
		g.oOneOfTypeConsts(oneof)
		g.oOneOfSpecialMethods(oneof)
	}

	return g.lastErr
}

func (g *generator) oOneOfTypeConsts(oneof *desc.OneOfDescriptor) {
	typeName := composeOneOfAliasTypeName(g.msg, oneof)
	g.o(
		"// %s defines the possible types for oneof field %q.", typeName,
		oneof.GetName(),
	)
	g.o(`type %s int`, typeName)
	g.o(``)
	g.o(`const (`)
	g.i(1)

	noneChoiceName := composeOneOfNoneChoiceName(g.msg, oneof)
	g.o(`// %s indicates that none of the oneof choices is set.`, noneChoiceName)
	g.o(`%s %s = 0`, noneChoiceName, typeName)

	for i, choice := range oneof.GetChoices() {
		choiceField := g.msg.FieldsMap[choice.GetName()]
		choiceName := composeOneOfChoiceName(g.msg, choiceField)
		g.o(
			"// %s indicates that oneof field %q is set.", choiceName,
			choiceField.GetName(),
		)
		g.o(`%s %s = %d`, choiceName, typeName, i+1)
	}

	g.i(-1)
	g.o(`)`) // const
	g.o(``)
}

func (g *generator) oOneOfSpecialMethods(oneof *desc.OneOfDescriptor) {
	typeName := composeOneOfAliasTypeName(g.msg, oneof)

	// Generate the Type() method.

	funcName := capitalCamelCase(oneof.GetName()) + "Type"
	g.o(
		"// %s returns the type of the current stored oneof %q.", funcName,
		oneof.GetName(),
	)
	g.o(`// To set the type use one of the setters.`)
	g.o(`func (m *$MessageName) %s() %s {`, funcName, typeName)
	g.o(`	return %s(m.%s.FieldIndex())`, typeName, oneof.GetName())
	g.o(`}`)
	g.o(``)

	// Generate the Unset() method.

	funcName = capitalCamelCase(oneof.GetName()) + "Unset"
	g.o(
		"// %s unsets the oneof field %q, so that it contains none of the choices.",
		funcName, oneof.GetName(),
	)
	g.o(`func (m *$MessageName) %s() {`, funcName)
	g.o(`	m.%s = oneof.NewNone()`, oneof.GetName())
	g.o(`}`)
	g.o(``)
}

func (g *generator) oFlagsDeclAndConsts() error {
	// Declare the type alias.
	if len(g.msg.DecodedFlags) > 0 || len(g.msg.PresenceFlags) > 0 {
		g.o(`// %s is the type of the bit flags.`, g.msg.FlagsTypeAlias)
		g.o(`type %s %s`, g.msg.FlagsTypeAlias, g.msg.FlagsUnderlyingType)
	}

	// Declare the "decoded" consts.
	if len(g.msg.DecodedFlags) > 0 {
		g.o(`// Bitmasks that indicate that the particular nested message is decoded.`)
		for _, bitDef := range g.msg.DecodedFlags {
			g.o(
				"const %s %s = 0x%X", bitDef.flagName, g.msg.FlagsTypeAlias,
				bitDef.maskVal,
			)
		}
		g.o(``)
	}

	// Declare the "presence" consts.
	if len(g.msg.PresenceFlags) > 0 {
		g.o(`// Bitmasks that indicate that the particular field is present.`)
		for _, bitDef := range g.msg.PresenceFlags {
			g.o(
				"const %s %s = 0x%X", bitDef.flagName, g.msg.FlagsTypeAlias,
				bitDef.maskVal,
			)
		}
		g.o(``)
	}
	return g.lastErr
}

// calcOneOfFieldIndex calculates the index of the current field within the list of
// all choices of the oneof field to which this field belongs.
// Returns -1 if the current field does not belong to a oneof field.
func (g *generator) calcOneOfFieldIndex() int {
	fieldIdx := -1
	if g.field.GetOneOf() == nil {
		return fieldIdx
	}

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

// Returns true if the current field belongs to a oneof field.
func (g *generator) isOneOfField() bool {
	return g.calcOneOfFieldIndex() >= 0
}

func (g *generator) oEnumDecl(enum *Enum) error {
	// Declare type alias for enum.
	g.oComment(getLeadingComment(enum.GetSourceInfo()))
	g.o(`type %[1]s uint32`, enum.GetName())
	g.o(``)

	// Declare enum consts.
	g.o(`const (`)

	g.i(1)
	for i, value := range enum.GetValues() {
		comment := getLeadingComment(value.GetSourceInfo())
		if i > 0 && comment != "" {
			g.o(``)
		}
		g.oComment(comment)
		g.o(
			"%[1]s_%[2]s %[1]s = %[3]d", enum.GetName(), value.GetName(),
			value.GetNumber(),
		)
	}
	g.i(-1)

	g.o(`)`) // const
	g.o(``)

	return g.lastErr
}

func (g *generator) oFreeMethod() error {
	g.o(
		`
func (m *$MessageName) Free() {
	$messagePool.Release(m)
}
`,
	)
	return g.lastErr
}
