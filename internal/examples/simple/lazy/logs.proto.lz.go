package simple

import (
	"sync"

	lazyproto "github.com/tigrannajaryan/exp-lazyproto"
	"github.com/tigrannajaryan/molecule"
	"github.com/tigrannajaryan/molecule/src/codec"
)

// ====================== Generated for message LogsData ======================

// LogsData contains all log data
type LogsData struct {
	protoMessage lazyproto.ProtoMessage
	// List of ResourceLogs
	resourceLogs []*ResourceLogs
}

func NewLogsData(bytes []byte) *LogsData {
	m := logsDataPool.Get()
	m.protoMessage.Bytes = bytes
	m.decode()
	return m
}

func (m *LogsData) Free() {
	logsDataPool.Release(m)
}

// Bitmasks that indicate that the particular nested message is decoded.
const flagLogsDataResourceLogsDecoded = 0x0000000000000002

func (m *LogsData) ResourceLogs() []*ResourceLogs {
	if m.protoMessage.Flags&flagLogsDataResourceLogsDecoded == 0 {
		// Decode nested message(s).
		for i := range m.resourceLogs {
			m.resourceLogs[i].decode()
		}
		m.protoMessage.Flags |= flagLogsDataResourceLogsDecoded
	}
	return m.resourceLogs
}

func (m *LogsData) SetResourceLogs(v []*ResourceLogs) {
	m.resourceLogs = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.resourceLogs {
		elem.protoMessage.Parent = &m.protoMessage
	}

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
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
		if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
			m.protoMessage.MarkModified()
		}
	}
}

func (m *LogsData) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)

	// Count all repeated fields. We need one counter per field.
	resourceLogsCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 1 {
				resourceLogsCount++
			}
		},
	)

	// Pre-allocate slices for repeated fields.
	m.resourceLogs = resourceLogsPool.GetSlice(resourceLogsCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.protoMessage.Bytes)

	// Set slice indexes to 0 to begin iterating over repeated fields.
	resourceLogsCount = 0

	// Iterate and decode the fields.
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode resourceLogs.
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.resourceLogs[resourceLogsCount]
				resourceLogsCount++
				elem.protoMessage.Parent = &m.protoMessage
				elem.protoMessage.Bytes = v
			}
			return true, nil
		},
	)
}

var preparedLogsDataResourceLogs = molecule.PrepareEmbeddedField(1)

func (m *LogsData) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		// Marshal resourceLogs
		for _, elem := range m.resourceLogs {
			token := ps.BeginEmbedded()
			elem.Marshal(ps)
			ps.EndEmbeddedPrepared(token, preparedLogsDataResourceLogs)
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(m.protoMessage.Bytes)
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

// ====================== Generated for message ResourceLogs ======================

type ResourceLogs struct {
	protoMessage lazyproto.ProtoMessage
	// The Resource
	resource *Resource
	// List of ScopeLogs
	scopeLogs []*ScopeLogs
	schemaUrl string
}

func NewResourceLogs(bytes []byte) *ResourceLogs {
	m := resourceLogsPool.Get()
	m.protoMessage.Bytes = bytes
	m.decode()
	return m
}

func (m *ResourceLogs) Free() {
	resourceLogsPool.Release(m)
}

// Bitmasks that indicate that the particular nested message is decoded.
const flagResourceLogsResourceDecoded = 0x0000000000000002
const flagResourceLogsScopeLogsDecoded = 0x0000000000000004

func (m *ResourceLogs) Resource() *Resource {
	if m.protoMessage.Flags&flagResourceLogsResourceDecoded == 0 {
		// Decode nested message(s).
		if m.resource != nil {
			m.resource.decode()
		}
		m.protoMessage.Flags |= flagResourceLogsResourceDecoded
	}
	return m.resource
}

func (m *ResourceLogs) SetResource(v *Resource) {
	m.resource = v

	// Make sure the field's Parent points to this message.
	m.resource.protoMessage.Parent = &m.protoMessage

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *ResourceLogs) ScopeLogs() []*ScopeLogs {
	if m.protoMessage.Flags&flagResourceLogsScopeLogsDecoded == 0 {
		// Decode nested message(s).
		for i := range m.scopeLogs {
			m.scopeLogs[i].decode()
		}
		m.protoMessage.Flags |= flagResourceLogsScopeLogsDecoded
	}
	return m.scopeLogs
}

func (m *ResourceLogs) SetScopeLogs(v []*ScopeLogs) {
	m.scopeLogs = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.scopeLogs {
		elem.protoMessage.Parent = &m.protoMessage
	}

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
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
		if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
			m.protoMessage.MarkModified()
		}
	}
}

func (m *ResourceLogs) SchemaUrl() string {
	return m.schemaUrl
}

func (m *ResourceLogs) SetSchemaUrl(v string) {
	m.schemaUrl = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *ResourceLogs) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)

	// Count all repeated fields. We need one counter per field.
	scopeLogsCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 2 {
				scopeLogsCount++
			}
		},
	)

	// Pre-allocate slices for repeated fields.
	m.scopeLogs = scopeLogsPool.GetSlice(scopeLogsCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.protoMessage.Bytes)

	// Set slice indexes to 0 to begin iterating over repeated fields.
	scopeLogsCount = 0

	// Iterate and decode the fields.
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode resource.
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.resource = resourcePool.Get()
				m.resource.protoMessage.Parent = &m.protoMessage
				m.resource.protoMessage.Bytes = v
			case 2:
				// Decode scopeLogs.
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.scopeLogs[scopeLogsCount]
				scopeLogsCount++
				elem.protoMessage.Parent = &m.protoMessage
				elem.protoMessage.Bytes = v
			case 3:
				// Decode schemaUrl.
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.schemaUrl = v
			}
			return true, nil
		},
	)
}

var preparedResourceLogsResource = molecule.PrepareEmbeddedField(1)
var preparedResourceLogsScopeLogs = molecule.PrepareEmbeddedField(2)
var preparedResourceLogsSchemaUrl = molecule.PrepareStringField(3)

func (m *ResourceLogs) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		// Marshal resource
		if m.resource != nil {
			token := ps.BeginEmbedded()
			m.resource.Marshal(ps)
			ps.EndEmbeddedPrepared(token, preparedResourceLogsResource)
		}
		// Marshal scopeLogs
		for _, elem := range m.scopeLogs {
			token := ps.BeginEmbedded()
			elem.Marshal(ps)
			ps.EndEmbeddedPrepared(token, preparedResourceLogsScopeLogs)
		}
		// Marshal schemaUrl
		ps.StringPrepared(preparedResourceLogsSchemaUrl, m.schemaUrl)
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(m.protoMessage.Bytes)
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

// ====================== Generated for message Resource ======================

type Resource struct {
	protoMessage           lazyproto.ProtoMessage
	attributes             []*KeyValue
	droppedAttributesCount uint32
}

func NewResource(bytes []byte) *Resource {
	m := resourcePool.Get()
	m.protoMessage.Bytes = bytes
	m.decode()
	return m
}

func (m *Resource) Free() {
	resourcePool.Release(m)
}

// Bitmasks that indicate that the particular nested message is decoded.
const flagResourceAttributesDecoded = 0x0000000000000002

func (m *Resource) Attributes() []*KeyValue {
	if m.protoMessage.Flags&flagResourceAttributesDecoded == 0 {
		// Decode nested message(s).
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.protoMessage.Flags |= flagResourceAttributesDecoded
	}
	return m.attributes
}

func (m *Resource) SetAttributes(v []*KeyValue) {
	m.attributes = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.attributes {
		elem.protoMessage.Parent = &m.protoMessage
	}

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
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
		if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
			m.protoMessage.MarkModified()
		}
	}
}

func (m *Resource) DroppedAttributesCount() uint32 {
	return m.droppedAttributesCount
}

func (m *Resource) SetDroppedAttributesCount(v uint32) {
	m.droppedAttributesCount = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *Resource) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)

	// Count all repeated fields. We need one counter per field.
	attributesCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 1 {
				attributesCount++
			}
		},
	)

	// Pre-allocate slices for repeated fields.
	m.attributes = keyValuePool.GetSlice(attributesCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.protoMessage.Bytes)

	// Set slice indexes to 0 to begin iterating over repeated fields.
	attributesCount = 0

	// Iterate and decode the fields.
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode attributes.
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.attributes[attributesCount]
				attributesCount++
				elem.protoMessage.Parent = &m.protoMessage
				elem.protoMessage.Bytes = v
			case 2:
				// Decode droppedAttributesCount.
				v, err := value.AsUint32()
				if err != nil {
					return false, err
				}
				m.droppedAttributesCount = v
			}
			return true, nil
		},
	)
}

var preparedResourceAttributes = molecule.PrepareEmbeddedField(1)
var preparedResourceDroppedAttributesCount = molecule.PrepareUint32Field(2)

func (m *Resource) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		// Marshal attributes
		for _, elem := range m.attributes {
			token := ps.BeginEmbedded()
			elem.Marshal(ps)
			ps.EndEmbeddedPrepared(token, preparedResourceAttributes)
		}
		// Marshal droppedAttributesCount
		ps.Uint32Prepared(preparedResourceDroppedAttributesCount, m.droppedAttributesCount)
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(m.protoMessage.Bytes)
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

// ====================== Generated for message ScopeLogs ======================

// A collection of Logs produced by a Scope.
type ScopeLogs struct {
	protoMessage lazyproto.ProtoMessage
	scope        *InstrumentationScope
	// A list of log records.
	logRecords []*LogRecord
	// This schema_url applies to all logs in the "logs" field.
	schemaUrl string
}

func NewScopeLogs(bytes []byte) *ScopeLogs {
	m := scopeLogsPool.Get()
	m.protoMessage.Bytes = bytes
	m.decode()
	return m
}

func (m *ScopeLogs) Free() {
	scopeLogsPool.Release(m)
}

// Bitmasks that indicate that the particular nested message is decoded.
const flagScopeLogsScopeDecoded = 0x0000000000000002
const flagScopeLogsLogRecordsDecoded = 0x0000000000000004

func (m *ScopeLogs) Scope() *InstrumentationScope {
	if m.protoMessage.Flags&flagScopeLogsScopeDecoded == 0 {
		// Decode nested message(s).
		if m.scope != nil {
			m.scope.decode()
		}
		m.protoMessage.Flags |= flagScopeLogsScopeDecoded
	}
	return m.scope
}

func (m *ScopeLogs) SetScope(v *InstrumentationScope) {
	m.scope = v

	// Make sure the field's Parent points to this message.
	m.scope.protoMessage.Parent = &m.protoMessage

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *ScopeLogs) LogRecords() []*LogRecord {
	if m.protoMessage.Flags&flagScopeLogsLogRecordsDecoded == 0 {
		// Decode nested message(s).
		for i := range m.logRecords {
			m.logRecords[i].decode()
		}
		m.protoMessage.Flags |= flagScopeLogsLogRecordsDecoded
	}
	return m.logRecords
}

func (m *ScopeLogs) SetLogRecords(v []*LogRecord) {
	m.logRecords = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.logRecords {
		elem.protoMessage.Parent = &m.protoMessage
	}

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
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
		if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
			m.protoMessage.MarkModified()
		}
	}
}

func (m *ScopeLogs) SchemaUrl() string {
	return m.schemaUrl
}

func (m *ScopeLogs) SetSchemaUrl(v string) {
	m.schemaUrl = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *ScopeLogs) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)

	// Count all repeated fields. We need one counter per field.
	logRecordsCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 2 {
				logRecordsCount++
			}
		},
	)

	// Pre-allocate slices for repeated fields.
	m.logRecords = logRecordPool.GetSlice(logRecordsCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.protoMessage.Bytes)

	// Set slice indexes to 0 to begin iterating over repeated fields.
	logRecordsCount = 0

	// Iterate and decode the fields.
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode scope.
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.scope = instrumentationScopePool.Get()
				m.scope.protoMessage.Parent = &m.protoMessage
				m.scope.protoMessage.Bytes = v
			case 2:
				// Decode logRecords.
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.logRecords[logRecordsCount]
				logRecordsCount++
				elem.protoMessage.Parent = &m.protoMessage
				elem.protoMessage.Bytes = v
			case 3:
				// Decode schemaUrl.
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.schemaUrl = v
			}
			return true, nil
		},
	)
}

var preparedScopeLogsScope = molecule.PrepareEmbeddedField(1)
var preparedScopeLogsLogRecords = molecule.PrepareEmbeddedField(2)
var preparedScopeLogsSchemaUrl = molecule.PrepareStringField(3)

func (m *ScopeLogs) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		// Marshal scope
		if m.scope != nil {
			token := ps.BeginEmbedded()
			m.scope.Marshal(ps)
			ps.EndEmbeddedPrepared(token, preparedScopeLogsScope)
		}
		// Marshal logRecords
		for _, elem := range m.logRecords {
			token := ps.BeginEmbedded()
			elem.Marshal(ps)
			ps.EndEmbeddedPrepared(token, preparedScopeLogsLogRecords)
		}
		// Marshal schemaUrl
		ps.StringPrepared(preparedScopeLogsSchemaUrl, m.schemaUrl)
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(m.protoMessage.Bytes)
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

// ====================== Generated for message InstrumentationScope ======================

type InstrumentationScope struct {
	protoMessage           lazyproto.ProtoMessage
	name                   string
	version                string
	attributes             []*KeyValue
	droppedAttributesCount uint32
}

func NewInstrumentationScope(bytes []byte) *InstrumentationScope {
	m := instrumentationScopePool.Get()
	m.protoMessage.Bytes = bytes
	m.decode()
	return m
}

func (m *InstrumentationScope) Free() {
	instrumentationScopePool.Release(m)
}

// Bitmasks that indicate that the particular nested message is decoded.
const flagInstrumentationScopeAttributesDecoded = 0x0000000000000002

func (m *InstrumentationScope) Name() string {
	return m.name
}

func (m *InstrumentationScope) SetName(v string) {
	m.name = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *InstrumentationScope) Version() string {
	return m.version
}

func (m *InstrumentationScope) SetVersion(v string) {
	m.version = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *InstrumentationScope) Attributes() []*KeyValue {
	if m.protoMessage.Flags&flagInstrumentationScopeAttributesDecoded == 0 {
		// Decode nested message(s).
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.protoMessage.Flags |= flagInstrumentationScopeAttributesDecoded
	}
	return m.attributes
}

func (m *InstrumentationScope) SetAttributes(v []*KeyValue) {
	m.attributes = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.attributes {
		elem.protoMessage.Parent = &m.protoMessage
	}

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
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
		if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
			m.protoMessage.MarkModified()
		}
	}
}

func (m *InstrumentationScope) DroppedAttributesCount() uint32 {
	return m.droppedAttributesCount
}

func (m *InstrumentationScope) SetDroppedAttributesCount(v uint32) {
	m.droppedAttributesCount = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *InstrumentationScope) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)

	// Count all repeated fields. We need one counter per field.
	attributesCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 3 {
				attributesCount++
			}
		},
	)

	// Pre-allocate slices for repeated fields.
	m.attributes = keyValuePool.GetSlice(attributesCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.protoMessage.Bytes)

	// Set slice indexes to 0 to begin iterating over repeated fields.
	attributesCount = 0

	// Iterate and decode the fields.
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode name.
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.name = v
			case 2:
				// Decode version.
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.version = v
			case 3:
				// Decode attributes.
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.attributes[attributesCount]
				attributesCount++
				elem.protoMessage.Parent = &m.protoMessage
				elem.protoMessage.Bytes = v
			case 4:
				// Decode droppedAttributesCount.
				v, err := value.AsUint32()
				if err != nil {
					return false, err
				}
				m.droppedAttributesCount = v
			}
			return true, nil
		},
	)
}

var preparedInstrumentationScopeName = molecule.PrepareStringField(1)
var preparedInstrumentationScopeVersion = molecule.PrepareStringField(2)
var preparedInstrumentationScopeAttributes = molecule.PrepareEmbeddedField(3)
var preparedInstrumentationScopeDroppedAttributesCount = molecule.PrepareUint32Field(4)

func (m *InstrumentationScope) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		// Marshal name
		ps.StringPrepared(preparedInstrumentationScopeName, m.name)
		// Marshal version
		ps.StringPrepared(preparedInstrumentationScopeVersion, m.version)
		// Marshal attributes
		for _, elem := range m.attributes {
			token := ps.BeginEmbedded()
			elem.Marshal(ps)
			ps.EndEmbeddedPrepared(token, preparedInstrumentationScopeAttributes)
		}
		// Marshal droppedAttributesCount
		ps.Uint32Prepared(preparedInstrumentationScopeDroppedAttributesCount, m.droppedAttributesCount)
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(m.protoMessage.Bytes)
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

// ====================== Generated for message LogRecord ======================

type LogRecord struct {
	protoMessage           lazyproto.ProtoMessage
	timeUnixNano           uint64
	observedTimeUnixNano   uint64
	severityText           string
	attributes             []*KeyValue
	droppedAttributesCount uint32
	flags                  uint32
}

func NewLogRecord(bytes []byte) *LogRecord {
	m := logRecordPool.Get()
	m.protoMessage.Bytes = bytes
	m.decode()
	return m
}

func (m *LogRecord) Free() {
	logRecordPool.Release(m)
}

// Bitmasks that indicate that the particular nested message is decoded.
const flagLogRecordAttributesDecoded = 0x0000000000000002

func (m *LogRecord) TimeUnixNano() uint64 {
	return m.timeUnixNano
}

func (m *LogRecord) SetTimeUnixNano(v uint64) {
	m.timeUnixNano = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *LogRecord) ObservedTimeUnixNano() uint64 {
	return m.observedTimeUnixNano
}

func (m *LogRecord) SetObservedTimeUnixNano(v uint64) {
	m.observedTimeUnixNano = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *LogRecord) SeverityText() string {
	return m.severityText
}

func (m *LogRecord) SetSeverityText(v string) {
	m.severityText = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *LogRecord) Attributes() []*KeyValue {
	if m.protoMessage.Flags&flagLogRecordAttributesDecoded == 0 {
		// Decode nested message(s).
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.protoMessage.Flags |= flagLogRecordAttributesDecoded
	}
	return m.attributes
}

func (m *LogRecord) SetAttributes(v []*KeyValue) {
	m.attributes = v

	// Make sure the field's Parent points to this message.
	for _, elem := range m.attributes {
		elem.protoMessage.Parent = &m.protoMessage
	}

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
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
		if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
			m.protoMessage.MarkModified()
		}
	}
}

func (m *LogRecord) DroppedAttributesCount() uint32 {
	return m.droppedAttributesCount
}

func (m *LogRecord) SetDroppedAttributesCount(v uint32) {
	m.droppedAttributesCount = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *LogRecord) Flags() uint32 {
	return m.flags
}

func (m *LogRecord) SetFlags(v uint32) {
	m.flags = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *LogRecord) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)

	// Count all repeated fields. We need one counter per field.
	attributesCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 6 {
				attributesCount++
			}
		},
	)

	// Pre-allocate slices for repeated fields.
	m.attributes = keyValuePool.GetSlice(attributesCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.protoMessage.Bytes)

	// Set slice indexes to 0 to begin iterating over repeated fields.
	attributesCount = 0

	// Iterate and decode the fields.
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode timeUnixNano.
				v, err := value.AsFixed64()
				if err != nil {
					return false, err
				}
				m.timeUnixNano = v
			case 11:
				// Decode observedTimeUnixNano.
				v, err := value.AsFixed64()
				if err != nil {
					return false, err
				}
				m.observedTimeUnixNano = v
			case 3:
				// Decode severityText.
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.severityText = v
			case 6:
				// Decode attributes.
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				elem := m.attributes[attributesCount]
				attributesCount++
				elem.protoMessage.Parent = &m.protoMessage
				elem.protoMessage.Bytes = v
			case 7:
				// Decode droppedAttributesCount.
				v, err := value.AsUint32()
				if err != nil {
					return false, err
				}
				m.droppedAttributesCount = v
			case 8:
				// Decode flags.
				v, err := value.AsFixed32()
				if err != nil {
					return false, err
				}
				m.flags = v
			}
			return true, nil
		},
	)
}

var preparedLogRecordTimeUnixNano = molecule.PrepareFixed64Field(1)
var preparedLogRecordObservedTimeUnixNano = molecule.PrepareFixed64Field(11)
var preparedLogRecordSeverityText = molecule.PrepareStringField(3)
var preparedLogRecordAttributes = molecule.PrepareEmbeddedField(6)
var preparedLogRecordDroppedAttributesCount = molecule.PrepareUint32Field(7)
var preparedLogRecordFlags = molecule.PrepareFixed32Field(8)

func (m *LogRecord) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		// Marshal timeUnixNano
		ps.Fixed64Prepared(preparedLogRecordTimeUnixNano, m.timeUnixNano)
		// Marshal observedTimeUnixNano
		ps.Fixed64Prepared(preparedLogRecordObservedTimeUnixNano, m.observedTimeUnixNano)
		// Marshal severityText
		ps.StringPrepared(preparedLogRecordSeverityText, m.severityText)
		// Marshal attributes
		for _, elem := range m.attributes {
			token := ps.BeginEmbedded()
			elem.Marshal(ps)
			ps.EndEmbeddedPrepared(token, preparedLogRecordAttributes)
		}
		// Marshal droppedAttributesCount
		ps.Uint32Prepared(preparedLogRecordDroppedAttributesCount, m.droppedAttributesCount)
		// Marshal flags
		ps.Fixed32Prepared(preparedLogRecordFlags, m.flags)
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(m.protoMessage.Bytes)
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

// ====================== Generated for message KeyValue ======================

type KeyValue struct {
	protoMessage lazyproto.ProtoMessage
	key          string
	value        string
}

func NewKeyValue(bytes []byte) *KeyValue {
	m := keyValuePool.Get()
	m.protoMessage.Bytes = bytes
	m.decode()
	return m
}

func (m *KeyValue) Free() {
	keyValuePool.Release(m)
}

func (m *KeyValue) Key() string {
	return m.key
}

func (m *KeyValue) SetKey(v string) {
	m.key = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *KeyValue) Value() string {
	return m.value
}

func (m *KeyValue) SetValue(v string) {
	m.value = v

	// Mark this message modified, if not already.
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *KeyValue) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)

	// Iterate and decode the fields.
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode key.
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.key = v
			case 2:
				// Decode value.
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.value = v
			}
			return true, nil
		},
	)
}

var preparedKeyValueKey = molecule.PrepareStringField(1)
var preparedKeyValueValue = molecule.PrepareStringField(2)

func (m *KeyValue) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		// Marshal key
		ps.StringPrepared(preparedKeyValueKey, m.key)
		// Marshal value
		ps.StringPrepared(preparedKeyValueValue, m.value)
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(m.protoMessage.Bytes)
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

	// Zero-initialize the released element.
	*elem = KeyValue{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}
