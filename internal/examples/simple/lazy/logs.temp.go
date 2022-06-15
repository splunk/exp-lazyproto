package simple

import "github.com/tigrannajaryan/exp-lazyproto/internal/protomessage"

func (m *LogsData) ResourceLogsSlice() (r Slice[ResourceLogs]) {
	if m._flags&flags_LogsData_ResourceLogs_Decoded == 0 {
		m.decodeResourceLogs()
	}
	return Slice[ResourceLogs]{
		slice: &m.resourceLogs,
		owner: &m._protoMessage,
	}
}

func (m *LogRecord) AttributesSlice() (r Slice[KeyValue]) {
	if m._flags&flags_LogRecord_Attributes_Decoded == 0 {
		m.decodeAttributes()
	}
	return Slice[KeyValue]{
		slice: &m.attributes,
		owner: &m._protoMessage,
	}
}

type Slice[T any] struct {
	slice *[]*T
	owner *protomessage.ProtoMessage
}

func (s *Slice[T]) Len() int {
	return len(*s.slice)
}

func (s *Slice[T]) At(idx int) *T {
	return (*s.slice)[idx]
}

func (s *Slice[T]) RemoveIf(func(e *T) bool) {

}

func (s *Slice[T]) Range(f func(e *T) bool) {
	for _, e := range *s.slice {
		if !f(e) {
			break
		}
	}
}
