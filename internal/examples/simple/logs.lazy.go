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
	m := &LogsData{protoMessage: lazyproto.ProtoMessage{Bytes: bytes}}
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
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
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
				*m.resourceLogs[resourceLogsCount] = ResourceLogs{
					protoMessage: lazyproto.ProtoMessage{
						Parent: &m.protoMessage, Bytes: v,
					},
				}
				resourceLogsCount++
			}
			return true, nil
		},
	)
}

func (m *LogsData) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		// Marshal resourceLogs
		for _, elem := range m.resourceLogs {
			token := ps.BeginEmbedded()
			elem.Marshal(ps)
			ps.EndEmbedded(token, 1)
		}
	} else {
		// Message is unchanged. Used original bytes.
		ps.Raw(m.protoMessage.Bytes)
	}
	return nil
}

type ResourceLogs struct {
	protoMessage lazyproto.ProtoMessage
	resource     *Resource
	scopeLogs    []*ScopeLogs
}

const resourceLogsResourceDecoded = 2
const resourceLogsScopeLogsDecoded = 4

func (m *ResourceLogs) Resource() *Resource {
	if m.protoMessage.Flags&resourceLogsResourceDecoded == 0 {
		m.resource.decode()
		m.protoMessage.Flags |= resourceLogsResourceDecoded
	}
	return m.resource
}

func (m *ResourceLogs) ScopeLogs() []*ScopeLogs {
	if m.protoMessage.Flags&resourceLogsScopeLogsDecoded == 0 {
		for i := range m.scopeLogs {
			m.scopeLogs[i].decode()
		}
		m.protoMessage.Flags |= resourceLogsScopeLogsDecoded
	}
	return m.scopeLogs
}

func (m *ResourceLogs) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)

	lrCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 2 {
				lrCount++
			}
		},
	)
	//m.scopeLogs = make([]ScopeLogs, 0, lrCount)
	m.scopeLogs = scopeLogsPool.GetSlice(lrCount)

	lrIndex := 0
	buf.Reset(m.protoMessage.Bytes)
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.resource = resourcePool.Get()
				*m.resource = Resource{
					protoMessage: lazyproto.ProtoMessage{
						Parent: &m.protoMessage, Bytes: v,
					},
				}

			case 2:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				sl := m.scopeLogs[lrIndex]
				lrIndex++
				*sl = ScopeLogs{
					protoMessage: lazyproto.ProtoMessage{
						Parent: &m.protoMessage, Bytes: v,
					},
				}
			}
			return true, nil
		},
	)
}

func (m *ResourceLogs) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		if m.resource != nil {
			token := ps.BeginEmbedded()
			m.resource.Marshal(ps)
			ps.EndEmbedded(token, 1)
		}
		for _, logs := range m.scopeLogs {
			token := ps.BeginEmbedded()
			logs.Marshal(ps)
			ps.EndEmbedded(token, 2)
		}
	} else {
		ps.Raw(m.protoMessage.Bytes)
	}
	return nil
}

type Resource struct {
	protoMessage           lazyproto.ProtoMessage
	attributes             []*KeyValue
	DroppedAttributesCount uint32
}

const resourceAttributesDecoded = 2

func (m *Resource) Attributes() []*KeyValue {
	if m.protoMessage.Flags&resourceAttributesDecoded == 0 {
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.protoMessage.Flags |= resourceAttributesDecoded
	}
	return m.attributes
}

func (m *Resource) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)
	attrCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 2 {
				attrCount++
			}
		},
	)
	m.attributes = keyValuePool.GetSlice(attrCount)

	attrIndex := 0
	buf.Reset(m.protoMessage.Bytes)
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				kv := m.attributes[attrIndex]
				attrIndex++
				*kv = KeyValue{
					protoMessage: lazyproto.ProtoMessage{
						Parent: &m.protoMessage, Bytes: v,
					},
				}

			case 2:
				v, err := value.AsUint32()
				if err != nil {
					return false, err
				}
				m.DroppedAttributesCount = v
			}
			return true, nil
		},
	)
}

var resourceAttrKeyPrepared = molecule.PrepareEmbeddedField(1)
var resourceDroppedKeyPrepared = molecule.PrepareUint32Field(2)

func (m *Resource) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		for _, attr := range m.attributes {
			token := ps.BeginEmbedded()
			attr.Marshal(ps)
			//ps.EndEmbeddedPrepared(token, 1)
			ps.EndEmbeddedPrepared(token, resourceAttrKeyPrepared)
		}
		ps.Uint32Prepared(resourceDroppedKeyPrepared, m.DroppedAttributesCount)
	} else {
		ps.Raw(m.protoMessage.Bytes)
	}
	return nil
}

type ScopeLogs struct {
	protoMessage lazyproto.ProtoMessage
	logRecords   []*LogRecord
}

const scopeLogsLogRecordsDecoded = 2

func (m *ScopeLogs) LogRecords() []*LogRecord {
	if m.protoMessage.Flags&scopeLogsLogRecordsDecoded == 0 {
		for i := range m.logRecords {
			m.logRecords[i].decode()
		}
		m.protoMessage.Flags |= scopeLogsLogRecordsDecoded
	}
	return m.logRecords
}

func (m *ScopeLogs) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)

	// Count all repeated fields. We need one counter per field.
	logRecordsCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 1 {
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
				// Decode logRecords.
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				lr := m.logRecords[logRecordsCount]
				*lr = LogRecord{
					protoMessage: lazyproto.ProtoMessage{
						Parent: &m.protoMessage, Bytes: v,
					},
				}
				logRecordsCount++
			}
			return true, nil
		},
	)
}

func (m *ScopeLogs) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		// Marshal logRecords
		for _, elem := range m.logRecords {
			token := ps.BeginEmbedded()
			elem.Marshal(ps)
			ps.EndEmbedded(token, 1)
		}
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
	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Cut the required slice from the end of the pool.
		r := p.pool[len(p.pool)-count:]
		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]
		return r
	}

	// Create a new slice.
	r := make([]*ScopeLogs, count)

	// Initialize with what remains in the pool.
	i := 0
	for ; i < len(p.pool); i++ {
		r[i] = p.pool[i]
	}
	p.pool = nil

	if i < count {
		// Create remaining elements.
		storage := make([]ScopeLogs, count-i)
		j := 0
		for ; i < count; i++ {
			r[i] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *scopeLogsPoolType) ReleaseSlice(slice []*ScopeLogs) {
	for _, elem := range slice {
		// Release nested logRecords recursively to their pool.
		logRecordPool.ReleaseSlice(elem.logRecords)

		// Zero-initialize the released element.
		// *elem = ScopeLogs{}
		elem = elem
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *scopeLogsPoolType) Release(elem *ScopeLogs) {
	// Release nested logRecords recursively to their pool.
	logRecordPool.ReleaseSlice(elem.logRecords)

	// Zero-initialize the released element.
	*elem = ScopeLogs{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, elem)
}

// ====================== Generated for message LogRecord ======================

type LogRecord struct {
	protoMessage           lazyproto.ProtoMessage
	timeUnixNano           uint64
	attributes             []*KeyValue
	droppedAttributesCount uint32
}

func NewLogRecord(bytes []byte) *LogRecord {
	m := &LogRecord{protoMessage: lazyproto.ProtoMessage{Bytes: bytes}}
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
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *LogRecord) DroppedAttributesCount() uint32 {
	return m.droppedAttributesCount
}

func (m *LogRecord) SetDroppedAttributesCount(v uint32) {
	m.droppedAttributesCount = v
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *LogRecord) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)
	attributesCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 2 {
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
			case 2:
				// Decode attributes.
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				*m.attributes[attributesCount] = KeyValue{
					protoMessage: lazyproto.ProtoMessage{
						Parent: &m.protoMessage, Bytes: v,
					},
				}
				attributesCount++
			case 3:
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

var preparedLogRecordTimeUnixNano = molecule.PrepareFixed64Field(1)
var logRecordAttrPrepared = molecule.PrepareEmbeddedField(2)
var preparedLogRecordDroppedAttributesCount = molecule.PrepareUint32Field(3)

func (m *LogRecord) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		// Marshal timeUnixNano
		ps.Fixed64Prepared(preparedLogRecordTimeUnixNano, m.timeUnixNano)
		// Marshal attributes
		for _, elem := range m.attributes {
			token := ps.BeginEmbedded()
			elem.Marshal(ps)
			ps.EndEmbeddedPrepared(token, logRecordAttrPrepared)
		}
		// Marshal droppedAttributesCount
		ps.Uint32Prepared(
			preparedLogRecordDroppedAttributesCount, m.droppedAttributesCount,
		)
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
	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Cut the required slice from the end of the pool.
		r := p.pool[len(p.pool)-count:]
		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]
		return r
	}

	// Create a new slice.
	r := make([]*LogRecord, count)

	// Initialize with what remains in the pool.
	i := 0
	for ; i < len(p.pool); i++ {
		r[i] = p.pool[i]
	}
	p.pool = nil

	if i < count {
		// Create remaining elements.
		storage := make([]LogRecord, count-i)
		j := 0
		for ; i < count; i++ {
			r[i] = &storage[j]
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
	m := &KeyValue{protoMessage: lazyproto.ProtoMessage{Bytes: bytes}}
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
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified == 0 {
		m.protoMessage.MarkModified()
	}
}

func (m *KeyValue) Value() string {
	return m.value
}

func (m *KeyValue) SetValue(v string) {
	m.value = v
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
	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Cut the required slice from the end of the pool.
		r := p.pool[len(p.pool)-count:]
		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]
		return r
	}

	// Create a new slice.
	r := make([]*KeyValue, count)

	// Initialize with what remains in the pool.
	i := 0
	for ; i < len(p.pool); i++ {
		r[i] = p.pool[i]
	}
	p.pool = nil

	if i < count {
		// Create remaining elements.
		storage := make([]KeyValue, count-i)
		j := 0
		for ; i < count; i++ {
			r[i] = &storage[j]
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

	//for _, e1 := range p.pool {
	//	for _, e2 := range slice {
	//		if e1 == e2 {
	//			panic("duplicate")
	//		}
	//	}
	//}

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

//elem := m.attributes[attributesCount]
//elem.protoMessage = lazyproto.ProtoMessage{
//Parent: &m.protoMessage, Bytes: v,
//}
