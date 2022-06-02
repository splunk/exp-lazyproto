package simple

import (
	lazyproto "github.com/tigrannajaryan/exp-lazyproto"
	"github.com/tigrannajaryan/molecule"
	"github.com/tigrannajaryan/molecule/src/codec"
)

type LogsData struct {
	protoMessage lazyproto.ProtoMessage
	resourceLogs []ResourceLogs
}

const logsDataResourceLogsDecoded = 2

func NewLogsData(bytes []byte) *LogsData {
	m := &LogsData{protoMessage: lazyproto.ProtoMessage{Bytes: bytes}}
	m.decode()
	return m
}

func (m *LogsData) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)

	lrCount := 0
	molecule.MessageFieldNums(
		buf, func(fieldNum int32) {
			if fieldNum == 1 {
				lrCount++
			}
		},
	)
	m.resourceLogs = make([]ResourceLogs, 0, lrCount)

	buf.Reset(m.protoMessage.Bytes)
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
						protoMessage: lazyproto.ProtoMessage{
							Parent: &m.protoMessage, Bytes: v,
						},
					},
				)
			}
			return true, nil
		},
	)
}

func (m *LogsData) GetResourceLogs() *[]ResourceLogs {
	if m.protoMessage.Flags&logsDataResourceLogsDecoded == 0 {
		for i := range m.resourceLogs {
			m.resourceLogs[i].decode()
		}
		m.protoMessage.Flags |= logsDataResourceLogsDecoded
	}
	return &m.resourceLogs
}

func (m *LogsData) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		for _, logs := range m.resourceLogs {
			token := ps.BeginEmbedded()
			logs.Marshal(ps)
			ps.EndEmbedded(token, 1)
		}
	} else {
		ps.Raw(m.protoMessage.Bytes)
	}
	return nil
}

type ResourceLogs struct {
	protoMessage lazyproto.ProtoMessage
	resource     *Resource
	scopeLogs    []ScopeLogs
}

const resourceLogsResourceDecoded = 2
const resourceLogsScopeLogsDecoded = 4

func (m *ResourceLogs) GetResource() **Resource {
	if m.protoMessage.Flags&resourceLogsResourceDecoded == 0 {
		m.resource.decode()
		m.protoMessage.Flags |= resourceLogsResourceDecoded
	}
	return &m.resource
}

func (m *ResourceLogs) GetScopeLogs() *[]ScopeLogs {
	if m.protoMessage.Flags&resourceLogsScopeLogsDecoded == 0 {
		for i := range m.scopeLogs {
			m.scopeLogs[i].decode()
		}
		m.protoMessage.Flags |= resourceLogsScopeLogsDecoded
	}
	return &m.scopeLogs
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
	m.scopeLogs = scopeLogsPool.GetScopeLogss(lrCount)

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
				m.resource = &Resource{
					protoMessage: lazyproto.ProtoMessage{
						Parent: &m.protoMessage, Bytes: v,
					},
				}

			case 2:
				v, err := value.AsBytesUnsafe()
				if err != nil {
					return false, err
				}
				sl := &m.scopeLogs[lrIndex]
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
	attributes             []KeyValue
	DroppedAttributesCount uint32
}

const resourceAttributesDecoded = 2

func (m *Resource) GetAttributes() *[]KeyValue {
	if m.protoMessage.Flags&resourceAttributesDecoded == 0 {
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.protoMessage.Flags |= resourceAttributesDecoded
	}
	return &m.attributes
}

func (m *Resource) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)
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
						protoMessage: lazyproto.ProtoMessage{
							Parent: &m.protoMessage, Bytes: v,
						},
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
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		for _, attr := range m.attributes {
			token := ps.BeginEmbedded()
			attr.Marshal(ps)
			//ps.EndEmbedded(token, 1)
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
	logRecords   []LogRecord
}

const scopeLogsLogRecordsDecoded = 2

func (m *ScopeLogs) GetLogRecords() *[]LogRecord {
	if m.protoMessage.Flags&scopeLogsLogRecordsDecoded == 0 {
		for i := range m.logRecords {
			m.logRecords[i].decode()
		}
		m.protoMessage.Flags |= scopeLogsLogRecordsDecoded
	}
	return &m.logRecords
}

func (m *ScopeLogs) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)

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
	buf.Reset(m.protoMessage.Bytes)
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
					protoMessage: lazyproto.ProtoMessage{
						Parent: &m.protoMessage, Bytes: v,
					},
				}
			}
			return true, nil
		},
	)
}

func (m *ScopeLogs) Marshal(ps *molecule.ProtoStream) error {
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		for _, logRecord := range m.logRecords {
			token := ps.BeginEmbedded()
			logRecord.Marshal(ps)
			ps.EndEmbedded(token, 1)
		}
	} else {
		ps.Raw(m.protoMessage.Bytes)
	}
	return nil
}

type LogRecord struct {
	protoMessage           lazyproto.ProtoMessage
	timeUnixNano           uint64
	attributes             []*KeyValue
	droppedAttributesCount uint32
}

const logRecordAttributesDecoded = 2

func (m *LogRecord) GetAttributes() *[]*KeyValue {
	if m.protoMessage.Flags&logRecordAttributesDecoded == 0 {
		for i := range m.attributes {
			m.attributes[i].decode()
		}
		m.protoMessage.Flags |= logRecordAttributesDecoded
	}
	return &m.attributes
}

func (m *LogRecord) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)
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
	buf.Reset(m.protoMessage.Bytes)
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
				kv := m.attributes[attrIndex]
				attrIndex++
				*kv = KeyValue{
					protoMessage: lazyproto.ProtoMessage{
						Parent: &m.protoMessage, Bytes: v,
					},
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
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		ps.Fixed64Prepared(logRecordTimePrepared, m.timeUnixNano)

		for _, attr := range m.attributes {
			token := ps.BeginEmbedded()
			attr.Marshal(ps)
			ps.EndEmbeddedPrepared(token, logRecordAttrPrepared)
		}

		ps.Uint32Prepared(logRecordDroppedPrepared, m.droppedAttributesCount)
	} else {
		ps.Raw(m.protoMessage.Bytes)
	}
	return nil
}

type KeyValue struct {
	protoMessage lazyproto.ProtoMessage
	key          string
	value        string
}

func (m *KeyValue) Key() string {
	return m.key
}

func (m *KeyValue) SetKey(s string) {
	m.key = s
	m.protoMessage.MarkModified() // TODO: check if this is inlined.
}

func (m *KeyValue) decode() {
	buf := codec.NewBuffer(m.protoMessage.Bytes)
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
	if m.protoMessage.Flags&lazyproto.FlagsMessageModified != 0 {
		ps.PreparedString(keyValuePreparedKey, m.key)
		ps.PreparedString(keyValuePreparedValue, m.value)
	} else {
		ps.Raw(m.protoMessage.Bytes)
	}
	return nil
}
