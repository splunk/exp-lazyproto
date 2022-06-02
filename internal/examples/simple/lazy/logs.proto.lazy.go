package simple

type LogsData struct {
	resource_logs []*ResourceLogs
}
type ResourceLogs struct {
	resource   *Resource
	scope_logs []*ScopeLogs
}
type Resource struct {
	attributes               []*KeyValue
	dropped_attributes_count uint32
}
type KeyValue struct {
	key   string
	value string
}
type ScopeLogs struct {
	log_records []*LogRecord
}
type LogRecord struct {
	time_unix_nano           uint64
	attributes               []*KeyValue
	dropped_attributes_count uint32
}
