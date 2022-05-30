package simple

import (
	"github.com/richardartoul/molecule"
	"github.com/richardartoul/molecule/src/codec"
)

type LogsData struct {
	bytes        []byte
	resourceLogs []ResourceLogs

	fieldsDecoded uint64
}

const logsDataResourceLogsDecoded = 1

func NewLogsData(bytes []byte) *LogsData {
	m := &LogsData{bytes: bytes}
	m.decode()
	return m
}

func (m *LogsData) decode() {
	buf := codec.NewBuffer(m.bytes)
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.resourceLogs = append(m.resourceLogs, ResourceLogs{bytes: v})
			}
			return true, nil
		},
	)
}

func (m *LogsData) GetResourceLogs() *[]ResourceLogs {
	if m.fieldsDecoded&logsDataResourceLogsDecoded == 0 {
		for i := range m.resourceLogs {
			m.resourceLogs[i].decode()
		}
		m.fieldsDecoded |= logsDataResourceLogsDecoded
	}
	return &m.resourceLogs
}

type ResourceLogs struct {
	bytes     []byte
	resource  *Resource
	scopeLogs []ScopeLogs

	fieldsDecoded uint64
}

const resourceLogsResourceDecoded = 1
const resourceLogsScopeLogsDecoded = 2

func (m *ResourceLogs) GetResource() **Resource {
	if m.fieldsDecoded&resourceLogsResourceDecoded == 0 {
		m.resource.decode()
		m.fieldsDecoded |= resourceLogsResourceDecoded
	}
	return &m.resource
}

func (m *ResourceLogs) GetScopeLogs() *[]ScopeLogs {
	if m.fieldsDecoded&resourceLogsScopeLogsDecoded == 0 {
		for i := range m.scopeLogs {
			m.scopeLogs[i].decode()
		}
		m.fieldsDecoded |= resourceLogsScopeLogsDecoded
	}
	return &m.scopeLogs
}

func (m *ResourceLogs) decode() {
	buf := codec.NewBuffer(m.bytes)
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.resource = &Resource{bytes: v}

			case 2:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.scopeLogs = append(m.scopeLogs, ScopeLogs{bytes: v})
			}
			return true, nil
		},
	)
}

type Resource struct {
	bytes                  []byte
	attributes             []KeyValue
	DroppedAttributesCount uint32

	fieldsDecoded uint64
}

const resourceAttributesDecoded = 1

func (m *Resource) GetAttributes() *[]KeyValue {
	if m.fieldsDecoded&resourceAttributesDecoded == 0 {
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.fieldsDecoded |= resourceAttributesDecoded
	}
	return &m.attributes
}

func (m *Resource) decode() {
	buf := codec.NewBuffer(m.bytes)
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.attributes = append(m.attributes, KeyValue{bytes: v})

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

type KeyValue struct {
	bytes []byte
	Key   string
	Value string

	fieldsDecoded uint64
}

func (m *KeyValue) decode() {
	buf := codec.NewBuffer(m.bytes)
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.Key = v

			case 2:
				v, err := value.AsStringUnsafe()
				if err != nil {
					return false, err
				}
				m.Value = v
			}
			return true, nil
		},
	)
}

type ScopeLogs struct {
	bytes      []byte
	logRecords []LogRecord

	fieldsDecoded uint64
}

const scopeLogsLogRecordsDecoded = 1

func (m *ScopeLogs) GetLogRecords() *[]LogRecord {
	if m.fieldsDecoded&scopeLogsLogRecordsDecoded == 0 {
		for i := range m.logRecords {
			m.logRecords[i].decode()
		}
		m.fieldsDecoded |= scopeLogsLogRecordsDecoded
	}
	return &m.logRecords
}

func (m *ScopeLogs) decode() {
	buf := codec.NewBuffer(m.bytes)
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.logRecords = append(m.logRecords, LogRecord{bytes: v})
			}
			return true, nil
		},
	)
}

type LogRecord struct {
	bytes                  []byte
	TimeUnixNano           uint64
	attributes             []KeyValue
	DroppedAttributesCount uint32

	fieldsDecoded uint64
}

const logRecordAttributesDecoded = 1

func (m *LogRecord) GetAttributes() *[]KeyValue {
	if m.fieldsDecoded&logRecordAttributesDecoded == 0 {
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.fieldsDecoded |= logRecordAttributesDecoded
	}
	return &m.attributes
}

func (m *LogRecord) decode() {
	buf := codec.NewBuffer(m.bytes)
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				v, err := value.AsFixed64()
				if err != nil {
					return false, err
				}
				m.TimeUnixNano = v

			case 2:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.attributes = append(m.attributes, KeyValue{bytes: v})

			case 3:
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
