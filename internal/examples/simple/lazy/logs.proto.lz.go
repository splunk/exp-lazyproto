package simple

import (
	"github.com/richardartoul/molecule"
	"github.com/richardartoul/molecule/src/codec"
)

const flagsMessageModified = 1

type ProtoMessage struct {
	flags  uint64
	parent *ProtoMessage
	bytes  []byte
}

// LogsData contains all log data
type LogsData struct {
	ProtoMessage
	// List of ResourceLogs
	resourceLogs []*ResourceLogs
}

func NewLogsData(bytes []byte) *LogsData {
	m := &LogsData{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *LogsData) decode() {
	buf := codec.NewBuffer(m.bytes)

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
	m.resourceLogs = make([]*ResourceLogs, 0, resourceLogsCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.bytes)

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
				m.resourceLogs[resourceLogsCount] = &ResourceLogs{
					ProtoMessage: ProtoMessage{bytes: v, parent: &m.ProtoMessage},
				}
				resourceLogsCount++
			}
			return true, nil
		},
	)
}

// Bitmasks that indicate that the particular nested message is decoded.
const flagLogsDataResourceLogsDecoded = 0x0000000000000002

func (m *LogsData) GetResourceLogs() []*ResourceLogs {
	if m.flags&flagLogsDataResourceLogsDecoded == 0 {
		// Decode nested message(s).
		for i := range m.resourceLogs {
			m.resourceLogs[i].decode()
		}
		m.flags |= flagLogsDataResourceLogsDecoded
	}
	return m.resourceLogs
}

type ResourceLogs struct {
	ProtoMessage
	// The Resource
	resource *Resource
	// List of ScopeLogs
	scopeLogs []*ScopeLogs
}

func NewResourceLogs(bytes []byte) *ResourceLogs {
	m := &ResourceLogs{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *ResourceLogs) decode() {
	buf := codec.NewBuffer(m.bytes)

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
	m.scopeLogs = make([]*ScopeLogs, 0, scopeLogsCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.bytes)

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
				m.resource = &Resource{
					ProtoMessage: ProtoMessage{bytes: v, parent: &m.ProtoMessage},
				}
			case 2:
				// Decode scopeLogs.
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				m.scopeLogs[scopeLogsCount] = &ScopeLogs{
					ProtoMessage: ProtoMessage{bytes: v, parent: &m.ProtoMessage},
				}
				scopeLogsCount++
			}
			return true, nil
		},
	)
}

// Bitmasks that indicate that the particular nested message is decoded.
const flagResourceLogsResourceDecoded = 0x0000000000000002
const flagResourceLogsScopeLogsDecoded = 0x0000000000000004

func (m *ResourceLogs) GetResource() *Resource {
	if m.flags&flagResourceLogsResourceDecoded == 0 {
		// Decode nested message(s).
		m.resource.decode()
		m.flags |= flagResourceLogsResourceDecoded
	}
	return m.resource
}
func (m *ResourceLogs) GetScopeLogs() []*ScopeLogs {
	if m.flags&flagResourceLogsScopeLogsDecoded == 0 {
		// Decode nested message(s).
		for i := range m.scopeLogs {
			m.scopeLogs[i].decode()
		}
		m.flags |= flagResourceLogsScopeLogsDecoded
	}
	return m.scopeLogs
}

type Resource struct {
	ProtoMessage
	attributes             []*KeyValue
	droppedAttributesCount uint32
}

func NewResource(bytes []byte) *Resource {
	m := &Resource{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *Resource) decode() {
	buf := codec.NewBuffer(m.bytes)

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
	m.attributes = make([]*KeyValue, 0, attributesCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.bytes)

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
				m.attributes[attributesCount] = &KeyValue{
					ProtoMessage: ProtoMessage{bytes: v, parent: &m.ProtoMessage},
				}
				attributesCount++
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

// Bitmasks that indicate that the particular nested message is decoded.
const flagResourceAttributesDecoded = 0x0000000000000002

func (m *Resource) GetAttributes() []*KeyValue {
	if m.flags&flagResourceAttributesDecoded == 0 {
		// Decode nested message(s).
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.flags |= flagResourceAttributesDecoded
	}
	return m.attributes
}
func (m *Resource) GetDroppedAttributesCount() uint32 {
	return m.droppedAttributesCount
}

// A collection of Logs produced by a Scope.
type ScopeLogs struct {
	ProtoMessage
	// A list of log records.
	logRecords []*LogRecord
}

func NewScopeLogs(bytes []byte) *ScopeLogs {
	m := &ScopeLogs{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *ScopeLogs) decode() {
	buf := codec.NewBuffer(m.bytes)

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
	m.logRecords = make([]*LogRecord, 0, logRecordsCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.bytes)

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
				// The slice is pre-allocated, assign to the appropriate index.
				m.logRecords[logRecordsCount] = &LogRecord{
					ProtoMessage: ProtoMessage{bytes: v, parent: &m.ProtoMessage},
				}
				logRecordsCount++
			}
			return true, nil
		},
	)
}

// Bitmasks that indicate that the particular nested message is decoded.
const flagScopeLogsLogRecordsDecoded = 0x0000000000000002

func (m *ScopeLogs) GetLogRecords() []*LogRecord {
	if m.flags&flagScopeLogsLogRecordsDecoded == 0 {
		// Decode nested message(s).
		for i := range m.logRecords {
			m.logRecords[i].decode()
		}
		m.flags |= flagScopeLogsLogRecordsDecoded
	}
	return m.logRecords
}

type LogRecord struct {
	ProtoMessage
	timeUnixNano           uint64
	attributes             []*KeyValue
	droppedAttributesCount uint32
}

func NewLogRecord(bytes []byte) *LogRecord {
	m := &LogRecord{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *LogRecord) decode() {
	buf := codec.NewBuffer(m.bytes)

	// Count all repeated fields. We need one counter per field.
	attributesCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 2 {
				attributesCount++
			}
		},
	)

	// Pre-allocate slices for repeated fields.
	m.attributes = make([]*KeyValue, 0, attributesCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.bytes)

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
				m.attributes[attributesCount] = &KeyValue{
					ProtoMessage: ProtoMessage{bytes: v, parent: &m.ProtoMessage},
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

// Bitmasks that indicate that the particular nested message is decoded.
const flagLogRecordAttributesDecoded = 0x0000000000000002

func (m *LogRecord) GetTimeUnixNano() uint64 {
	return m.timeUnixNano
}
func (m *LogRecord) GetAttributes() []*KeyValue {
	if m.flags&flagLogRecordAttributesDecoded == 0 {
		// Decode nested message(s).
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.flags |= flagLogRecordAttributesDecoded
	}
	return m.attributes
}
func (m *LogRecord) GetDroppedAttributesCount() uint32 {
	return m.droppedAttributesCount
}

type KeyValue struct {
	ProtoMessage
	key   string
	value string
}

func NewKeyValue(bytes []byte) *KeyValue {
	m := &KeyValue{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *KeyValue) decode() {
	buf := codec.NewBuffer(m.bytes)

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

func (m *KeyValue) GetKey() string {
	return m.key
}
func (m *KeyValue) GetValue() string {
	return m.value
}
