package simple

import (
	"github.com/richardartoul/molecule"
	"github.com/richardartoul/molecule/src/codec"
)

const flagsMessageModified = 1

type ProtoMessage struct {
	flags  uint64
	parent *ProtoMessage
}

type LogsData struct {
	ProtoMessage
	bytes        []byte
	resourceLogs []ResourceLogs
}

const logsDataResourceLogsDecoded = 2

func NewLogsData(bytes []byte) *LogsData {
	m := &LogsData{bytes: bytes}
	m.decode()
	return m
}

func (m *LogsData) decode() {
	buf := codec.NewBuffer(m.bytes)

	lrCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 1 {
				lrCount++
			}
		},
	)
	m.resourceLogs = make([]ResourceLogs, 0, lrCount)

	buf.Reset(m.bytes)
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.resourceLogs = append(
					m.resourceLogs,
					ResourceLogs{
						bytes:        v,
						ProtoMessage: ProtoMessage{parent: &m.ProtoMessage},
					},
				)
			}
			return true, nil
		},
	)
}

func (m *LogsData) GetResourceLogs() *[]ResourceLogs {
	if m.flags&logsDataResourceLogsDecoded == 0 {
		for i := range m.resourceLogs {
			m.resourceLogs[i].decode()
		}
		m.flags |= logsDataResourceLogsDecoded
	}
	return &m.resourceLogs
}

func (m *LogsData) Marshal(ps *molecule.ProtoStream) error {
	if m.flags&flagsMessageModified != 0 {
		for _, logs := range m.resourceLogs {
			token := ps.BeginEmbedded()
			logs.Marshal(ps)
			ps.EndEmbedded(token, 1)
		}
	} else {
		ps.Raw(m.bytes)
	}
	return nil
}

type ResourceLogs struct {
	ProtoMessage
	bytes     []byte
	resource  *Resource
	scopeLogs []ScopeLogs
}

const resourceLogsResourceDecoded = 2
const resourceLogsScopeLogsDecoded = 4

func (m *ResourceLogs) GetResource() **Resource {
	if m.flags&resourceLogsResourceDecoded == 0 {
		m.resource.decode()
		m.flags |= resourceLogsResourceDecoded
	}
	return &m.resource
}

func (m *ResourceLogs) GetScopeLogs() *[]ScopeLogs {
	if m.flags&resourceLogsScopeLogsDecoded == 0 {
		for i := range m.scopeLogs {
			m.scopeLogs[i].decode()
		}
		m.flags |= resourceLogsScopeLogsDecoded
	}
	return &m.scopeLogs
}

func (m *ResourceLogs) decode() {
	buf := codec.NewBuffer(m.bytes)

	lrCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 2 {
				lrCount++
			}
		},
	)
	//m.scopeLogs = make([]ScopeLogs, 0, lrCount)
	m.scopeLogs = scopeLogsPool.GetScopeLogss(lrCount)

	lrIndex := 0
	buf.Reset(m.bytes)
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				m.resource = &Resource{
					bytes:        v,
					ProtoMessage: ProtoMessage{parent: &m.ProtoMessage},
				}

			case 2:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				sl := &m.scopeLogs[lrIndex]
				lrIndex++
				*sl = ScopeLogs{
					bytes:        v,
					ProtoMessage: ProtoMessage{parent: &m.ProtoMessage},
				}
			}
			return true, nil
		},
	)
}

func (m *ResourceLogs) Marshal(ps *molecule.ProtoStream) error {
	if m.flags&flagsMessageModified != 0 {
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
		ps.Raw(m.bytes)
	}
	return nil
}

type Resource struct {
	ProtoMessage
	bytes                  []byte
	attributes             []KeyValue
	DroppedAttributesCount uint32
}

const resourceAttributesDecoded = 2

func (m *Resource) GetAttributes() *[]KeyValue {
	if m.flags&resourceAttributesDecoded == 0 {
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.flags |= resourceAttributesDecoded
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
				m.attributes = append(
					m.attributes, KeyValue{
						bytes:        v,
						ProtoMessage: ProtoMessage{parent: &m.ProtoMessage},
					},
				)

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
	if m.flags&flagsMessageModified != 0 {
		for _, attr := range m.attributes {
			token := ps.BeginEmbedded()
			attr.Marshal(ps)
			//ps.EndEmbedded(token, 1)
			ps.EndEmbeddedPrepared(token, resourceAttrKeyPrepared)
		}
		ps.Uint32Prepared(resourceDroppedKeyPrepared, m.DroppedAttributesCount)
	} else {
		ps.Raw(m.bytes)
	}
	return nil
}

type ScopeLogs struct {
	ProtoMessage
	bytes      []byte
	logRecords []LogRecord
}

const scopeLogsLogRecordsDecoded = 2

func (m *ScopeLogs) GetLogRecords() *[]LogRecord {
	if m.flags&scopeLogsLogRecordsDecoded == 0 {
		for i := range m.logRecords {
			m.logRecords[i].decode()
		}
		m.flags |= scopeLogsLogRecordsDecoded
	}
	return &m.logRecords
}

func (m *ScopeLogs) decode() {
	buf := codec.NewBuffer(m.bytes)

	lrCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 1 {
				lrCount++
			}
		},
	)
	//m.logRecords = make([]LogRecord, 0, lrCount)
	m.logRecords = logRecordPool.GetLogRecords(lrCount)

	lrIndex := 0
	buf.Reset(m.bytes)
	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				lr := &m.logRecords[lrIndex]
				*lr = LogRecord{
					bytes:        v,
					ProtoMessage: ProtoMessage{parent: &m.ProtoMessage},
				}
			}
			return true, nil
		},
	)
}

func (m *ScopeLogs) Marshal(ps *molecule.ProtoStream) error {
	if m.flags&flagsMessageModified != 0 {
		for _, logRecord := range m.logRecords {
			token := ps.BeginEmbedded()
			logRecord.Marshal(ps)
			ps.EndEmbedded(token, 1)
		}
	} else {
		ps.Raw(m.bytes)
	}
	return nil
}

type LogRecord struct {
	ProtoMessage
	bytes                  []byte
	timeUnixNano           uint64
	attributes             []KeyValue
	droppedAttributesCount uint32
}

const logRecordAttributesDecoded = 2

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
	attrCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 2 {
				attrCount++
			}
		},
	)
	//m.attributes = make([]*KeyValue, 0, attrCount)
	m.attributes = keyValuePool.GetKeyValues(attrCount)

	attrIndex := 0
	buf.Reset(m.bytes)
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
				kv := &m.attributes[attrIndex]
				attrIndex++
				*kv = KeyValue{
					bytes:        v,
					ProtoMessage: ProtoMessage{parent: &m.ProtoMessage},
				}

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
	if m.flags&flagsMessageModified != 0 {
		ps.Fixed64Prepared(logRecordTimePrepared, m.timeUnixNano)

		for _, attr := range m.attributes {
			token := ps.BeginEmbedded()
			attr.Marshal(ps)
			ps.EndEmbeddedPrepared(token, logRecordAttrPrepared)
		}

		ps.Uint32Prepared(logRecordDroppedPrepared, m.droppedAttributesCount)
	} else {
		ps.Raw(m.bytes)
	}
	return nil
}

type KeyValue struct {
	ProtoMessage
	bytes []byte
	key   string
	value string
}

func (m *KeyValue) Key() string {
	return m.key
}

func (m *KeyValue) SetKey(s string) {
	m.key = s
	m.markModified() // TODO: check if this is inlined.
}

func (m *ProtoMessage) markModified() {
	if m.flags&flagsMessageModified != 0 {
		return
	}

	m.flags |= flagsMessageModified
	parent := m.parent
	for parent != nil {
		if parent.flags&flagsMessageModified != 0 {
			break
		}
		parent.flags |= flagsMessageModified
		parent = parent.parent
	}
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
				m.key = v

			case 2:
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

var keyValuePreparedKey = molecule.PrepareStringField(1)
var keyValuePreparedValue = molecule.PrepareStringField(2)

func (m *KeyValue) Marshal(ps *molecule.ProtoStream) error {
	if m.flags&flagsMessageModified != 0 {
		ps.PreparedString(keyValuePreparedKey, m.key)
		ps.PreparedString(keyValuePreparedValue, m.value)
	} else {
		ps.Raw(m.bytes)
	}
	return nil
}
