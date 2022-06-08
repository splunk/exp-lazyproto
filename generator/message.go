package generator

import (
	"strings"

	"github.com/jhump/protoreflect/desc"
)

type File struct {
	Messages map[string]*Message
}

type Message struct {
	desc.MessageDescriptor

	Parent   *Message
	FullName string

	Fields    []*Field
	FieldsMap map[string]*Field

	FlagsType           string
	FlagsBitCount       int
	HasEmbeddedMessages bool
	FlagsUnderlyingType string

	DecodedFlags     []flagBitDef
	PresenceFlags    []flagBitDef
	DecodedFlagName  map[*Field]string
	PresenceFlagName map[*Field]string
}

type flagBitDef struct {
	flagName string
	maskVal  uint64
}

func NewMessage(parent *Message, descr *desc.MessageDescriptor) *Message {
	m := &Message{
		MessageDescriptor: *descr,
		FieldsMap:         map[string]*Field{},
		DecodedFlagName:   map[*Field]string{},
		PresenceFlagName:  map[*Field]string{},
	}

	if parent != nil {
		m.FullName = parent.GetName() + "_" + descr.GetName()
	} else {
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
