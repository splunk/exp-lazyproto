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
		m.FieldsMap[name] = f
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
	if f.camelName != "" {
		return strings.ToUpper(f.camelName[:1]) + f.camelName[1:]
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
