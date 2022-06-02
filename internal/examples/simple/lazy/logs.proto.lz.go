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
	resource_logs []*ResourceLogs
}

func NewLogsData(bytes []byte) *LogsData {
	m := &LogsData{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *LogsData) decode() {
	buf := codec.NewBuffer(m.bytes)

	// Count all repeated fields. We need one counter per field.
	resource_logsCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 1 {
				resource_logsCount++
			}
		},
	)

	// Pre-allocate slices for repeated fields.
	m.resource_logs = make([]*ResourceLogs, 0, resource_logsCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.bytes)

	// Set slice indexes to 0 to begin iterating over repeated fields.
	resource_logsCount = 0
	// Iterate and decode the fields.
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode resource_logs
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				m.resource_logs[resource_logsCount] = &ResourceLogs{
					ProtoMessage: ProtoMessage{bytes: v, parent: &m.ProtoMessage},
				}
				resource_logsCount++
			}
			return true, nil
		},
	)
}

// Bitmasks that indicate that the particular nested message is decoded
const flagLogsDataresource_logsDecoded = 0x0000000000000002

func (m *LogsData) Getresource_logs() []*ResourceLogs {
	if m.flags&flagLogsDataresource_logsDecoded == 0 {
		// Decode nested message(s)
		for i := range m.resource_logs {
			m.resource_logs[i].decode()
		}
		m.flags |= flagLogsDataresource_logsDecoded
	}
	return m.resource_logs
}

type ResourceLogs struct {
	ProtoMessage
	// The Resource
	resource *Resource
	// List of ScopeLogs
	scope_logs []*ScopeLogs
}

func NewResourceLogs(bytes []byte) *ResourceLogs {
	m := &ResourceLogs{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *ResourceLogs) decode() {
	buf := codec.NewBuffer(m.bytes)

	// Count all repeated fields. We need one counter per field.
	scope_logsCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 2 {
				scope_logsCount++
			}
		},
	)

	// Pre-allocate slices for repeated fields.
	m.scope_logs = make([]*ScopeLogs, 0, scope_logsCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.bytes)

	// Set slice indexes to 0 to begin iterating over repeated fields.
	scope_logsCount = 0
	// Iterate and decode the fields.
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode resource
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.resource = &Resource{
					ProtoMessage: ProtoMessage{bytes: v, parent: &m.ProtoMessage},
				}
			case 2:
				// Decode scope_logs
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				m.scope_logs[scope_logsCount] = &ScopeLogs{
					ProtoMessage: ProtoMessage{bytes: v, parent: &m.ProtoMessage},
				}
				scope_logsCount++
			}
			return true, nil
		},
	)
}

// Bitmasks that indicate that the particular nested message is decoded
const flagResourceLogsresourceDecoded = 0x0000000000000002
const flagResourceLogsscope_logsDecoded = 0x0000000000000004

func (m *ResourceLogs) Getresource() *Resource {
	if m.flags&flagResourceLogsresourceDecoded == 0 {
		// Decode nested message(s)
		m.resource.decode()
		m.flags |= flagResourceLogsresourceDecoded
	}
	return m.resource
}
func (m *ResourceLogs) Getscope_logs() []*ScopeLogs {
	if m.flags&flagResourceLogsscope_logsDecoded == 0 {
		// Decode nested message(s)
		for i := range m.scope_logs {
			m.scope_logs[i].decode()
		}
		m.flags |= flagResourceLogsscope_logsDecoded
	}
	return m.scope_logs
}

type Resource struct {
	ProtoMessage
	attributes               []*KeyValue
	dropped_attributes_count uint32
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
				// Decode attributes
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
				// Decode dropped_attributes_count
				v, err := value.AsUint32()
				if err != nil {
					return false, err
				}
				m.dropped_attributes_count = v
			}
			return true, nil
		},
	)
}

// Bitmasks that indicate that the particular nested message is decoded
const flagResourceattributesDecoded = 0x0000000000000002

func (m *Resource) Getattributes() []*KeyValue {
	if m.flags&flagResourceattributesDecoded == 0 {
		// Decode nested message(s)
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.flags |= flagResourceattributesDecoded
	}
	return m.attributes
}
func (m *Resource) Getdropped_attributes_count() uint32 {
	return m.dropped_attributes_count
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
				// Decode key
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.key = v
			case 2:
				// Decode value
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

func (m *KeyValue) Getkey() string {
	return m.key
}
func (m *KeyValue) Getvalue() string {
	return m.value
}

// A collection of Logs produced by a Scope.
type ScopeLogs struct {
	ProtoMessage
	// A list of log records.
	log_records []*LogRecord
}

func NewScopeLogs(bytes []byte) *ScopeLogs {
	m := &ScopeLogs{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *ScopeLogs) decode() {
	buf := codec.NewBuffer(m.bytes)

	// Count all repeated fields. We need one counter per field.
	log_recordsCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 1 {
				log_recordsCount++
			}
		},
	)

	// Pre-allocate slices for repeated fields.
	m.log_records = make([]*LogRecord, 0, log_recordsCount)

	// Reset the buffer to start iterating over the fields again
	buf.Reset(m.bytes)

	// Set slice indexes to 0 to begin iterating over repeated fields.
	log_recordsCount = 0
	// Iterate and decode the fields.
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				// Decode log_records
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				// The slice is pre-allocated, assign to the appropriate index.
				m.log_records[log_recordsCount] = &LogRecord{
					ProtoMessage: ProtoMessage{bytes: v, parent: &m.ProtoMessage},
				}
				log_recordsCount++
			}
			return true, nil
		},
	)
}

// Bitmasks that indicate that the particular nested message is decoded
const flagScopeLogslog_recordsDecoded = 0x0000000000000002

func (m *ScopeLogs) Getlog_records() []*LogRecord {
	if m.flags&flagScopeLogslog_recordsDecoded == 0 {
		// Decode nested message(s)
		for i := range m.log_records {
			m.log_records[i].decode()
		}
		m.flags |= flagScopeLogslog_recordsDecoded
	}
	return m.log_records
}

type LogRecord struct {
	ProtoMessage
	time_unix_nano           uint64
	attributes               []*KeyValue
	dropped_attributes_count uint32
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
				// Decode time_unix_nano
				v, err := value.AsFixed64()
				if err != nil {
					return false, err
				}
				m.time_unix_nano = v
			case 2:
				// Decode attributes
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
				// Decode dropped_attributes_count
				v, err := value.AsUint32()
				if err != nil {
					return false, err
				}
				m.dropped_attributes_count = v
			}
			return true, nil
		},
	)
}

// Bitmasks that indicate that the particular nested message is decoded
const flagLogRecordattributesDecoded = 0x0000000000000002

func (m *LogRecord) Gettime_unix_nano() uint64 {
	return m.time_unix_nano
}
func (m *LogRecord) Getattributes() []*KeyValue {
	if m.flags&flagLogRecordattributesDecoded == 0 {
		// Decode nested message(s)
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.flags |= flagLogRecordattributesDecoded
	}
	return m.attributes
}
func (m *LogRecord) Getdropped_attributes_count() uint32 {
	return m.dropped_attributes_count
}
