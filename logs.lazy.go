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

func (m *LogsData) Marshal(ps *molecule.ProtoStream) error {
	for _, logs := range m.resourceLogs {
		//ps.Embedded(
		//	1, func() error {
		//		return logs.Marshal(ps)
		//	},
		//)
		token := ps.BeginEmbedded()
		logs.Marshal(ps)
		ps.EndEmbedded(token, 1)
	}
	return nil
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

func (m *ResourceLogs) Marshal(ps *molecule.ProtoStream) error {
	if m.resource != nil {
		//ps.Embedded(
		//	1, func() error {
		//		return m.resource.Marshal(ps)
		//	},
		//)
		token := ps.BeginEmbedded()
		m.resource.Marshal(ps)
		ps.EndEmbedded(token, 1)
	}
	for _, logs := range m.scopeLogs {
		//ps.Embedded(
		//	2, func() error {
		//		return logs.Marshal(ps)
		//	},
		//)
		token := ps.BeginEmbedded()
		logs.Marshal(ps)
		ps.EndEmbedded(token, 2)
	}
	return nil
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

var resourceAttrKeyPrepared = molecule.PrepareEmbeddedField(1)
var resourceDroppedKeyPrepared = molecule.PrepareUint32Field(2)

func (m *Resource) Marshal(ps *molecule.ProtoStream) error {
	for _, attr := range m.attributes {
		//ps.Embedded(
		//	1, func() error {
		//		return attr.Marshal(ps)
		//	},
		//)
		token := ps.BeginEmbedded()
		attr.Marshal(ps)
		//ps.EndEmbedded(token, 1)
		ps.EndEmbeddedPrepared(token, resourceAttrKeyPrepared)
	}
	//ps.Uint32(2, m.DroppedAttributesCount)
	ps.Uint32Prepared(resourceDroppedKeyPrepared, m.DroppedAttributesCount)
	return nil
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

var keyValuePreparedKey = molecule.PrepareStringField(1)
var keyValuePreparedValue = molecule.PrepareStringField(2)

func (m *KeyValue) Marshal(ps *molecule.ProtoStream) error {
	ps.PreparedString(keyValuePreparedKey, m.Key)
	ps.PreparedString(keyValuePreparedValue, m.Value)
	//ps.String(1, m.Key)
	//ps.String(2, m.Value)
	return nil
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

func (m *ScopeLogs) Marshal(ps *molecule.ProtoStream) error {
	for _, logRecord := range m.logRecords {
		//ps.Embedded(
		//	1, func() error {
		//		return logRecord.Marshal(ps)
		//	},
		//)
		token := ps.BeginEmbedded()
		logRecord.Marshal(ps)
		ps.EndEmbedded(token, 1)
	}
	return nil
}

type LogRecord struct {
	bytes                  []byte
	timeUnixNano           uint64
	attributes             []KeyValue
	droppedAttributesCount uint32

	flags uint64
}

const logRecordAttributesDecoded = 1
const logRecordModified = 2

func (m *LogRecord) GetAttributes() *[]KeyValue {
	if m.flags&logRecordAttributesDecoded == 0 {
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.flags |= logRecordAttributesDecoded
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
				m.timeUnixNano = v

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
				m.droppedAttributesCount = v
			}
			return true, nil
		},
	)
}

var logRecordTimePrepared = molecule.PrepareFixed64Field(1)
var logRecordAttrPrepared = molecule.PrepareEmbeddedField(2)
var logRecordDroppedPrepared = molecule.PrepareUint32Field(3)

func (m *LogRecord) Marshal(ps *molecule.ProtoStream) error {
	//if m.flags&logRecordModified != 0 {
	ps.Fixed64Prepared(logRecordTimePrepared, m.timeUnixNano)

	for _, attr := range m.attributes {
		//ps.Embedded(
		//	2, func() error {
		//		return attr.Marshal(ps)
		//	},
		//)

		token := ps.BeginEmbedded()
		attr.Marshal(ps)
		ps.EndEmbeddedPrepared(token, logRecordAttrPrepared)
	}

	ps.Uint32Prepared(logRecordDroppedPrepared, m.droppedAttributesCount)
	//} else {
	//	return ps.Raw(m.bytes)
	//}
	return nil
}
