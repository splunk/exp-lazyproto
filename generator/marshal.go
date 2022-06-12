package generator

import (
	"fmt"
	"sort"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	_ "github.com/jhump/protoreflect/desc/protoparse"
)

func (g *generator) oMarshalMethod() error {
	// Output "prepared" field definitions.
	for _, field := range g.msg.Fields {
		g.setField(field)

		if field.IsRepeated() && field.GetType() != descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			// We don't use "prepared" for primitive repeated fields.
			continue
		}

		g.oPrepareMarshalField(field)
	}

	// Output the Marshal() func.

	g.o(``)
	if g.useSizedMarshaler {
		g.o(`func (m *$MessageName) Marshal(ps *sizedstream.ProtoStream) error {`)
	} else {
		g.o(`func (m *$MessageName) Marshal(ps *molecule.ProtoStream) error {`)
	}
	g.i(1)
	g.o(`if m._protoMessage.IsModified() {`)
	g.i(1)

	// Order the fields by their number to ensure marshaling is done in an
	// order that we can rely on in the tests (same order as other as Protobuf
	// libs so that we can compare the results).
	fields := make([]*Field, len(g.msg.Fields))
	copy(fields, g.msg.Fields)
	sort.Slice(
		fields, func(i, j int) bool {
			return fields[i].GetNumber() < fields[j].GetNumber()
		},
	)

	// Marshal one field at a time.

	for _, field := range fields {
		g.setField(field)

		if field.GetOneOf() != nil {
			fieldIndex := g.calcOneOfFieldIndex()

			// Only marshal oneof field once. We generate the marshaling code
			// when we see the field with index 0 in the list of choices.
			if fieldIndex != 0 {
				// Already generated this oneof, nothing else to do.
				continue
			}

			g.o(`// Marshal %q.`, g.field.GetOneOf().GetName())

			typeName := composeOneOfAliasTypeName(g.msg, g.field.GetOneOf())
			g.o(`// Switch on the type of the value stored in the oneof field.`)
			g.o(`switch %s(m.%s.FieldIndex()) {`, typeName, g.field.GetOneOf().GetName())

			// Add the "none" case.
			noneChoiceName := composeOneOfNoneChoiceName(g.msg, g.field.GetOneOf())
			g.o(`case %s:`, noneChoiceName)
			g.o(`	// Nothing to do, oneof is unset.`)

			// Generate cases for each choice of this oneof.
			for _, choice := range field.GetOneOf().GetChoices() {
				oneofField := g.msg.FieldsMap[choice.GetName()]
				g.setField(oneofField)
				typeName := composeOneOfChoiceName(g.msg, g.field)
				g.o(`case %s:`, typeName)
				g.i(1)
				g.oMarshalField()
				g.i(-1)
			}
			g.o(`}`)
		} else {
			g.oMarshalField()
		}
	}
	g.i(-1)
	g.o(`} else {`)
	g.o(`	// We have the original bytes and the message is unchanged. Use the original bytes.`)
	g.o(`	ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))`)
	g.o(`}`)
	g.o(`return nil`)
	g.i(-1)
	g.o(`}`)
	return g.lastErr
}

// Returns the name of the var which contains prepared data for the specified field.
func embeddedFieldPreparedVarName(msg *Message, field *Field) string {
	return fmt.Sprintf("prepared_%s_%s", msg.GetName(), field.GetCapitalName())
}

func convertProtoTypeToOneOfType(protoType string) string {
	switch protoType {
	case "SFixed64":
		return "Int64"
	case "SFixed32":
		return "Int64"
	default:
		return protoType
	}
}

func (g *generator) oMarshalPreparedField(protoTypeName string) {
	if g.field.GetOneOf() != nil {
		g.o(
			"ps.%[1]sPrepared(prepared_$MessageName_$FieldName, m.%[2]s.%[3]sVal())",
			protoTypeName,
			g.field.GetOneOf().GetName(),
			convertProtoTypeToOneOfType(protoTypeName),
		)
	} else {
		g.o(
			"ps.%sPrepared(prepared_$MessageName_$FieldName, m.$fieldName)",
			protoTypeName,
		)
	}
}

func (g *generator) oMarshalField() {
	g.o(`// Marshal "$fieldName".`)

	if g.field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
		g.oMarshalMessageTypeField()
		return
	} else if g.field.IsRepeated() {
		g.oMarshalPrimitiveRepeated()
		return
	}

	usePresence := g.options.WithPresence && !g.isOneOfField()

	if usePresence {
		g.o(`if m._flags&%s != 0 {`, g.msg.PresenceFlagName[g.field])
		g.i(1)
	}

	switch g.field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		g.oMarshalPreparedField("Bool")

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		g.oMarshalPreparedField("String")

	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		g.oMarshalPreparedField("Bytes")

	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		g.oMarshalPreparedField("Fixed64")

	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		g.oMarshalPreparedField("SFixed64")

	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		g.oMarshalPreparedField("Uint64")

	case descriptor.FieldDescriptorProto_TYPE_INT64:
		g.oMarshalPreparedField("Int64")

	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		g.oMarshalPreparedField("Fixed32")

	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		g.oMarshalPreparedField("Uint32")

	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		g.oMarshalPreparedField("Sint32")

	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		g.oMarshalPreparedField("Double")

	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		g.o(`ps.Uint32Prepared(prepared_$MessageName_$FieldName, uint32(m.$fieldName))`)

	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:

	default:
		g.lastErr = fmt.Errorf("unsupported field type %v", g.field.GetType())
	}

	if usePresence {
		g.i(-1)
		g.o(`}`)
	}
}

func (g *generator) oMarshalPrimitiveRepeated() {
	switch g.field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		g.oMarshalPrimitiveRepeatedField("Bool")

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		g.oMarshalPrimitiveRepeatedField("String")

	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		g.oMarshalPrimitiveRepeatedField("Bytes")

	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		g.oMarshalPrimitiveRepeatedField("Fixed64")

	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		g.oMarshalPrimitiveRepeatedField("SFixed64")

	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		g.oMarshalPrimitiveRepeatedField("Uint64")

	case descriptor.FieldDescriptorProto_TYPE_INT64:
		g.oMarshalPrimitiveRepeatedField("Int64")

	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		g.oMarshalPrimitiveRepeatedField("Fixed32")

	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		g.oMarshalPrimitiveRepeatedField("Uint32")

	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		g.oMarshalPrimitiveRepeatedField("Sint32")

	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		g.oMarshalPrimitiveRepeatedField("Double")

	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		g.o(`ps.Uint32Prepared(prepared_$MessageName_$FieldName, uint32(m.$fieldName))`)

	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:

	default:
		g.lastErr = fmt.Errorf("unsupported field type %v", g.field.GetType())
	}
}

func (g *generator) oMarshalPrimitiveRepeatedField(protoTypeName string) {
	g.o(`ps.%sPacked(%d, m.$fieldName)`, protoTypeName, g.field.GetNumber())
}

func (g *generator) oMarshalMessageTypeField() {
	if g.field.IsRepeated() {
		g.o(`for _, elem := range m.$fieldName {`)
		g.o(`	token := ps.BeginEmbedded()`)
		g.o(`	if err := elem.Marshal(ps); err != nil {`)
		g.o(`		return err`)
		g.o(`	}`)
		g.o(
			"	ps.EndEmbeddedPrepared(token, %s)",
			embeddedFieldPreparedVarName(g.msg, g.field),
		)
	} else {
		if g.field.GetOneOf() != nil {
			g.o(
				"$fieldName := (*$FieldMessageTypeName)(m.%s.PtrVal())",
				g.field.GetOneOf().GetName(),
			)
		} else {
			g.o(`$fieldName := m.$fieldName`)
		}

		g.o(`if $fieldName != nil {`)
		g.o(`	token := ps.BeginEmbedded()`)
		g.o(`	if err := $fieldName.Marshal(ps); err != nil {`)
		g.o(`		return err`)
		g.o(`	}`)
		g.o(
			"	ps.EndEmbeddedPrepared(token, %s)",
			embeddedFieldPreparedVarName(g.msg, g.field),
		)
	}
	g.o(`}`)
}

func (g *generator) oPrepareMarshalField(field *Field) {
	switch field.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		g.o(g.preparedFieldDecl(g.msg, field, "Bool"))

	case descriptor.FieldDescriptorProto_TYPE_STRING:
		g.o(g.preparedFieldDecl(g.msg, field, "String"))

	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		g.o(g.preparedFieldDecl(g.msg, field, "Bytes"))

	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		g.o(g.preparedFieldDecl(g.msg, field, "Fixed64"))

	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		g.o(g.preparedFieldDecl(g.msg, field, "Fixed64"))

	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		g.o(g.preparedFieldDecl(g.msg, field, "Uint64"))

	case descriptor.FieldDescriptorProto_TYPE_INT64:
		g.o(g.preparedFieldDecl(g.msg, field, "Int64"))

	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		g.o(g.preparedFieldDecl(g.msg, field, "Fixed32"))

	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		g.o(g.preparedFieldDecl(g.msg, field, "Uint32"))

	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		g.o(g.preparedFieldDecl(g.msg, field, "Sint32"))

	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		g.o(g.preparedFieldDecl(g.msg, field, "Double"))

	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		g.o(g.preparedFieldDecl(g.msg, field, "Uint32"))

	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		g.o(g.preparedFieldDecl(g.msg, field, "Embedded"))

	default:
		g.lastErr = fmt.Errorf("unsupported field type %v", field.GetType())
	}
}

func (g *generator) preparedFieldDecl(
	msg *Message, field *Field, prepareFuncNamePrefix string,
) string {
	if g.useSizedMarshaler {
		return fmt.Sprintf(
			"var prepared_%s_%s = sizedstream.Prepare%sField(%d)", msg.GetName(),
			field.GetCapitalName(), prepareFuncNamePrefix, field.GetNumber(),
		)
	} else {
		return fmt.Sprintf(
			"var prepared_%s_%s = molecule.Prepare%sField(%d)", msg.GetName(),
			field.GetCapitalName(), prepareFuncNamePrefix, field.GetNumber(),
		)
	}
}
