package simple

import (
	"sync"
	"unsafe"

	"github.com/tigrannajaryan/lazyproto/internal/oneof"
	"github.com/tigrannajaryan/lazyproto/internal/protomessage"

	"github.com/tigrannajaryan/molecule"
	"github.com/tigrannajaryan/molecule/src/codec"
)

// SeverityNumber values
type SeverityNumber uint32

const (
	// SeverityNumber is not specified
	SeverityNumber_SEVERITY_NUMBER_UNSPECIFIED SeverityNumber = 0
	SeverityNumber_SEVERITY_NUMBER_TRACE       SeverityNumber = 1
	SeverityNumber_SEVERITY_NUMBER_TRACE2      SeverityNumber = 2
	SeverityNumber_SEVERITY_NUMBER_TRACE3      SeverityNumber = 3
	SeverityNumber_SEVERITY_NUMBER_TRACE4      SeverityNumber = 4
	SeverityNumber_SEVERITY_NUMBER_DEBUG       SeverityNumber = 5
	SeverityNumber_SEVERITY_NUMBER_DEBUG2      SeverityNumber = 6
	SeverityNumber_SEVERITY_NUMBER_DEBUG3      SeverityNumber = 7
	SeverityNumber_SEVERITY_NUMBER_DEBUG4      SeverityNumber = 8
	SeverityNumber_SEVERITY_NUMBER_INFO        SeverityNumber = 9
	SeverityNumber_SEVERITY_NUMBER_INFO2       SeverityNumber = 10
	SeverityNumber_SEVERITY_NUMBER_INFO3       SeverityNumber = 11
	SeverityNumber_SEVERITY_NUMBER_INFO4       SeverityNumber = 12
	SeverityNumber_SEVERITY_NUMBER_WARN        SeverityNumber = 13
	SeverityNumber_SEVERITY_NUMBER_WARN2       SeverityNumber = 14
	SeverityNumber_SEVERITY_NUMBER_WARN3       SeverityNumber = 15
	SeverityNumber_SEVERITY_NUMBER_WARN4       SeverityNumber = 16
	SeverityNumber_SEVERITY_NUMBER_ERROR       SeverityNumber = 17
	SeverityNumber_SEVERITY_NUMBER_ERROR2      SeverityNumber = 18
	SeverityNumber_SEVERITY_NUMBER_ERROR3      SeverityNumber = 19
	SeverityNumber_SEVERITY_NUMBER_ERROR4      SeverityNumber = 20
	SeverityNumber_SEVERITY_NUMBER_FATAL       SeverityNumber = 21
	SeverityNumber_SEVERITY_NUMBER_FATAL2      SeverityNumber = 22
	SeverityNumber_SEVERITY_NUMBER_FATAL3      SeverityNumber = 23
	SeverityNumber_SEVERITY_NUMBER_FATAL4      SeverityNumber = 24
)

// ====================== LogsData message implementation ======================

// LogsData contains all log data
type LogsData struct {
	_protoMessage protomessage.ProtoMessage
	_flags        LogsDataFlags

	// List of ResourceLogs
	resourceLogs []*ResourceLogs
}

func UnmarshalLogsData(bytes []byte) (*LogsData, error) {
	m := logsDataPool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *LogsData) Free() {
	logsDataPool.Release(m)
}

// LogsDataFlags is the type of the bit flags.
type LogsDataFlags uint8

// Bitmasks that indicate that the particular nested message is decoded.
const flagLogsDataResourceLogsDecoded LogsDataFlags = 0x1

// HasResourceLogs returns true if the resourceLogs is present.
func (m *LogsData) HasResourceLogs() bool {
	return len(m.resourceLogs) > 0
}

// ResourceLogs returns the value of the resourceLogs.
func (m *LogsData) ResourceLogs() []*ResourceLogs {
	if m._flags&flagLogsDataResourceLogsDecoded == 0 {
		// Decode nested message(s).
		for i := range m.resourceLogs {
			// TODO: decide how to handle decoding errors.
			_ = m.resourceLogs[i].decode()
		}
		m._flags |= flagLogsDataResourceLogsDecoded
	}
	return m.resourceLogs
}

// SetResourceLogs sets the value of the resourceLogs.
func (m *LogsData) SetResourceLogs(v []*ResourceLogs) {
	m.resourceLogs = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.resourceLogs {
		elem._protoMessage.Parent = &m._protoMessage
	}

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *LogsData) ResourceLogsRemoveIf(f func(*ResourceLogs) bool) {
	// Call getter to load the field.
	m.ResourceLogs()

	newLen := 0
	for i := 0; i < len(m.resourceLogs); i++ {
		if f(m.resourceLogs[i]) {
			continue
		}
		if newLen == i {
			// Nothing to move, element is at the right place.
			newLen++
			continue
		}
		m.resourceLogs[newLen] = m.resourceLogs[i]
		newLen++
	}
	if newLen != len(m.resourceLogs) {
		m.resourceLogs = m.resourceLogs[:newLen]
		// Mark this message modified, if not already.
		m._protoMessage.MarkModified()
	}
}

func (m *LogsData) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// If the user makes a mistake and takes a copy of this struct before decoding it
	// the "decoded" flag will incorrectly set on the copy, not the original, but nested
	// messages will be marked as "decoded". Next time we try to access the nested
	// message via getter func this decode() func will be called again and will overwrite
	// nested message values, but the getter func will no longer attempt to decode it
	// because the flag "decoded" flag is incorrectly set on nested message.
	// This will result in incorrect state of nested message returned by getter.
	// To make sure we correctly decode even after this mistake we reset all "decoded"
	// and "presence" flags here.
	m._flags = 0

	// Count all repeated fields. We need one counter per field.
	resourceLogsCount := 0
	err := molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 1 {
				resourceLogsCount++
			}
		},
	)
	if err != nil {
		return err
	}

	// Pre-allocate slices for repeated fields.
	m.resourceLogs = resourceLogsPool.GetSlice(resourceLogsCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// Set slice indexes to 0 to begin iterating over repeated fields.
	resourceLogsCount = 0
	// Iterate and decode the fields.
	err2 := molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode "resourceLogs".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.resourceLogs[resourceLogsCount]
				resourceLogsCount++
				elem._protoMessage.Parent = &m._protoMessage
				elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
			}
			return true, nil
		},
	)
	if err2 != nil {
		return err2
	}
	return nil
}

var preparedLogsDataResourceLogs = molecule.PrepareEmbeddedField(1)

func (m *LogsData) Marshal(ps *molecule.ProtoStream) error {
	if m._protoMessage.IsModified() {
		// Marshal "resourceLogs".
		for _, elem := range m.resourceLogs {
			token := ps.BeginEmbedded()
			if err := elem.Marshal(ps); err != nil {
				return err
			}
			ps.EndEmbeddedPrepared(token, preparedLogsDataResourceLogs)
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))
	}
	return nil
}

// Pool of LogsData structs.
type logsDataPoolType struct {
	pool []*LogsData
	mux  sync.Mutex
}

var logsDataPool = logsDataPoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *logsDataPoolType) Get() *LogsData {
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
	return &LogsData{}
}

func (p *logsDataPoolType) GetSlice(count int) []*LogsData {
	// Create a new slice.
	r := make([]*LogsData, count)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return r
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]LogsData, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *logsDataPoolType) ReleaseSlice(slice []*LogsData) {
	for _, elem := range slice {
		// Release nested resourceLogs recursively to their pool.
		resourceLogsPool.ReleaseSlice(elem.resourceLogs)

		// Zero-initialize the released element.
		*elem = LogsData{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *logsDataPoolType) Release(elem *LogsData) {
	// Release nested resourceLogs recursively to their pool.
	resourceLogsPool.ReleaseSlice(elem.resourceLogs)

	// Zero-initialize the released element.
	*elem = LogsData{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}

// ====================== ResourceLogs message implementation ======================

type ResourceLogs struct {
	_protoMessage protomessage.ProtoMessage
	_flags        ResourceLogsFlags

	// The Resource
	resource *Resource
	// List of ScopeLogs
	scopeLogs []*ScopeLogs
	schemaUrl string
}

func UnmarshalResourceLogs(bytes []byte) (*ResourceLogs, error) {
	m := resourceLogsPool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *ResourceLogs) Free() {
	resourceLogsPool.Release(m)
}

// ResourceLogsFlags is the type of the bit flags.
type ResourceLogsFlags uint8

// Bitmasks that indicate that the particular nested message is decoded.
const flagResourceLogsResourceDecoded ResourceLogsFlags = 0x1
const flagResourceLogsScopeLogsDecoded ResourceLogsFlags = 0x2

// Bitmasks that indicate that the particular field is present.
const flagResourceLogsSchemaUrlPresent ResourceLogsFlags = 0x4

// HasResource returns true if the resource is present.
func (m *ResourceLogs) HasResource() bool {
	return m.resource != nil
}

// Resource returns the value of the resource.
func (m *ResourceLogs) Resource() *Resource {
	if m._flags&flagResourceLogsResourceDecoded == 0 {
		// Decode nested message(s).
		resource := m.resource
		if resource != nil {
			// TODO: decide how to handle decoding errors.
			_ = resource.decode()
		}
		m._flags |= flagResourceLogsResourceDecoded
	}
	return m.resource
}

// SetResource sets the value of the resource.
func (m *ResourceLogs) SetResource(v *Resource) {
	m.resource = v

	// Make sure the field's Parent points to this message.
	v._protoMessage.Parent = &m._protoMessage

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasScopeLogs returns true if the scopeLogs is present.
func (m *ResourceLogs) HasScopeLogs() bool {
	return len(m.scopeLogs) > 0
}

// ScopeLogs returns the value of the scopeLogs.
func (m *ResourceLogs) ScopeLogs() []*ScopeLogs {
	if m._flags&flagResourceLogsScopeLogsDecoded == 0 {
		// Decode nested message(s).
		for i := range m.scopeLogs {
			// TODO: decide how to handle decoding errors.
			_ = m.scopeLogs[i].decode()
		}
		m._flags |= flagResourceLogsScopeLogsDecoded
	}
	return m.scopeLogs
}

// SetScopeLogs sets the value of the scopeLogs.
func (m *ResourceLogs) SetScopeLogs(v []*ScopeLogs) {
	m.scopeLogs = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.scopeLogs {
		elem._protoMessage.Parent = &m._protoMessage
	}

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *ResourceLogs) ScopeLogsRemoveIf(f func(*ScopeLogs) bool) {
	// Call getter to load the field.
	m.ScopeLogs()

	newLen := 0
	for i := 0; i < len(m.scopeLogs); i++ {
		if f(m.scopeLogs[i]) {
			continue
		}
		if newLen == i {
			// Nothing to move, element is at the right place.
			newLen++
			continue
		}
		m.scopeLogs[newLen] = m.scopeLogs[i]
		newLen++
	}
	if newLen != len(m.scopeLogs) {
		m.scopeLogs = m.scopeLogs[:newLen]
		// Mark this message modified, if not already.
		m._protoMessage.MarkModified()
	}
}

// HasSchemaUrl returns true if the schemaUrl is present.
func (m *ResourceLogs) HasSchemaUrl() bool {
	return m._flags&flagResourceLogsSchemaUrlPresent != 0
}

// SchemaUrl returns the value of the schemaUrl.
func (m *ResourceLogs) SchemaUrl() string {
	return m.schemaUrl
}

// SetSchemaUrl sets the value of the schemaUrl.
func (m *ResourceLogs) SetSchemaUrl(v string) {
	m.schemaUrl = v
	m._flags |= flagResourceLogsSchemaUrlPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *ResourceLogs) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// If the user makes a mistake and takes a copy of this struct before decoding it
	// the "decoded" flag will incorrectly set on the copy, not the original, but nested
	// messages will be marked as "decoded". Next time we try to access the nested
	// message via getter func this decode() func will be called again and will overwrite
	// nested message values, but the getter func will no longer attempt to decode it
	// because the flag "decoded" flag is incorrectly set on nested message.
	// This will result in incorrect state of nested message returned by getter.
	// To make sure we correctly decode even after this mistake we reset all "decoded"
	// and "presence" flags here.
	m._flags = 0

	// Count all repeated fields. We need one counter per field.
	scopeLogsCount := 0
	err := molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 2 {
				scopeLogsCount++
			}
		},
	)
	if err != nil {
		return err
	}

	// Pre-allocate slices for repeated fields.
	m.scopeLogs = scopeLogsPool.GetSlice(scopeLogsCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// Set slice indexes to 0 to begin iterating over repeated fields.
	scopeLogsCount = 0
	// Iterate and decode the fields.
	err2 := molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode "resource".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.resource = resourcePool.Get()
				m.resource._protoMessage.Parent = &m._protoMessage
				m.resource._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
			case 2:
				// Decode "scopeLogs".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.scopeLogs[scopeLogsCount]
				scopeLogsCount++
				elem._protoMessage.Parent = &m._protoMessage
				elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
			case 3:
				// Decode "schemaUrl".
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.schemaUrl = v
				m._flags |= flagResourceLogsSchemaUrlPresent
			}
			return true, nil
		},
	)
	if err2 != nil {
		return err2
	}
	return nil
}

var preparedResourceLogsResource = molecule.PrepareEmbeddedField(1)
var preparedResourceLogsScopeLogs = molecule.PrepareEmbeddedField(2)
var preparedResourceLogsSchemaUrl = molecule.PrepareStringField(3)

func (m *ResourceLogs) Marshal(ps *molecule.ProtoStream) error {
	if m._protoMessage.IsModified() {
		// Marshal "resource".
		elem := m.resource
		if elem != nil {
			token := ps.BeginEmbedded()
			if err := elem.Marshal(ps); err != nil {
				return err
			}
			ps.EndEmbeddedPrepared(token, preparedResourceLogsResource)
		}
		// Marshal "scopeLogs".
		for _, elem := range m.scopeLogs {
			token := ps.BeginEmbedded()
			if err := elem.Marshal(ps); err != nil {
				return err
			}
			ps.EndEmbeddedPrepared(token, preparedResourceLogsScopeLogs)
		}
		// Marshal "schemaUrl".
		if m._flags&flagResourceLogsSchemaUrlPresent != 0 {
			ps.StringPrepared(preparedResourceLogsSchemaUrl, m.schemaUrl)
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))
	}
	return nil
}

// Pool of ResourceLogs structs.
type resourceLogsPoolType struct {
	pool []*ResourceLogs
	mux  sync.Mutex
}

var resourceLogsPool = resourceLogsPoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *resourceLogsPoolType) Get() *ResourceLogs {
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
	return &ResourceLogs{}
}

func (p *resourceLogsPoolType) GetSlice(count int) []*ResourceLogs {
	// Create a new slice.
	r := make([]*ResourceLogs, count)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return r
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]ResourceLogs, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *resourceLogsPoolType) ReleaseSlice(slice []*ResourceLogs) {
	for _, elem := range slice {
		// Release nested resource recursively to their pool.
		if elem.resource != nil {
			resourcePool.Release(elem.resource)
		}
		// Release nested scopeLogs recursively to their pool.
		scopeLogsPool.ReleaseSlice(elem.scopeLogs)

		// Zero-initialize the released element.
		*elem = ResourceLogs{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *resourceLogsPoolType) Release(elem *ResourceLogs) {
	// Release nested resource recursively to their pool.
	if elem.resource != nil {
		resourcePool.Release(elem.resource)
	}
	// Release nested scopeLogs recursively to their pool.
	scopeLogsPool.ReleaseSlice(elem.scopeLogs)

	// Zero-initialize the released element.
	*elem = ResourceLogs{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}

// ====================== Resource message implementation ======================

type Resource struct {
	_protoMessage protomessage.ProtoMessage
	_flags        ResourceFlags

	attributes             []*KeyValue
	droppedAttributesCount uint32
}

func UnmarshalResource(bytes []byte) (*Resource, error) {
	m := resourcePool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Resource) Free() {
	resourcePool.Release(m)
}

// ResourceFlags is the type of the bit flags.
type ResourceFlags uint8

// Bitmasks that indicate that the particular nested message is decoded.
const flagResourceAttributesDecoded ResourceFlags = 0x1

// Bitmasks that indicate that the particular field is present.
const flagResourceDroppedAttributesCountPresent ResourceFlags = 0x2

// HasAttributes returns true if the attributes is present.
func (m *Resource) HasAttributes() bool {
	return len(m.attributes) > 0
}

// Attributes returns the value of the attributes.
func (m *Resource) Attributes() []*KeyValue {
	if m._flags&flagResourceAttributesDecoded == 0 {
		// Decode nested message(s).
		for i := range m.attributes {
			// TODO: decide how to handle decoding errors.
			_ = m.attributes[i].decode()
		}
		m._flags |= flagResourceAttributesDecoded
	}
	return m.attributes
}

// SetAttributes sets the value of the attributes.
func (m *Resource) SetAttributes(v []*KeyValue) {
	m.attributes = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.attributes {
		elem._protoMessage.Parent = &m._protoMessage
	}

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *Resource) AttributesRemoveIf(f func(*KeyValue) bool) {
	// Call getter to load the field.
	m.Attributes()

	newLen := 0
	for i := 0; i < len(m.attributes); i++ {
		if f(m.attributes[i]) {
			continue
		}
		if newLen == i {
			// Nothing to move, element is at the right place.
			newLen++
			continue
		}
		m.attributes[newLen] = m.attributes[i]
		newLen++
	}
	if newLen != len(m.attributes) {
		m.attributes = m.attributes[:newLen]
		// Mark this message modified, if not already.
		m._protoMessage.MarkModified()
	}
}

// HasDroppedAttributesCount returns true if the droppedAttributesCount is present.
func (m *Resource) HasDroppedAttributesCount() bool {
	return m._flags&flagResourceDroppedAttributesCountPresent != 0
}

// DroppedAttributesCount returns the value of the droppedAttributesCount.
func (m *Resource) DroppedAttributesCount() uint32 {
	return m.droppedAttributesCount
}

// SetDroppedAttributesCount sets the value of the droppedAttributesCount.
func (m *Resource) SetDroppedAttributesCount(v uint32) {
	m.droppedAttributesCount = v
	m._flags |= flagResourceDroppedAttributesCountPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *Resource) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// If the user makes a mistake and takes a copy of this struct before decoding it
	// the "decoded" flag will incorrectly set on the copy, not the original, but nested
	// messages will be marked as "decoded". Next time we try to access the nested
	// message via getter func this decode() func will be called again and will overwrite
	// nested message values, but the getter func will no longer attempt to decode it
	// because the flag "decoded" flag is incorrectly set on nested message.
	// This will result in incorrect state of nested message returned by getter.
	// To make sure we correctly decode even after this mistake we reset all "decoded"
	// and "presence" flags here.
	m._flags = 0

	// Count all repeated fields. We need one counter per field.
	attributesCount := 0
	err := molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 1 {
				attributesCount++
			}
		},
	)
	if err != nil {
		return err
	}

	// Pre-allocate slices for repeated fields.
	m.attributes = keyValuePool.GetSlice(attributesCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// Set slice indexes to 0 to begin iterating over repeated fields.
	attributesCount = 0
	// Iterate and decode the fields.
	err2 := molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode "attributes".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.attributes[attributesCount]
				attributesCount++
				elem._protoMessage.Parent = &m._protoMessage
				elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
			case 2:
				// Decode "droppedAttributesCount".
				v, err := value.AsUint32()
				if err != nil {
					return false, err
				}
				m.droppedAttributesCount = v
				m._flags |= flagResourceDroppedAttributesCountPresent
			}
			return true, nil
		},
	)
	if err2 != nil {
		return err2
	}
	return nil
}

var preparedResourceAttributes = molecule.PrepareEmbeddedField(1)
var preparedResourceDroppedAttributesCount = molecule.PrepareUint32Field(2)

func (m *Resource) Marshal(ps *molecule.ProtoStream) error {
	if m._protoMessage.IsModified() {
		// Marshal "attributes".
		for _, elem := range m.attributes {
			token := ps.BeginEmbedded()
			if err := elem.Marshal(ps); err != nil {
				return err
			}
			ps.EndEmbeddedPrepared(token, preparedResourceAttributes)
		}
		// Marshal "droppedAttributesCount".
		if m._flags&flagResourceDroppedAttributesCountPresent != 0 {
			ps.Uint32Prepared(preparedResourceDroppedAttributesCount, m.droppedAttributesCount)
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))
	}
	return nil
}

// Pool of Resource structs.
type resourcePoolType struct {
	pool []*Resource
	mux  sync.Mutex
}

var resourcePool = resourcePoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *resourcePoolType) Get() *Resource {
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
	return &Resource{}
}

func (p *resourcePoolType) GetSlice(count int) []*Resource {
	// Create a new slice.
	r := make([]*Resource, count)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return r
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]Resource, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *resourcePoolType) ReleaseSlice(slice []*Resource) {
	for _, elem := range slice {
		// Release nested attributes recursively to their pool.
		keyValuePool.ReleaseSlice(elem.attributes)

		// Zero-initialize the released element.
		*elem = Resource{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *resourcePoolType) Release(elem *Resource) {
	// Release nested attributes recursively to their pool.
	keyValuePool.ReleaseSlice(elem.attributes)

	// Zero-initialize the released element.
	*elem = Resource{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}

// ====================== ScopeLogs message implementation ======================

// A collection of Logs produced by a Scope.
type ScopeLogs struct {
	_protoMessage protomessage.ProtoMessage
	_flags        ScopeLogsFlags

	scope *InstrumentationScope
	// A list of log records.
	logRecords []*LogRecord
	// This schema_url applies to all logs in the "logs" field.
	schemaUrl string
}

func UnmarshalScopeLogs(bytes []byte) (*ScopeLogs, error) {
	m := scopeLogsPool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *ScopeLogs) Free() {
	scopeLogsPool.Release(m)
}

// ScopeLogsFlags is the type of the bit flags.
type ScopeLogsFlags uint8

// Bitmasks that indicate that the particular nested message is decoded.
const flagScopeLogsScopeDecoded ScopeLogsFlags = 0x1
const flagScopeLogsLogRecordsDecoded ScopeLogsFlags = 0x2

// Bitmasks that indicate that the particular field is present.
const flagScopeLogsSchemaUrlPresent ScopeLogsFlags = 0x4

// HasScope returns true if the scope is present.
func (m *ScopeLogs) HasScope() bool {
	return m.scope != nil
}

// Scope returns the value of the scope.
func (m *ScopeLogs) Scope() *InstrumentationScope {
	if m._flags&flagScopeLogsScopeDecoded == 0 {
		// Decode nested message(s).
		scope := m.scope
		if scope != nil {
			// TODO: decide how to handle decoding errors.
			_ = scope.decode()
		}
		m._flags |= flagScopeLogsScopeDecoded
	}
	return m.scope
}

// SetScope sets the value of the scope.
func (m *ScopeLogs) SetScope(v *InstrumentationScope) {
	m.scope = v

	// Make sure the field's Parent points to this message.
	v._protoMessage.Parent = &m._protoMessage

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasLogRecords returns true if the logRecords is present.
func (m *ScopeLogs) HasLogRecords() bool {
	return len(m.logRecords) > 0
}

// LogRecords returns the value of the logRecords.
func (m *ScopeLogs) LogRecords() []*LogRecord {
	if m._flags&flagScopeLogsLogRecordsDecoded == 0 {
		// Decode nested message(s).
		for i := range m.logRecords {
			// TODO: decide how to handle decoding errors.
			_ = m.logRecords[i].decode()
		}
		m._flags |= flagScopeLogsLogRecordsDecoded
	}
	return m.logRecords
}

// SetLogRecords sets the value of the logRecords.
func (m *ScopeLogs) SetLogRecords(v []*LogRecord) {
	m.logRecords = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.logRecords {
		elem._protoMessage.Parent = &m._protoMessage
	}

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *ScopeLogs) LogRecordsRemoveIf(f func(*LogRecord) bool) {
	// Call getter to load the field.
	m.LogRecords()

	newLen := 0
	for i := 0; i < len(m.logRecords); i++ {
		if f(m.logRecords[i]) {
			continue
		}
		if newLen == i {
			// Nothing to move, element is at the right place.
			newLen++
			continue
		}
		m.logRecords[newLen] = m.logRecords[i]
		newLen++
	}
	if newLen != len(m.logRecords) {
		m.logRecords = m.logRecords[:newLen]
		// Mark this message modified, if not already.
		m._protoMessage.MarkModified()
	}
}

// HasSchemaUrl returns true if the schemaUrl is present.
func (m *ScopeLogs) HasSchemaUrl() bool {
	return m._flags&flagScopeLogsSchemaUrlPresent != 0
}

// SchemaUrl returns the value of the schemaUrl.
func (m *ScopeLogs) SchemaUrl() string {
	return m.schemaUrl
}

// SetSchemaUrl sets the value of the schemaUrl.
func (m *ScopeLogs) SetSchemaUrl(v string) {
	m.schemaUrl = v
	m._flags |= flagScopeLogsSchemaUrlPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *ScopeLogs) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// If the user makes a mistake and takes a copy of this struct before decoding it
	// the "decoded" flag will incorrectly set on the copy, not the original, but nested
	// messages will be marked as "decoded". Next time we try to access the nested
	// message via getter func this decode() func will be called again and will overwrite
	// nested message values, but the getter func will no longer attempt to decode it
	// because the flag "decoded" flag is incorrectly set on nested message.
	// This will result in incorrect state of nested message returned by getter.
	// To make sure we correctly decode even after this mistake we reset all "decoded"
	// and "presence" flags here.
	m._flags = 0

	// Count all repeated fields. We need one counter per field.
	logRecordsCount := 0
	err := molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 2 {
				logRecordsCount++
			}
		},
	)
	if err != nil {
		return err
	}

	// Pre-allocate slices for repeated fields.
	m.logRecords = logRecordPool.GetSlice(logRecordsCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// Set slice indexes to 0 to begin iterating over repeated fields.
	logRecordsCount = 0
	// Iterate and decode the fields.
	err2 := molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode "scope".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.scope = instrumentationScopePool.Get()
				m.scope._protoMessage.Parent = &m._protoMessage
				m.scope._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
			case 2:
				// Decode "logRecords".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.logRecords[logRecordsCount]
				logRecordsCount++
				elem._protoMessage.Parent = &m._protoMessage
				elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
			case 3:
				// Decode "schemaUrl".
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.schemaUrl = v
				m._flags |= flagScopeLogsSchemaUrlPresent
			}
			return true, nil
		},
	)
	if err2 != nil {
		return err2
	}
	return nil
}

var preparedScopeLogsScope = molecule.PrepareEmbeddedField(1)
var preparedScopeLogsLogRecords = molecule.PrepareEmbeddedField(2)
var preparedScopeLogsSchemaUrl = molecule.PrepareStringField(3)

func (m *ScopeLogs) Marshal(ps *molecule.ProtoStream) error {
	if m._protoMessage.IsModified() {
		// Marshal "scope".
		elem := m.scope
		if elem != nil {
			token := ps.BeginEmbedded()
			if err := elem.Marshal(ps); err != nil {
				return err
			}
			ps.EndEmbeddedPrepared(token, preparedScopeLogsScope)
		}
		// Marshal "logRecords".
		for _, elem := range m.logRecords {
			token := ps.BeginEmbedded()
			if err := elem.Marshal(ps); err != nil {
				return err
			}
			ps.EndEmbeddedPrepared(token, preparedScopeLogsLogRecords)
		}
		// Marshal "schemaUrl".
		if m._flags&flagScopeLogsSchemaUrlPresent != 0 {
			ps.StringPrepared(preparedScopeLogsSchemaUrl, m.schemaUrl)
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))
	}
	return nil
}

// Pool of ScopeLogs structs.
type scopeLogsPoolType struct {
	pool []*ScopeLogs
	mux  sync.Mutex
}

var scopeLogsPool = scopeLogsPoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *scopeLogsPoolType) Get() *ScopeLogs {
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
	return &ScopeLogs{}
}

func (p *scopeLogsPoolType) GetSlice(count int) []*ScopeLogs {
	// Create a new slice.
	r := make([]*ScopeLogs, count)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return r
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]ScopeLogs, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *scopeLogsPoolType) ReleaseSlice(slice []*ScopeLogs) {
	for _, elem := range slice {
		// Release nested scope recursively to their pool.
		if elem.scope != nil {
			instrumentationScopePool.Release(elem.scope)
		}
		// Release nested logRecords recursively to their pool.
		logRecordPool.ReleaseSlice(elem.logRecords)

		// Zero-initialize the released element.
		*elem = ScopeLogs{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *scopeLogsPoolType) Release(elem *ScopeLogs) {
	// Release nested scope recursively to their pool.
	if elem.scope != nil {
		instrumentationScopePool.Release(elem.scope)
	}
	// Release nested logRecords recursively to their pool.
	logRecordPool.ReleaseSlice(elem.logRecords)

	// Zero-initialize the released element.
	*elem = ScopeLogs{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}

// ====================== InstrumentationScope message implementation ======================

type InstrumentationScope struct {
	_protoMessage protomessage.ProtoMessage
	_flags        InstrumentationScopeFlags

	name                   string
	version                string
	attributes             []*KeyValue
	droppedAttributesCount uint32
}

func UnmarshalInstrumentationScope(bytes []byte) (*InstrumentationScope, error) {
	m := instrumentationScopePool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *InstrumentationScope) Free() {
	instrumentationScopePool.Release(m)
}

// InstrumentationScopeFlags is the type of the bit flags.
type InstrumentationScopeFlags uint8

// Bitmasks that indicate that the particular nested message is decoded.
const flagInstrumentationScopeAttributesDecoded InstrumentationScopeFlags = 0x1

// Bitmasks that indicate that the particular field is present.
const flagInstrumentationScopeNamePresent InstrumentationScopeFlags = 0x2
const flagInstrumentationScopeVersionPresent InstrumentationScopeFlags = 0x4
const flagInstrumentationScopeDroppedAttributesCountPresent InstrumentationScopeFlags = 0x8

// HasName returns true if the name is present.
func (m *InstrumentationScope) HasName() bool {
	return m._flags&flagInstrumentationScopeNamePresent != 0
}

// Name returns the value of the name.
func (m *InstrumentationScope) Name() string {
	return m.name
}

// SetName sets the value of the name.
func (m *InstrumentationScope) SetName(v string) {
	m.name = v
	m._flags |= flagInstrumentationScopeNamePresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasVersion returns true if the version is present.
func (m *InstrumentationScope) HasVersion() bool {
	return m._flags&flagInstrumentationScopeVersionPresent != 0
}

// Version returns the value of the version.
func (m *InstrumentationScope) Version() string {
	return m.version
}

// SetVersion sets the value of the version.
func (m *InstrumentationScope) SetVersion(v string) {
	m.version = v
	m._flags |= flagInstrumentationScopeVersionPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasAttributes returns true if the attributes is present.
func (m *InstrumentationScope) HasAttributes() bool {
	return len(m.attributes) > 0
}

// Attributes returns the value of the attributes.
func (m *InstrumentationScope) Attributes() []*KeyValue {
	if m._flags&flagInstrumentationScopeAttributesDecoded == 0 {
		// Decode nested message(s).
		for i := range m.attributes {
			// TODO: decide how to handle decoding errors.
			_ = m.attributes[i].decode()
		}
		m._flags |= flagInstrumentationScopeAttributesDecoded
	}
	return m.attributes
}

// SetAttributes sets the value of the attributes.
func (m *InstrumentationScope) SetAttributes(v []*KeyValue) {
	m.attributes = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.attributes {
		elem._protoMessage.Parent = &m._protoMessage
	}

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *InstrumentationScope) AttributesRemoveIf(f func(*KeyValue) bool) {
	// Call getter to load the field.
	m.Attributes()

	newLen := 0
	for i := 0; i < len(m.attributes); i++ {
		if f(m.attributes[i]) {
			continue
		}
		if newLen == i {
			// Nothing to move, element is at the right place.
			newLen++
			continue
		}
		m.attributes[newLen] = m.attributes[i]
		newLen++
	}
	if newLen != len(m.attributes) {
		m.attributes = m.attributes[:newLen]
		// Mark this message modified, if not already.
		m._protoMessage.MarkModified()
	}
}

// HasDroppedAttributesCount returns true if the droppedAttributesCount is present.
func (m *InstrumentationScope) HasDroppedAttributesCount() bool {
	return m._flags&flagInstrumentationScopeDroppedAttributesCountPresent != 0
}

// DroppedAttributesCount returns the value of the droppedAttributesCount.
func (m *InstrumentationScope) DroppedAttributesCount() uint32 {
	return m.droppedAttributesCount
}

// SetDroppedAttributesCount sets the value of the droppedAttributesCount.
func (m *InstrumentationScope) SetDroppedAttributesCount(v uint32) {
	m.droppedAttributesCount = v
	m._flags |= flagInstrumentationScopeDroppedAttributesCountPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *InstrumentationScope) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// If the user makes a mistake and takes a copy of this struct before decoding it
	// the "decoded" flag will incorrectly set on the copy, not the original, but nested
	// messages will be marked as "decoded". Next time we try to access the nested
	// message via getter func this decode() func will be called again and will overwrite
	// nested message values, but the getter func will no longer attempt to decode it
	// because the flag "decoded" flag is incorrectly set on nested message.
	// This will result in incorrect state of nested message returned by getter.
	// To make sure we correctly decode even after this mistake we reset all "decoded"
	// and "presence" flags here.
	m._flags = 0

	// Count all repeated fields. We need one counter per field.
	attributesCount := 0
	err := molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 3 {
				attributesCount++
			}
		},
	)
	if err != nil {
		return err
	}

	// Pre-allocate slices for repeated fields.
	m.attributes = keyValuePool.GetSlice(attributesCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// Set slice indexes to 0 to begin iterating over repeated fields.
	attributesCount = 0
	// Iterate and decode the fields.
	err2 := molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode "name".
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.name = v
				m._flags |= flagInstrumentationScopeNamePresent
			case 2:
				// Decode "version".
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.version = v
				m._flags |= flagInstrumentationScopeVersionPresent
			case 3:
				// Decode "attributes".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.attributes[attributesCount]
				attributesCount++
				elem._protoMessage.Parent = &m._protoMessage
				elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
			case 4:
				// Decode "droppedAttributesCount".
				v, err := value.AsUint32()
				if err != nil {
					return false, err
				}
				m.droppedAttributesCount = v
				m._flags |= flagInstrumentationScopeDroppedAttributesCountPresent
			}
			return true, nil
		},
	)
	if err2 != nil {
		return err2
	}
	return nil
}

var preparedInstrumentationScopeName = molecule.PrepareStringField(1)
var preparedInstrumentationScopeVersion = molecule.PrepareStringField(2)
var preparedInstrumentationScopeAttributes = molecule.PrepareEmbeddedField(3)
var preparedInstrumentationScopeDroppedAttributesCount = molecule.PrepareUint32Field(4)

func (m *InstrumentationScope) Marshal(ps *molecule.ProtoStream) error {
	if m._protoMessage.IsModified() {
		// Marshal "name".
		if m._flags&flagInstrumentationScopeNamePresent != 0 {
			ps.StringPrepared(preparedInstrumentationScopeName, m.name)
		}
		// Marshal "version".
		if m._flags&flagInstrumentationScopeVersionPresent != 0 {
			ps.StringPrepared(preparedInstrumentationScopeVersion, m.version)
		}
		// Marshal "attributes".
		for _, elem := range m.attributes {
			token := ps.BeginEmbedded()
			if err := elem.Marshal(ps); err != nil {
				return err
			}
			ps.EndEmbeddedPrepared(token, preparedInstrumentationScopeAttributes)
		}
		// Marshal "droppedAttributesCount".
		if m._flags&flagInstrumentationScopeDroppedAttributesCountPresent != 0 {
			ps.Uint32Prepared(preparedInstrumentationScopeDroppedAttributesCount, m.droppedAttributesCount)
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))
	}
	return nil
}

// Pool of InstrumentationScope structs.
type instrumentationScopePoolType struct {
	pool []*InstrumentationScope
	mux  sync.Mutex
}

var instrumentationScopePool = instrumentationScopePoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *instrumentationScopePoolType) Get() *InstrumentationScope {
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
	return &InstrumentationScope{}
}

func (p *instrumentationScopePoolType) GetSlice(count int) []*InstrumentationScope {
	// Create a new slice.
	r := make([]*InstrumentationScope, count)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return r
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]InstrumentationScope, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *instrumentationScopePoolType) ReleaseSlice(slice []*InstrumentationScope) {
	for _, elem := range slice {
		// Release nested attributes recursively to their pool.
		keyValuePool.ReleaseSlice(elem.attributes)

		// Zero-initialize the released element.
		*elem = InstrumentationScope{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *instrumentationScopePoolType) Release(elem *InstrumentationScope) {
	// Release nested attributes recursively to their pool.
	keyValuePool.ReleaseSlice(elem.attributes)

	// Zero-initialize the released element.
	*elem = InstrumentationScope{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}

// ====================== LogRecord message implementation ======================

type LogRecord struct {
	_protoMessage protomessage.ProtoMessage
	_flags        LogRecordFlags

	timeUnixNano           uint64
	observedTimeUnixNano   uint64
	severityNumber         SeverityNumber
	severityText           string
	attributes             []*KeyValue
	droppedAttributesCount uint32
	flags                  uint32
	traceId                []byte
	spanId                 []byte
}

func UnmarshalLogRecord(bytes []byte) (*LogRecord, error) {
	m := logRecordPool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *LogRecord) Free() {
	logRecordPool.Release(m)
}

// LogRecordFlags is the type of the bit flags.
type LogRecordFlags uint16

// Bitmasks that indicate that the particular nested message is decoded.
const flagLogRecordAttributesDecoded LogRecordFlags = 0x1

// Bitmasks that indicate that the particular field is present.
const flagLogRecordTimeUnixNanoPresent LogRecordFlags = 0x2
const flagLogRecordObservedTimeUnixNanoPresent LogRecordFlags = 0x4
const flagLogRecordSeverityNumberPresent LogRecordFlags = 0x8
const flagLogRecordSeverityTextPresent LogRecordFlags = 0x10
const flagLogRecordDroppedAttributesCountPresent LogRecordFlags = 0x20
const flagLogRecordFlagsPresent LogRecordFlags = 0x40
const flagLogRecordTraceIdPresent LogRecordFlags = 0x80
const flagLogRecordSpanIdPresent LogRecordFlags = 0x100

// HasTimeUnixNano returns true if the timeUnixNano is present.
func (m *LogRecord) HasTimeUnixNano() bool {
	return m._flags&flagLogRecordTimeUnixNanoPresent != 0
}

// TimeUnixNano returns the value of the timeUnixNano.
func (m *LogRecord) TimeUnixNano() uint64 {
	return m.timeUnixNano
}

// SetTimeUnixNano sets the value of the timeUnixNano.
func (m *LogRecord) SetTimeUnixNano(v uint64) {
	m.timeUnixNano = v
	m._flags |= flagLogRecordTimeUnixNanoPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasObservedTimeUnixNano returns true if the observedTimeUnixNano is present.
func (m *LogRecord) HasObservedTimeUnixNano() bool {
	return m._flags&flagLogRecordObservedTimeUnixNanoPresent != 0
}

// ObservedTimeUnixNano returns the value of the observedTimeUnixNano.
func (m *LogRecord) ObservedTimeUnixNano() uint64 {
	return m.observedTimeUnixNano
}

// SetObservedTimeUnixNano sets the value of the observedTimeUnixNano.
func (m *LogRecord) SetObservedTimeUnixNano(v uint64) {
	m.observedTimeUnixNano = v
	m._flags |= flagLogRecordObservedTimeUnixNanoPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasSeverityNumber returns true if the severityNumber is present.
func (m *LogRecord) HasSeverityNumber() bool {
	return m._flags&flagLogRecordSeverityNumberPresent != 0
}

// SeverityNumber returns the value of the severityNumber.
func (m *LogRecord) SeverityNumber() SeverityNumber {
	return m.severityNumber
}

// SetSeverityNumber sets the value of the severityNumber.
func (m *LogRecord) SetSeverityNumber(v SeverityNumber) {
	m.severityNumber = v
	m._flags |= flagLogRecordSeverityNumberPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasSeverityText returns true if the severityText is present.
func (m *LogRecord) HasSeverityText() bool {
	return m._flags&flagLogRecordSeverityTextPresent != 0
}

// SeverityText returns the value of the severityText.
func (m *LogRecord) SeverityText() string {
	return m.severityText
}

// SetSeverityText sets the value of the severityText.
func (m *LogRecord) SetSeverityText(v string) {
	m.severityText = v
	m._flags |= flagLogRecordSeverityTextPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasAttributes returns true if the attributes is present.
func (m *LogRecord) HasAttributes() bool {
	return len(m.attributes) > 0
}

// Attributes returns the value of the attributes.
func (m *LogRecord) Attributes() []*KeyValue {
	if m._flags&flagLogRecordAttributesDecoded == 0 {
		// Decode nested message(s).
		for i := range m.attributes {
			// TODO: decide how to handle decoding errors.
			_ = m.attributes[i].decode()
		}
		m._flags |= flagLogRecordAttributesDecoded
	}
	return m.attributes
}

// SetAttributes sets the value of the attributes.
func (m *LogRecord) SetAttributes(v []*KeyValue) {
	m.attributes = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.attributes {
		elem._protoMessage.Parent = &m._protoMessage
	}

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *LogRecord) AttributesRemoveIf(f func(*KeyValue) bool) {
	// Call getter to load the field.
	m.Attributes()

	newLen := 0
	for i := 0; i < len(m.attributes); i++ {
		if f(m.attributes[i]) {
			continue
		}
		if newLen == i {
			// Nothing to move, element is at the right place.
			newLen++
			continue
		}
		m.attributes[newLen] = m.attributes[i]
		newLen++
	}
	if newLen != len(m.attributes) {
		m.attributes = m.attributes[:newLen]
		// Mark this message modified, if not already.
		m._protoMessage.MarkModified()
	}
}

// HasDroppedAttributesCount returns true if the droppedAttributesCount is present.
func (m *LogRecord) HasDroppedAttributesCount() bool {
	return m._flags&flagLogRecordDroppedAttributesCountPresent != 0
}

// DroppedAttributesCount returns the value of the droppedAttributesCount.
func (m *LogRecord) DroppedAttributesCount() uint32 {
	return m.droppedAttributesCount
}

// SetDroppedAttributesCount sets the value of the droppedAttributesCount.
func (m *LogRecord) SetDroppedAttributesCount(v uint32) {
	m.droppedAttributesCount = v
	m._flags |= flagLogRecordDroppedAttributesCountPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasFlags returns true if the flags is present.
func (m *LogRecord) HasFlags() bool {
	return m._flags&flagLogRecordFlagsPresent != 0
}

// Flags returns the value of the flags.
func (m *LogRecord) Flags() uint32 {
	return m.flags
}

// SetFlags sets the value of the flags.
func (m *LogRecord) SetFlags(v uint32) {
	m.flags = v
	m._flags |= flagLogRecordFlagsPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasTraceId returns true if the traceId is present.
func (m *LogRecord) HasTraceId() bool {
	return m._flags&flagLogRecordTraceIdPresent != 0
}

// TraceId returns the value of the traceId.
func (m *LogRecord) TraceId() []byte {
	return m.traceId
}

// SetTraceId sets the value of the traceId.
func (m *LogRecord) SetTraceId(v []byte) {
	m.traceId = v
	m._flags |= flagLogRecordTraceIdPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasSpanId returns true if the spanId is present.
func (m *LogRecord) HasSpanId() bool {
	return m._flags&flagLogRecordSpanIdPresent != 0
}

// SpanId returns the value of the spanId.
func (m *LogRecord) SpanId() []byte {
	return m.spanId
}

// SetSpanId sets the value of the spanId.
func (m *LogRecord) SetSpanId(v []byte) {
	m.spanId = v
	m._flags |= flagLogRecordSpanIdPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *LogRecord) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// If the user makes a mistake and takes a copy of this struct before decoding it
	// the "decoded" flag will incorrectly set on the copy, not the original, but nested
	// messages will be marked as "decoded". Next time we try to access the nested
	// message via getter func this decode() func will be called again and will overwrite
	// nested message values, but the getter func will no longer attempt to decode it
	// because the flag "decoded" flag is incorrectly set on nested message.
	// This will result in incorrect state of nested message returned by getter.
	// To make sure we correctly decode even after this mistake we reset all "decoded"
	// and "presence" flags here.
	m._flags = 0

	// Count all repeated fields. We need one counter per field.
	attributesCount := 0
	err := molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 6 {
				attributesCount++
			}
		},
	)
	if err != nil {
		return err
	}

	// Pre-allocate slices for repeated fields.
	m.attributes = keyValuePool.GetSlice(attributesCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// Set slice indexes to 0 to begin iterating over repeated fields.
	attributesCount = 0
	// Iterate and decode the fields.
	err2 := molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode "timeUnixNano".
				v, err := value.AsFixed64()
				if err != nil {
					return false, err
				}
				m.timeUnixNano = v
				m._flags |= flagLogRecordTimeUnixNanoPresent
			case 11:
				// Decode "observedTimeUnixNano".
				v, err := value.AsFixed64()
				if err != nil {
					return false, err
				}
				m.observedTimeUnixNano = v
				m._flags |= flagLogRecordObservedTimeUnixNanoPresent
			case 2:
				// Decode "severityNumber".
				v, err := value.AsUint32()
				if err != nil {
					return false, err
				}
				m.severityNumber = SeverityNumber(v)
				m._flags |= flagLogRecordSeverityNumberPresent
			case 3:
				// Decode "severityText".
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.severityText = v
				m._flags |= flagLogRecordSeverityTextPresent
			case 6:
				// Decode "attributes".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.attributes[attributesCount]
				attributesCount++
				elem._protoMessage.Parent = &m._protoMessage
				elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
			case 7:
				// Decode "droppedAttributesCount".
				v, err := value.AsUint32()
				if err != nil {
					return false, err
				}
				m.droppedAttributesCount = v
				m._flags |= flagLogRecordDroppedAttributesCountPresent
			case 8:
				// Decode "flags".
				v, err := value.AsFixed32()
				if err != nil {
					return false, err
				}
				m.flags = v
				m._flags |= flagLogRecordFlagsPresent
			case 9:
				// Decode "traceId".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.traceId = v
				m._flags |= flagLogRecordTraceIdPresent
			case 10:
				// Decode "spanId".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.spanId = v
				m._flags |= flagLogRecordSpanIdPresent
			}
			return true, nil
		},
	)
	if err2 != nil {
		return err2
	}
	return nil
}

var preparedLogRecordTimeUnixNano = molecule.PrepareFixed64Field(1)
var preparedLogRecordObservedTimeUnixNano = molecule.PrepareFixed64Field(11)
var preparedLogRecordSeverityNumber = molecule.PrepareUint32Field(2)
var preparedLogRecordSeverityText = molecule.PrepareStringField(3)
var preparedLogRecordAttributes = molecule.PrepareEmbeddedField(6)
var preparedLogRecordDroppedAttributesCount = molecule.PrepareUint32Field(7)
var preparedLogRecordFlags = molecule.PrepareFixed32Field(8)
var preparedLogRecordTraceId = molecule.PrepareBytesField(9)
var preparedLogRecordSpanId = molecule.PrepareBytesField(10)

func (m *LogRecord) Marshal(ps *molecule.ProtoStream) error {
	if m._protoMessage.IsModified() {
		// Marshal "timeUnixNano".
		if m._flags&flagLogRecordTimeUnixNanoPresent != 0 {
			ps.Fixed64Prepared(preparedLogRecordTimeUnixNano, m.timeUnixNano)
		}
		// Marshal "severityNumber".
		if m._flags&flagLogRecordSeverityNumberPresent != 0 {
			ps.Uint32Prepared(preparedLogRecordSeverityNumber, uint32(m.severityNumber))
		}
		// Marshal "severityText".
		if m._flags&flagLogRecordSeverityTextPresent != 0 {
			ps.StringPrepared(preparedLogRecordSeverityText, m.severityText)
		}
		// Marshal "attributes".
		for _, elem := range m.attributes {
			token := ps.BeginEmbedded()
			if err := elem.Marshal(ps); err != nil {
				return err
			}
			ps.EndEmbeddedPrepared(token, preparedLogRecordAttributes)
		}
		// Marshal "droppedAttributesCount".
		if m._flags&flagLogRecordDroppedAttributesCountPresent != 0 {
			ps.Uint32Prepared(preparedLogRecordDroppedAttributesCount, m.droppedAttributesCount)
		}
		// Marshal "flags".
		if m._flags&flagLogRecordFlagsPresent != 0 {
			ps.Fixed32Prepared(preparedLogRecordFlags, m.flags)
		}
		// Marshal "traceId".
		if m._flags&flagLogRecordTraceIdPresent != 0 {
			ps.BytesPrepared(preparedLogRecordTraceId, m.traceId)
		}
		// Marshal "spanId".
		if m._flags&flagLogRecordSpanIdPresent != 0 {
			ps.BytesPrepared(preparedLogRecordSpanId, m.spanId)
		}
		// Marshal "observedTimeUnixNano".
		if m._flags&flagLogRecordObservedTimeUnixNanoPresent != 0 {
			ps.Fixed64Prepared(preparedLogRecordObservedTimeUnixNano, m.observedTimeUnixNano)
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))
	}
	return nil
}

// Pool of LogRecord structs.
type logRecordPoolType struct {
	pool []*LogRecord
	mux  sync.Mutex
}

var logRecordPool = logRecordPoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *logRecordPoolType) Get() *LogRecord {
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
	return &LogRecord{}
}

func (p *logRecordPoolType) GetSlice(count int) []*LogRecord {
	// Create a new slice.
	r := make([]*LogRecord, count)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return r
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]LogRecord, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *logRecordPoolType) ReleaseSlice(slice []*LogRecord) {
	for _, elem := range slice {
		// Release nested attributes recursively to their pool.
		keyValuePool.ReleaseSlice(elem.attributes)

		// Zero-initialize the released element.
		*elem = LogRecord{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *logRecordPoolType) Release(elem *LogRecord) {
	// Release nested attributes recursively to their pool.
	keyValuePool.ReleaseSlice(elem.attributes)

	// Zero-initialize the released element.
	*elem = LogRecord{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}

// ====================== KeyValue message implementation ======================

type KeyValue struct {
	_protoMessage protomessage.ProtoMessage
	_flags        KeyValueFlags

	key   string
	value *AnyValue
}

func UnmarshalKeyValue(bytes []byte) (*KeyValue, error) {
	m := keyValuePool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *KeyValue) Free() {
	keyValuePool.Release(m)
}

// KeyValueFlags is the type of the bit flags.
type KeyValueFlags uint8

// Bitmasks that indicate that the particular nested message is decoded.
const flagKeyValueValueDecoded KeyValueFlags = 0x1

// Bitmasks that indicate that the particular field is present.
const flagKeyValueKeyPresent KeyValueFlags = 0x2

// HasKey returns true if the key is present.
func (m *KeyValue) HasKey() bool {
	return m._flags&flagKeyValueKeyPresent != 0
}

// Key returns the value of the key.
func (m *KeyValue) Key() string {
	return m.key
}

// SetKey sets the value of the key.
func (m *KeyValue) SetKey(v string) {
	m.key = v
	m._flags |= flagKeyValueKeyPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasValue returns true if the value is present.
func (m *KeyValue) HasValue() bool {
	return m.value != nil
}

// Value returns the value of the value.
func (m *KeyValue) Value() *AnyValue {
	if m._flags&flagKeyValueValueDecoded == 0 {
		// Decode nested message(s).
		value := m.value
		if value != nil {
			// TODO: decide how to handle decoding errors.
			_ = value.decode()
		}
		m._flags |= flagKeyValueValueDecoded
	}
	return m.value
}

// SetValue sets the value of the value.
func (m *KeyValue) SetValue(v *AnyValue) {
	m.value = v

	// Make sure the field's Parent points to this message.
	v._protoMessage.Parent = &m._protoMessage

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *KeyValue) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// If the user makes a mistake and takes a copy of this struct before decoding it
	// the "decoded" flag will incorrectly set on the copy, not the original, but nested
	// messages will be marked as "decoded". Next time we try to access the nested
	// message via getter func this decode() func will be called again and will overwrite
	// nested message values, but the getter func will no longer attempt to decode it
	// because the flag "decoded" flag is incorrectly set on nested message.
	// This will result in incorrect state of nested message returned by getter.
	// To make sure we correctly decode even after this mistake we reset all "decoded"
	// and "presence" flags here.
	m._flags = 0

	// Iterate and decode the fields.
	err2 := molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode "key".
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.key = v
				m._flags |= flagKeyValueKeyPresent
			case 2:
				// Decode "value".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.value = anyValuePool.Get()
				m.value._protoMessage.Parent = &m._protoMessage
				m.value._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
			}
			return true, nil
		},
	)
	if err2 != nil {
		return err2
	}
	return nil
}

var preparedKeyValueKey = molecule.PrepareStringField(1)
var preparedKeyValueValue = molecule.PrepareEmbeddedField(2)

func (m *KeyValue) Marshal(ps *molecule.ProtoStream) error {
	if m._protoMessage.IsModified() {
		// Marshal "key".
		if m._flags&flagKeyValueKeyPresent != 0 {
			ps.StringPrepared(preparedKeyValueKey, m.key)
		}
		// Marshal "value".
		elem := m.value
		if elem != nil {
			token := ps.BeginEmbedded()
			if err := elem.Marshal(ps); err != nil {
				return err
			}
			ps.EndEmbeddedPrepared(token, preparedKeyValueValue)
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))
	}
	return nil
}

// Pool of KeyValue structs.
type keyValuePoolType struct {
	pool []*KeyValue
	mux  sync.Mutex
}

var keyValuePool = keyValuePoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *keyValuePoolType) Get() *KeyValue {
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
	return &KeyValue{}
}

func (p *keyValuePoolType) GetSlice(count int) []*KeyValue {
	// Create a new slice.
	r := make([]*KeyValue, count)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return r
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]KeyValue, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *keyValuePoolType) ReleaseSlice(slice []*KeyValue) {
	for _, elem := range slice {
		// Release nested value recursively to their pool.
		if elem.value != nil {
			anyValuePool.Release(elem.value)
		}

		// Zero-initialize the released element.
		*elem = KeyValue{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *keyValuePoolType) Release(elem *KeyValue) {
	// Release nested value recursively to their pool.
	if elem.value != nil {
		anyValuePool.Release(elem.value)
	}

	// Zero-initialize the released element.
	*elem = KeyValue{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}

// ====================== AnyValue message implementation ======================

type AnyValue struct {
	_protoMessage protomessage.ProtoMessage
	_flags        AnyValueFlags

	value oneof.OneOf
}

func UnmarshalAnyValue(bytes []byte) (*AnyValue, error) {
	m := anyValuePool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *AnyValue) Free() {
	anyValuePool.Release(m)
}

// AnyValueValue defines the possible types for oneof field "value".
type AnyValueValue int

const (
	// AnyValueValueNone indicates that none of the oneof choices is set.
	AnyValueValueNone AnyValueValue = 0
	// AnyValueStringValue indicates that oneof field "stringValue" is set.
	AnyValueStringValue AnyValueValue = 1
	// AnyValueBoolValue indicates that oneof field "boolValue" is set.
	AnyValueBoolValue AnyValueValue = 2
	// AnyValueIntValue indicates that oneof field "intValue" is set.
	AnyValueIntValue AnyValueValue = 3
	// AnyValueDoubleValue indicates that oneof field "doubleValue" is set.
	AnyValueDoubleValue AnyValueValue = 4
	// AnyValueArrayValue indicates that oneof field "arrayValue" is set.
	AnyValueArrayValue AnyValueValue = 5
	// AnyValueKvlistValue indicates that oneof field "kvlistValue" is set.
	AnyValueKvlistValue AnyValueValue = 6
	// AnyValueBytesValue indicates that oneof field "bytesValue" is set.
	AnyValueBytesValue AnyValueValue = 7
)

// ValueType returns the type of the current stored oneof "value".
// To set the type use one of the setters.
func (m *AnyValue) ValueType() AnyValueValue {
	return AnyValueValue(m.value.FieldIndex())
}

// ValueUnset unsets the oneof field "value", so that it contains none of the choices.
func (m *AnyValue) ValueUnset() {
	m.value = oneof.NewOneOfNone()
}

// AnyValueFlags is the type of the bit flags.
type AnyValueFlags uint8

// Bitmasks that indicate that the particular nested message is decoded.
const flagAnyValueArrayValueDecoded AnyValueFlags = 0x1
const flagAnyValueKvlistValueDecoded AnyValueFlags = 0x2

// StringValue returns the value of the stringValue.
// If the field "value" is not set to "stringValue" then the returned value is undefined.
func (m *AnyValue) StringValue() string {
	return m.value.StringVal()
}

// SetStringValue sets the value of the stringValue.
// The oneof field "value" will be set to "stringValue".
func (m *AnyValue) SetStringValue(v string) {
	m.value = oneof.NewOneOfString(v, int(AnyValueStringValue))

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// BoolValue returns the value of the boolValue.
// If the field "value" is not set to "boolValue" then the returned value is undefined.
func (m *AnyValue) BoolValue() bool {
	return m.value.BoolVal()
}

// SetBoolValue sets the value of the boolValue.
// The oneof field "value" will be set to "boolValue".
func (m *AnyValue) SetBoolValue(v bool) {
	m.value = oneof.NewOneOfBool(v, int(AnyValueBoolValue))

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// IntValue returns the value of the intValue.
// If the field "value" is not set to "intValue" then the returned value is undefined.
func (m *AnyValue) IntValue() int64 {
	return m.value.Int64Val()
}

// SetIntValue sets the value of the intValue.
// The oneof field "value" will be set to "intValue".
func (m *AnyValue) SetIntValue(v int64) {
	m.value = oneof.NewOneOfInt64(v, int(AnyValueIntValue))

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// DoubleValue returns the value of the doubleValue.
// If the field "value" is not set to "doubleValue" then the returned value is undefined.
func (m *AnyValue) DoubleValue() float64 {
	return m.value.DoubleVal()
}

// SetDoubleValue sets the value of the doubleValue.
// The oneof field "value" will be set to "doubleValue".
func (m *AnyValue) SetDoubleValue(v float64) {
	m.value = oneof.NewOneOfDouble(v, int(AnyValueDoubleValue))

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// ArrayValue returns the value of the arrayValue.
// If the field "value" is not set to "arrayValue" then the returned value is undefined.
func (m *AnyValue) ArrayValue() *ArrayValue {
	if m._flags&flagAnyValueArrayValueDecoded == 0 {
		// Decode nested message(s).
		if m.value.FieldIndex() == int(AnyValueArrayValue) {
			arrayValue := (*ArrayValue)(m.value.PtrVal())
			if arrayValue != nil {
				// TODO: decide how to handle decoding errors.
				_ = arrayValue.decode()
			}
		}
		m._flags |= flagAnyValueArrayValueDecoded
	}
	return (*ArrayValue)(m.value.PtrVal())
}

// SetArrayValue sets the value of the arrayValue.
// The oneof field "value" will be set to "arrayValue".
func (m *AnyValue) SetArrayValue(v *ArrayValue) {
	m.value = oneof.NewOneOfPtr(unsafe.Pointer(v), int(AnyValueArrayValue))

	// Make sure the field's Parent points to this message.
	v._protoMessage.Parent = &m._protoMessage

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// KvlistValue returns the value of the kvlistValue.
// If the field "value" is not set to "kvlistValue" then the returned value is undefined.
func (m *AnyValue) KvlistValue() *KeyValueList {
	if m._flags&flagAnyValueKvlistValueDecoded == 0 {
		// Decode nested message(s).
		if m.value.FieldIndex() == int(AnyValueKvlistValue) {
			kvlistValue := (*KeyValueList)(m.value.PtrVal())
			if kvlistValue != nil {
				// TODO: decide how to handle decoding errors.
				_ = kvlistValue.decode()
			}
		}
		m._flags |= flagAnyValueKvlistValueDecoded
	}
	return (*KeyValueList)(m.value.PtrVal())
}

// SetKvlistValue sets the value of the kvlistValue.
// The oneof field "value" will be set to "kvlistValue".
func (m *AnyValue) SetKvlistValue(v *KeyValueList) {
	m.value = oneof.NewOneOfPtr(unsafe.Pointer(v), int(AnyValueKvlistValue))

	// Make sure the field's Parent points to this message.
	v._protoMessage.Parent = &m._protoMessage

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// BytesValue returns the value of the bytesValue.
// If the field "value" is not set to "bytesValue" then the returned value is undefined.
func (m *AnyValue) BytesValue() []byte {
	return m.value.BytesVal()
}

// SetBytesValue sets the value of the bytesValue.
// The oneof field "value" will be set to "bytesValue".
func (m *AnyValue) SetBytesValue(v []byte) {
	m.value = oneof.NewOneOfBytes(v, int(AnyValueBytesValue))

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *AnyValue) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// If the user makes a mistake and takes a copy of this struct before decoding it
	// the "decoded" flag will incorrectly set on the copy, not the original, but nested
	// messages will be marked as "decoded". Next time we try to access the nested
	// message via getter func this decode() func will be called again and will overwrite
	// nested message values, but the getter func will no longer attempt to decode it
	// because the flag "decoded" flag is incorrectly set on nested message.
	// This will result in incorrect state of nested message returned by getter.
	// To make sure we correctly decode even after this mistake we reset all "decoded"
	// and "presence" flags here.
	m._flags = 0

	// Iterate and decode the fields.
	err2 := molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode "stringValue".
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.value = oneof.NewOneOfString(v, int(AnyValueStringValue))
			case 2:
				// Decode "boolValue".
				v, err := value.AsBool()
				if err != nil {
					return false, err
				}
				m.value = oneof.NewOneOfBool(v, int(AnyValueBoolValue))
			case 3:
				// Decode "intValue".
				v, err := value.AsInt64()
				if err != nil {
					return false, err
				}
				m.value = oneof.NewOneOfInt64(v, int(AnyValueIntValue))
			case 4:
				// Decode "doubleValue".
				v, err := value.AsDouble()
				if err != nil {
					return false, err
				}
				m.value = oneof.NewOneOfDouble(v, int(AnyValueDoubleValue))
			case 5:
				// Decode "arrayValue".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				elem := arrayValuePool.Get()
				elem._protoMessage.Parent = &m._protoMessage
				elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
				m.value = oneof.NewOneOfPtr(unsafe.Pointer(elem), int(AnyValueArrayValue))
			case 6:
				// Decode "kvlistValue".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				elem := keyValueListPool.Get()
				elem._protoMessage.Parent = &m._protoMessage
				elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
				m.value = oneof.NewOneOfPtr(unsafe.Pointer(elem), int(AnyValueKvlistValue))
			case 7:
				// Decode "bytesValue".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.value = oneof.NewOneOfBytes(v, int(AnyValueBytesValue))
			}
			return true, nil
		},
	)
	if err2 != nil {
		return err2
	}
	return nil
}

var preparedAnyValueStringValue = molecule.PrepareStringField(1)
var preparedAnyValueBoolValue = molecule.PrepareBoolField(2)
var preparedAnyValueIntValue = molecule.PrepareInt64Field(3)
var preparedAnyValueDoubleValue = molecule.PrepareDoubleField(4)
var preparedAnyValueArrayValue = molecule.PrepareEmbeddedField(5)
var preparedAnyValueKvlistValue = molecule.PrepareEmbeddedField(6)
var preparedAnyValueBytesValue = molecule.PrepareBytesField(7)

func (m *AnyValue) Marshal(ps *molecule.ProtoStream) error {
	if m._protoMessage.IsModified() {
		// Marshal "value".
		switch AnyValueValue(m.value.FieldIndex()) {
		case AnyValueValueNone:
			// Nothing to do, oneof is unset.
		case AnyValueStringValue:
			// Marshal "stringValue".
			ps.StringPrepared(preparedAnyValueStringValue, m.value.StringVal())
		case AnyValueBoolValue:
			// Marshal "boolValue".
			ps.BoolPrepared(preparedAnyValueBoolValue, m.value.BoolVal())
		case AnyValueIntValue:
			// Marshal "intValue".
			ps.Int64Prepared(preparedAnyValueIntValue, m.value.Int64Val())
		case AnyValueDoubleValue:
			// Marshal "doubleValue".
			ps.DoublePrepared(preparedAnyValueDoubleValue, m.value.DoubleVal())
		case AnyValueArrayValue:
			// Marshal "arrayValue".
			elem := (*ArrayValue)(m.value.PtrVal())
			if elem != nil {
				token := ps.BeginEmbedded()
				if err := elem.Marshal(ps); err != nil {
					return err
				}
				ps.EndEmbeddedPrepared(token, preparedAnyValueArrayValue)
			}
		case AnyValueKvlistValue:
			// Marshal "kvlistValue".
			elem := (*KeyValueList)(m.value.PtrVal())
			if elem != nil {
				token := ps.BeginEmbedded()
				if err := elem.Marshal(ps); err != nil {
					return err
				}
				ps.EndEmbeddedPrepared(token, preparedAnyValueKvlistValue)
			}
		case AnyValueBytesValue:
			// Marshal "bytesValue".
			ps.BytesPrepared(preparedAnyValueBytesValue, m.value.BytesVal())
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))
	}
	return nil
}

// Pool of AnyValue structs.
type anyValuePoolType struct {
	pool []*AnyValue
	mux  sync.Mutex
}

var anyValuePool = anyValuePoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *anyValuePoolType) Get() *AnyValue {
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
	return &AnyValue{}
}

func (p *anyValuePoolType) GetSlice(count int) []*AnyValue {
	// Create a new slice.
	r := make([]*AnyValue, count)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return r
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]AnyValue, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *anyValuePoolType) ReleaseSlice(slice []*AnyValue) {
	for _, elem := range slice {
		switch AnyValueValue(elem.value.FieldIndex()) {
		case AnyValueArrayValue:
			ptr := (*ArrayValue)(elem.value.PtrVal())
			if ptr != nil {
				arrayValuePool.Release(ptr)
			}
		case AnyValueKvlistValue:
			ptr := (*KeyValueList)(elem.value.PtrVal())
			if ptr != nil {
				keyValueListPool.Release(ptr)
			}
		}

		// Zero-initialize the released element.
		*elem = AnyValue{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *anyValuePoolType) Release(elem *AnyValue) {
	switch AnyValueValue(elem.value.FieldIndex()) {
	case AnyValueArrayValue:
		ptr := (*ArrayValue)(elem.value.PtrVal())
		if ptr != nil {
			arrayValuePool.Release(ptr)
		}
	case AnyValueKvlistValue:
		ptr := (*KeyValueList)(elem.value.PtrVal())
		if ptr != nil {
			keyValueListPool.Release(ptr)
		}
	}

	// Zero-initialize the released element.
	*elem = AnyValue{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}

// ====================== ArrayValue message implementation ======================

type ArrayValue struct {
	_protoMessage protomessage.ProtoMessage
	_flags        ArrayValueFlags

	values []*AnyValue
}

func UnmarshalArrayValue(bytes []byte) (*ArrayValue, error) {
	m := arrayValuePool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *ArrayValue) Free() {
	arrayValuePool.Release(m)
}

// ArrayValueFlags is the type of the bit flags.
type ArrayValueFlags uint8

// Bitmasks that indicate that the particular nested message is decoded.
const flagArrayValueValuesDecoded ArrayValueFlags = 0x1

// HasValues returns true if the values is present.
func (m *ArrayValue) HasValues() bool {
	return len(m.values) > 0
}

// Values returns the value of the values.
func (m *ArrayValue) Values() []*AnyValue {
	if m._flags&flagArrayValueValuesDecoded == 0 {
		// Decode nested message(s).
		for i := range m.values {
			// TODO: decide how to handle decoding errors.
			_ = m.values[i].decode()
		}
		m._flags |= flagArrayValueValuesDecoded
	}
	return m.values
}

// SetValues sets the value of the values.
func (m *ArrayValue) SetValues(v []*AnyValue) {
	m.values = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.values {
		elem._protoMessage.Parent = &m._protoMessage
	}

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *ArrayValue) ValuesRemoveIf(f func(*AnyValue) bool) {
	// Call getter to load the field.
	m.Values()

	newLen := 0
	for i := 0; i < len(m.values); i++ {
		if f(m.values[i]) {
			continue
		}
		if newLen == i {
			// Nothing to move, element is at the right place.
			newLen++
			continue
		}
		m.values[newLen] = m.values[i]
		newLen++
	}
	if newLen != len(m.values) {
		m.values = m.values[:newLen]
		// Mark this message modified, if not already.
		m._protoMessage.MarkModified()
	}
}

func (m *ArrayValue) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// If the user makes a mistake and takes a copy of this struct before decoding it
	// the "decoded" flag will incorrectly set on the copy, not the original, but nested
	// messages will be marked as "decoded". Next time we try to access the nested
	// message via getter func this decode() func will be called again and will overwrite
	// nested message values, but the getter func will no longer attempt to decode it
	// because the flag "decoded" flag is incorrectly set on nested message.
	// This will result in incorrect state of nested message returned by getter.
	// To make sure we correctly decode even after this mistake we reset all "decoded"
	// and "presence" flags here.
	m._flags = 0

	// Count all repeated fields. We need one counter per field.
	valuesCount := 0
	err := molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 1 {
				valuesCount++
			}
		},
	)
	if err != nil {
		return err
	}

	// Pre-allocate slices for repeated fields.
	m.values = anyValuePool.GetSlice(valuesCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// Set slice indexes to 0 to begin iterating over repeated fields.
	valuesCount = 0
	// Iterate and decode the fields.
	err2 := molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode "values".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.values[valuesCount]
				valuesCount++
				elem._protoMessage.Parent = &m._protoMessage
				elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
			}
			return true, nil
		},
	)
	if err2 != nil {
		return err2
	}
	return nil
}

var preparedArrayValueValues = molecule.PrepareEmbeddedField(1)

func (m *ArrayValue) Marshal(ps *molecule.ProtoStream) error {
	if m._protoMessage.IsModified() {
		// Marshal "values".
		for _, elem := range m.values {
			token := ps.BeginEmbedded()
			if err := elem.Marshal(ps); err != nil {
				return err
			}
			ps.EndEmbeddedPrepared(token, preparedArrayValueValues)
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))
	}
	return nil
}

// Pool of ArrayValue structs.
type arrayValuePoolType struct {
	pool []*ArrayValue
	mux  sync.Mutex
}

var arrayValuePool = arrayValuePoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *arrayValuePoolType) Get() *ArrayValue {
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
	return &ArrayValue{}
}

func (p *arrayValuePoolType) GetSlice(count int) []*ArrayValue {
	// Create a new slice.
	r := make([]*ArrayValue, count)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return r
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]ArrayValue, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *arrayValuePoolType) ReleaseSlice(slice []*ArrayValue) {
	for _, elem := range slice {
		// Release nested values recursively to their pool.
		anyValuePool.ReleaseSlice(elem.values)

		// Zero-initialize the released element.
		*elem = ArrayValue{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *arrayValuePoolType) Release(elem *ArrayValue) {
	// Release nested values recursively to their pool.
	anyValuePool.ReleaseSlice(elem.values)

	// Zero-initialize the released element.
	*elem = ArrayValue{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}

// ====================== KeyValueList message implementation ======================

type KeyValueList struct {
	_protoMessage protomessage.ProtoMessage
	_flags        KeyValueListFlags

	values []*KeyValue
}

func UnmarshalKeyValueList(bytes []byte) (*KeyValueList, error) {
	m := keyValueListPool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *KeyValueList) Free() {
	keyValueListPool.Release(m)
}

// KeyValueListFlags is the type of the bit flags.
type KeyValueListFlags uint8

// Bitmasks that indicate that the particular nested message is decoded.
const flagKeyValueListValuesDecoded KeyValueListFlags = 0x1

// HasValues returns true if the values is present.
func (m *KeyValueList) HasValues() bool {
	return len(m.values) > 0
}

// Values returns the value of the values.
func (m *KeyValueList) Values() []*KeyValue {
	if m._flags&flagKeyValueListValuesDecoded == 0 {
		// Decode nested message(s).
		for i := range m.values {
			// TODO: decide how to handle decoding errors.
			_ = m.values[i].decode()
		}
		m._flags |= flagKeyValueListValuesDecoded
	}
	return m.values
}

// SetValues sets the value of the values.
func (m *KeyValueList) SetValues(v []*KeyValue) {
	m.values = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.values {
		elem._protoMessage.Parent = &m._protoMessage
	}

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *KeyValueList) ValuesRemoveIf(f func(*KeyValue) bool) {
	// Call getter to load the field.
	m.Values()

	newLen := 0
	for i := 0; i < len(m.values); i++ {
		if f(m.values[i]) {
			continue
		}
		if newLen == i {
			// Nothing to move, element is at the right place.
			newLen++
			continue
		}
		m.values[newLen] = m.values[i]
		newLen++
	}
	if newLen != len(m.values) {
		m.values = m.values[:newLen]
		// Mark this message modified, if not already.
		m._protoMessage.MarkModified()
	}
}

func (m *KeyValueList) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// If the user makes a mistake and takes a copy of this struct before decoding it
	// the "decoded" flag will incorrectly set on the copy, not the original, but nested
	// messages will be marked as "decoded". Next time we try to access the nested
	// message via getter func this decode() func will be called again and will overwrite
	// nested message values, but the getter func will no longer attempt to decode it
	// because the flag "decoded" flag is incorrectly set on nested message.
	// This will result in incorrect state of nested message returned by getter.
	// To make sure we correctly decode even after this mistake we reset all "decoded"
	// and "presence" flags here.
	m._flags = 0

	// Count all repeated fields. We need one counter per field.
	valuesCount := 0
	err := molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 1 {
				valuesCount++
			}
		},
	)
	if err != nil {
		return err
	}

	// Pre-allocate slices for repeated fields.
	m.values = keyValuePool.GetSlice(valuesCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// Set slice indexes to 0 to begin iterating over repeated fields.
	valuesCount = 0
	// Iterate and decode the fields.
	err2 := molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode "values".
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.values[valuesCount]
				valuesCount++
				elem._protoMessage.Parent = &m._protoMessage
				elem._protoMessage.Bytes = protomessage.BytesViewFromBytes(v)
			}
			return true, nil
		},
	)
	if err2 != nil {
		return err2
	}
	return nil
}

var preparedKeyValueListValues = molecule.PrepareEmbeddedField(1)

func (m *KeyValueList) Marshal(ps *molecule.ProtoStream) error {
	if m._protoMessage.IsModified() {
		// Marshal "values".
		for _, elem := range m.values {
			token := ps.BeginEmbedded()
			if err := elem.Marshal(ps); err != nil {
				return err
			}
			ps.EndEmbeddedPrepared(token, preparedKeyValueListValues)
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))
	}
	return nil
}

// Pool of KeyValueList structs.
type keyValueListPoolType struct {
	pool []*KeyValueList
	mux  sync.Mutex
}

var keyValueListPool = keyValueListPoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *keyValueListPoolType) Get() *KeyValueList {
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
	return &KeyValueList{}
}

func (p *keyValueListPoolType) GetSlice(count int) []*KeyValueList {
	// Create a new slice.
	r := make([]*KeyValueList, count)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return r
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]KeyValueList, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *keyValueListPoolType) ReleaseSlice(slice []*KeyValueList) {
	for _, elem := range slice {
		// Release nested values recursively to their pool.
		keyValuePool.ReleaseSlice(elem.values)

		// Zero-initialize the released element.
		*elem = KeyValueList{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *keyValueListPoolType) Release(elem *KeyValueList) {
	// Release nested values recursively to their pool.
	keyValuePool.ReleaseSlice(elem.values)

	// Zero-initialize the released element.
	*elem = KeyValueList{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}

// ====================== PlainMessage message implementation ======================

type PlainMessage struct {
	_protoMessage protomessage.ProtoMessage
	_flags        PlainMessageFlags

	key   string
	value string
}

func UnmarshalPlainMessage(bytes []byte) (*PlainMessage, error) {
	m := plainMessagePool.Get()
	m._protoMessage.Bytes = protomessage.BytesViewFromBytes(bytes)
	if err := m.decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *PlainMessage) Free() {
	plainMessagePool.Release(m)
}

// PlainMessageFlags is the type of the bit flags.
type PlainMessageFlags uint8

// Bitmasks that indicate that the particular field is present.
const flagPlainMessageKeyPresent PlainMessageFlags = 0x1
const flagPlainMessageValuePresent PlainMessageFlags = 0x2

// HasKey returns true if the key is present.
func (m *PlainMessage) HasKey() bool {
	return m._flags&flagPlainMessageKeyPresent != 0
}

// Key returns the value of the key.
func (m *PlainMessage) Key() string {
	return m.key
}

// SetKey sets the value of the key.
func (m *PlainMessage) SetKey(v string) {
	m.key = v
	m._flags |= flagPlainMessageKeyPresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

// HasValue returns true if the value is present.
func (m *PlainMessage) HasValue() bool {
	return m._flags&flagPlainMessageValuePresent != 0
}

// Value returns the value of the value.
func (m *PlainMessage) Value() string {
	return m.value
}

// SetValue sets the value of the value.
func (m *PlainMessage) SetValue(v string) {
	m.value = v
	m._flags |= flagPlainMessageValuePresent

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}

func (m *PlainMessage) decode() error {
	buf := codec.NewBuffer(protomessage.BytesFromBytesView(m._protoMessage.Bytes))

	// If the user makes a mistake and takes a copy of this struct before decoding it
	// the "decoded" flag will incorrectly set on the copy, not the original, but nested
	// messages will be marked as "decoded". Next time we try to access the nested
	// message via getter func this decode() func will be called again and will overwrite
	// nested message values, but the getter func will no longer attempt to decode it
	// because the flag "decoded" flag is incorrectly set on nested message.
	// This will result in incorrect state of nested message returned by getter.
	// To make sure we correctly decode even after this mistake we reset all "decoded"
	// and "presence" flags here.
	m._flags = 0

	// Iterate and decode the fields.
	err2 := molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode "key".
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.key = v
				m._flags |= flagPlainMessageKeyPresent
			case 2:
				// Decode "value".
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.value = v
				m._flags |= flagPlainMessageValuePresent
			}
			return true, nil
		},
	)
	if err2 != nil {
		return err2
	}
	return nil
}

var preparedPlainMessageKey = molecule.PrepareStringField(1)
var preparedPlainMessageValue = molecule.PrepareStringField(2)

func (m *PlainMessage) Marshal(ps *molecule.ProtoStream) error {
	if m._protoMessage.IsModified() {
		// Marshal "key".
		if m._flags&flagPlainMessageKeyPresent != 0 {
			ps.StringPrepared(preparedPlainMessageKey, m.key)
		}
		// Marshal "value".
		if m._flags&flagPlainMessageValuePresent != 0 {
			ps.StringPrepared(preparedPlainMessageValue, m.value)
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(protomessage.BytesFromBytesView(m._protoMessage.Bytes))
	}
	return nil
}

// Pool of PlainMessage structs.
type plainMessagePoolType struct {
	pool []*PlainMessage
	mux  sync.Mutex
}

var plainMessagePool = plainMessagePoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *plainMessagePoolType) Get() *PlainMessage {
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
	return &PlainMessage{}
}

func (p *plainMessagePoolType) GetSlice(count int) []*PlainMessage {
	// Create a new slice.
	r := make([]*PlainMessage, count)

	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Copy the elements from the end of the pool.
		copy(r, p.pool[len(p.pool)-count:])

		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]

		return r
	}

	// Initialize with what remains in the pool.
	copied := copy(r, p.pool)
	p.pool = nil

	if copied < count {
		// Create remaining elements.
		storage := make([]PlainMessage, count-copied)
		j := 0
		for ; copied < count; copied++ {
			r[copied] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *plainMessagePoolType) ReleaseSlice(slice []*PlainMessage) {
	for _, elem := range slice {

		// Zero-initialize the released element.
		*elem = PlainMessage{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *plainMessagePoolType) Release(elem *PlainMessage) {

	// Zero-initialize the released element.
	*elem = PlainMessage{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}