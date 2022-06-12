package generator

import (
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	_ "github.com/jhump/protoreflect/desc/protoparse"
)

func (g *generator) oFieldsAccessors() error {
	for _, field := range g.msg.Fields {
		g.setField(field)

		if g.options.WithPresence {
			if err := g.oHasMethod(); err != nil {
				return err
			}
		}

		if err := g.oFieldGetter(); err != nil {
			return err
		}

		if err := g.oFieldSetter(); err != nil {
			return err
		}

		if g.field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE && g.field.IsRepeated() {
			if err := g.oFieldSliceMethods(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *generator) oFieldGetter() error {
	g.o(`// $FieldName returns the value of the $fieldName.`)

	if g.field.GetOneOf() != nil {
		g.o(
			`// If the field "%s" is not set to "$fieldName" then the returned value is undefined.`,
			g.field.GetOneOf().GetName(),
		)
	}

	goType := g.convertTypeToGo(g.field)

	g.o(`func (m *$MessageName) $FieldName() %s {`, goType)

	g.i(1)
	if g.field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		g.o(`if m._flags&%s == 0 {`, g.msg.DecodedFlagName[g.field])
		g.i(1)
		g.o(`// Decode nested message(s).`)
		if g.field.IsRepeated() {
			g.o(`for i := range m.$fieldName {`)
			g.o(`	// TODO: decide how to handle decoding errors.`)
			g.o(`	_ = m.$fieldName[i].decode()`)
			g.o(`}`)
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
				g.o(`$fieldName := m.$fieldName`)
			}

			g.o(`if $fieldName != nil {`)
			g.o(`	// TODO: decide how to handle decoding errors.`)
			g.o(`	_ = $fieldName.decode()`)
			g.o(`}`)

			if g.field.GetOneOf() != nil {
				g.i(-1)
				g.o(`}`)
			}
		}

		g.o(`m._flags |= %s`, g.msg.DecodedFlagName[g.field])

		g.i(-1)
		g.o(`}`)
	}

	if g.field.GetOneOf() != nil {
		switch g.field.GetType() {
		case descriptor.FieldDescriptorProto_TYPE_BOOL:
			g.o(`return m.%s.BoolVal()`, g.field.GetOneOf().GetName())

		case descriptor.FieldDescriptorProto_TYPE_INT64,
			descriptor.FieldDescriptorProto_TYPE_SFIXED64:
			g.o(`return m.%s.Int64Val()`, g.field.GetOneOf().GetName())

		case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
			g.o(`return m.%s.DoubleVal()`, g.field.GetOneOf().GetName())
		case descriptor.FieldDescriptorProto_TYPE_STRING:
			g.o(`return m.%s.StringVal()`, g.field.GetOneOf().GetName())
		case descriptor.FieldDescriptorProto_TYPE_BYTES:
			g.o(`return m.%s.BytesVal()`, g.field.GetOneOf().GetName())
		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			g.o(`return (%s)(m.%s.PtrVal())`, goType, g.field.GetOneOf().GetName())
		default:
			return fmt.Errorf("unsupported oneof field type %v", g.field.GetType())
		}
	} else {
		g.o(`return m.$fieldName`)
	}

	g.i(-1)
	g.o(`}`)
	g.o(``)

	return g.lastErr
}

func (g *generator) oFieldSetter() error {
	g.o(`// Set$FieldName sets the value of the $fieldName.`)

	if g.field.GetOneOf() != nil {
		g.o(
			`// The oneof field "%s" will be set to "$fieldName".`,
			g.field.GetOneOf().GetName(),
		)
	}

	g.o(`func (m *$MessageName) Set$FieldName(v %s) {`, g.convertTypeToGo(g.field))

	if g.field.GetOneOf() != nil {
		choiceName := composeOneOfChoiceName(g.msg, g.field)
		oneofName := g.field.GetOneOf().GetName()

		g.i(1)
		switch g.field.GetType() {
		case descriptor.FieldDescriptorProto_TYPE_BOOL:
			g.o("m.%s = oneof.NewBool(v, int(%s))", oneofName, choiceName)

		case descriptor.FieldDescriptorProto_TYPE_INT64,
			descriptor.FieldDescriptorProto_TYPE_SFIXED64:
			g.o("m.%s = oneof.NewInt64(v, int(%s))", oneofName, choiceName)

		case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
			g.o("m.%s = oneof.NewDouble(v, int(%s))", oneofName, choiceName)

		case descriptor.FieldDescriptorProto_TYPE_STRING:
			g.o("m.%s = oneof.NewString(v, int(%s))", oneofName, choiceName)

		case descriptor.FieldDescriptorProto_TYPE_BYTES:
			g.o("m.%s = oneof.NewBytes(v, int(%s))", oneofName, choiceName)

		case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
			g.o("m.%s = oneof.NewPtr(unsafe.Pointer(v), int(%s))", oneofName, choiceName)

		default:
			return fmt.Errorf("unsupported oneof field type %v", g.field.GetType())
		}
		g.i(-1)
	} else {
		g.o(`	m.$fieldName = v`)
		if g.options.WithPresence {
			if g.field.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
				g.o(`m._flags |= %s`, g.msg.PresenceFlagName[g.field])
			}
		}
	}

	if g.field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		g.o(``)
		g.o(`	// Make sure the field's Parent points to this message.`)
		if g.field.IsRepeated() {
			g.o(`	for _, elem := range m.$fieldName {`)
			g.o(`		elem._protoMessage.Parent = &m._protoMessage`)
			g.o(`	}`)
		} else {
			g.o(`	v._protoMessage.Parent = &m._protoMessage`)
		}
	}
	g.o(``)
	g.o(`	// Mark this message modified, if not already.`)
	g.o(`	m._protoMessage.MarkModified()`)
	g.o(`}`)
	g.o(``)

	return g.lastErr
}

func (g *generator) oHasMethod() error {
	if g.field.GetOneOf() != nil {
		// Has() func is not needed for oneof fields since they have the Type() func.
		return nil
	}

	g.o(`// Has$FieldName returns true if the $fieldName is present.`)
	g.o(`func (m *$MessageName) Has$FieldName() bool {`)

	g.i(1)

	if g.field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		if g.field.IsRepeated() {
			g.o(`return len(m.$fieldName) > 0`)
		} else {
			g.o(`return m.$fieldName != nil`)
		}
	} else {
		if g.field.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			g.o(`return m._flags & %s != 0`, g.msg.PresenceFlagName[g.field])
		}
	}

	g.i(-1)
	g.o(`}`)
	g.o(``)

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
		if m.$fieldName[newLen] != nil {
			m.$fieldName[newLen].Free()
		}
		m.$fieldName[newLen] = m.$fieldName[i]
		m.$fieldName[i] = nil
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
