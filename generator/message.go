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
	Fields    []*Field
	FieldsMap map[string]*Field
}

func NewMessage(descr *desc.MessageDescriptor) *Message {
	m := &Message{
		MessageDescriptor: *descr,
		FieldsMap:         map[string]*Field{},
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
