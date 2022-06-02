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
type LogsData struct {
	ProtoMessage
	resource_logs []*ResourceLogs
}

func NewLogsData(bytes []byte) *LogsData {
	m := &LogsData{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *LogsData) decode() {
	buf := codec.NewBuffer(m.bytes)

	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
			}
			return true, nil
		},
	)
}

type ResourceLogs struct {
	ProtoMessage
	resource   *Resource
	scope_logs []*ScopeLogs
}

func NewResourceLogs(bytes []byte) *ResourceLogs {
	m := &ResourceLogs{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *ResourceLogs) decode() {
	buf := codec.NewBuffer(m.bytes)

	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
			case 2:
			}
			return true, nil
		},
	)
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

	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
			case 2:
			}
			return true, nil
		},
	)
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

type ScopeLogs struct {
	ProtoMessage
	log_records []*LogRecord
}

func NewScopeLogs(bytes []byte) *ScopeLogs {
	m := &ScopeLogs{ProtoMessage: ProtoMessage{bytes: bytes}}
	m.decode()
	return m
}

func (m *ScopeLogs) decode() {
	buf := codec.NewBuffer(m.bytes)

	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 1:
			}
			return true, nil
		},
	)
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

	molecule.MessageEach(
		buf, func(fieldNum int32, value molecule.Value) (bool, error) {
			switch fieldNum {
			case 3:
			case 1:
			case 2:
			}
			return true, nil
		},
	)
}
