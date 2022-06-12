package generator

import (
	"strings"

	"github.com/jhump/protoreflect/desc"
)

// File that is parsed from a source proto.
type File struct {
	Messages map[string]*Message
}

// Message is the representation of a source proto "message".
type Message struct {
	desc.MessageDescriptor

	// FullName is the full unique name of this message (includes parent names).
	FullName string

	// Fields of this message.
	Fields []*Field

	// FieldsMap maps the name of the field to its representation.
	FieldsMap map[string]*Field

	// FlagsTypeAlias is the name of the type alias for the _flags field.
	FlagsTypeAlias string

	// FlagsTypeAlias is the underlying unsigned integer type of the _flags field.
	FlagsUnderlyingType string

	// FlagsBitCount is the number of the bits that the _flags field contains.
	FlagsBitCount int

	// DecodedFlags are definitions for "decoded" flags. Each such flag indicates
	// that a particular field of the message is already decoded.
	DecodedFlags []flagBitDef
	// DecodedFlagName maps the field to its "decoded" const name.
	DecodedFlagName map[*Field]string

	// DecodedFlags are definitions for "present" flags. Each such flag indicates
	// that a particular field of the message is present. Used for "Has" methods.
	PresenceFlags []flagBitDef
	// PresenceFlagName maps the field to its "present" const name.
	PresenceFlagName map[*Field]string
}

type flagBitDef struct {
	// The name of the const.
	flagName string

	// The value of the const.
	maskVal uint64
}

func NewMessage(parent *Message, descr *desc.MessageDescriptor) *Message {
	m := &Message{
		MessageDescriptor: *descr,
		FieldsMap:         map[string]*Field{},
		DecodedFlagName:   map[*Field]string{},
		PresenceFlagName:  map[*Field]string{},
	}

	if parent != nil {
		// This is message declared inside another message.
		// Compose the name using parent's name to make sure the FullName is unique.
		m.FullName = parent.GetName() + "_" + descr.GetName()
	} else {
		// This is a top-level message.
		m.FullName = descr.GetName()
	}

	for _, field := range descr.GetFields() {
		name := camelCase(field.GetName())

		f := &Field{
			FieldDescriptor: *field,
			camelName:       name,
		}
		m.Fields = append(m.Fields, f)
		m.FieldsMap[field.GetName()] = f
	}
	return m
}

func (f *Message) GetName() string {
	return f.FullName
}

// Field is the representation of a source proto field.
type Field struct {
	desc.FieldDescriptor
	camelName string
}

func (f *Field) GetName() string {
	return f.camelName
}

func (f *Field) GetCapitalName() string {
	return capitalCamelCase(f.camelName)
}

func capitalCamelCase(s string) string {
	if s != "" {
		return strings.ToUpper(s[:1]) + s[1:]
	}
	return ""
}

func camelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		p := parts[i]
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, "")
}

// Enum is the representation of a source proto enum.
type Enum struct {
	desc.EnumDescriptor
	FullName string
}

func NewEnum(parent *Message, descr *desc.EnumDescriptor) *Enum {
	e := &Enum{
		EnumDescriptor: *descr,
	}

	if parent != nil {
		e.FullName = parent.GetName() + "_" + descr.GetName()
	} else {
		e.FullName = descr.GetName()
	}
	return e
}

func (e *Enum) GetName() string {
	return e.FullName
}
