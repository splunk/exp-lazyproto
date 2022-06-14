package generator

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
)

func unexportedName(name string) string {
	return strings.ToLower(name[0:1]) + name[1:]
}

func (g *generator) oPool() error {
	g.oPoolStruct()
	g.oPoolGetFuncs()
	g.oPoolReleaseSliceFunc()
	g.oPoolReleaseFunc()

	return g.lastErr
}

func (g *generator) oPoolStruct() {
	g.o(
		`
// Pool of $MessageName structs.
type $messagePoolType struct {
	pool []*$MessageName
	mux  sync.Mutex
}

var $messagePool = $messagePoolType{}`,
	)
}

func (g *generator) oPoolGetFuncs() {
	g.o(
		`
// Get one element from the pool. Creates a new element if the pool is empty.
func (p *$messagePoolType) Get() *$MessageName {
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
	return &$MessageName{}
}

func (p *$messagePoolType) GetSlice(r []*$MessageName) {
	// Create a new slice.
	count := len(r)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]$MessageName, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}
}`,
	)
}

func getPoolName(msgName string) string {
	return unexportedName(msgName) + "Pool"
}

func (g *generator) oPoolReleaseElem() {
	for _, field := range g.msg.Fields {
		g.setField(field)

		if field.GetOneOf() != nil {
			// A oneof field.
			fieldIndex := g.calcOneOfFieldIndex()
			if fieldIndex == 0 {
				// We generate all oneof cases when we see the first field. Skip for the rest.
				g.oPoolReleaseOneofField()
			}
		} else if field.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			// Embedded message that is not a oneof field.
			g.o(`// Release nested $fieldName recursively to their pool.`)
			if field.IsRepeated() {
				g.o(`$fieldTypeMessagePool.ReleaseSlice(elem.$fieldName)`)
			} else {
				g.o(`if elem.$fieldName != nil {`)
				g.o(`	$fieldTypeMessagePool.Release(elem.$fieldName)`)
				g.o(`}`)
			}
		}
	}

	g.o(``)

	g.oResetElem()
}

func (g *generator) oPoolReleaseOneofField() {
	typeName := composeOneOfAliasTypeName(g.msg, g.field.GetOneOf())
	g.o(`switch %s(elem.%s.FieldIndex()) {`, typeName, g.field.GetOneOf().GetName())
	for _, choice := range g.field.GetOneOf().GetChoices() {
		choiceField := g.msg.FieldsMap[choice.GetName()]
		g.setField(choiceField)
		if choiceField.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
			choiceName := composeOneOfChoiceName(g.msg, choiceField)
			g.o(`case %s:`, choiceName)
			g.i(1)
			g.o(
				"ptr := (*$FieldMessageTypeName)(elem.%s.PtrVal())",
				g.field.GetOneOf().GetName(),
			)
			g.o(`if ptr != nil {`)
			g.o(`	$fieldTypeMessagePool.Release(ptr)`)
			g.o(`}`)
			g.i(-1)
		}
	}
	g.o(`}`)
}

func (g *generator) oResetElem() {
	g.o(`// Reset the released element.`)
	g.o(`elem._protoMessage = protomessage.ProtoMessage{}`)
	if g.msg.FlagsBitCount > 0 {
		g.o(`elem._flags = 0`)
	}
	for _, field := range g.msg.Fields {
		g.setField(field)

		if field.IsRepeated() {
			g.o(`elem.$fieldName = elem.$fieldName[:0]`)
		} else if field.GetOneOf() != nil {
			idx := g.calcOneOfFieldIndex()
			if idx == 0 {
				g.o(`elem.%s = oneof.NewNone()`, field.GetOneOf().GetName())
			}
		} else {
			// Not a repeated field and not a oneof. Need to assign a zero value.
			var zeroVal string
			switch field.GetType() {
			case descriptor.FieldDescriptorProto_TYPE_BOOL:
				zeroVal = "false"

			case descriptor.FieldDescriptorProto_TYPE_FIXED64,
				descriptor.FieldDescriptorProto_TYPE_SFIXED64,
				descriptor.FieldDescriptorProto_TYPE_UINT64,
				descriptor.FieldDescriptorProto_TYPE_INT64,
				descriptor.FieldDescriptorProto_TYPE_FIXED32,
				descriptor.FieldDescriptorProto_TYPE_SFIXED32,
				descriptor.FieldDescriptorProto_TYPE_SINT32,
				descriptor.FieldDescriptorProto_TYPE_UINT32,
				descriptor.FieldDescriptorProto_TYPE_DOUBLE:
				zeroVal = "0"

			case descriptor.FieldDescriptorProto_TYPE_STRING:
				zeroVal = `""`
			case descriptor.FieldDescriptorProto_TYPE_BYTES:
				zeroVal = "nil"
			case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
				zeroVal = "nil"
			case descriptor.FieldDescriptorProto_TYPE_ENUM:
				zeroVal = g.enumDescrToEnum[field.GetEnumType()].GetName() + "(0)"
			default:
				g.lastErr = fmt.Errorf("unsupported field type %v", field.GetType())
			}

			g.o(`elem.$fieldName = %s`, zeroVal)
		}
	}
}

func (g *generator) oPoolReleaseSliceFunc() {
	g.o(``)
	g.o(
		`
// ReleaseSlice releases a slice of elements back to the pool.
func (p *$messagePoolType) ReleaseSlice(slice []*$MessageName) {
	for _, elem := range slice {`,
	)

	g.i(2)
	g.oPoolReleaseElem()
	g.i(-2)

	g.o(
		`	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}`,
	)
}

func (g *generator) oPoolReleaseFunc() {
	g.o(``)
	g.o(
		`
// Release an element back to the pool.
func (p *$messagePoolType) Release(elem *$MessageName) {`,
	)

	g.i(1)
	g.oPoolReleaseElem()
	g.i(-1)

	g.o(``)
	g.o(
		`
	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}`,
	)
}
