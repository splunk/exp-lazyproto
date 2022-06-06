// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: logs.proto

package logs

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// SeverityNumber values
type SeverityNumber int32

const (
	// SeverityNumber is not specified
	SeverityNumber_SEVERITY_NUMBER_UNSPECIFIED SeverityNumber = 0
	SeverityNumber_SEVERITY_NUMBER_TRACE       SeverityNumber = 1
	SeverityNumber_SEVERITY_NUMBER_TRACE2      SeverityNumber = 2
	SeverityNumber_SEVERITY_NUMBER_TRACE3      SeverityNumber = 3
	SeverityNumber_SEVERITY_NUMBER_TRACE4      SeverityNumber = 4
	SeverityNumber_SEVERITY_NUMBER_DEBUG       SeverityNumber = 5
	SeverityNumber_SEVERITY_NUMBER_DEBUG2      SeverityNumber = 6
	SeverityNumber_SEVERITY_NUMBER_DEBUG3      SeverityNumber = 7
	SeverityNumber_SEVERITY_NUMBER_DEBUG4      SeverityNumber = 8
	SeverityNumber_SEVERITY_NUMBER_INFO        SeverityNumber = 9
	SeverityNumber_SEVERITY_NUMBER_INFO2       SeverityNumber = 10
	SeverityNumber_SEVERITY_NUMBER_INFO3       SeverityNumber = 11
	SeverityNumber_SEVERITY_NUMBER_INFO4       SeverityNumber = 12
	SeverityNumber_SEVERITY_NUMBER_WARN        SeverityNumber = 13
	SeverityNumber_SEVERITY_NUMBER_WARN2       SeverityNumber = 14
	SeverityNumber_SEVERITY_NUMBER_WARN3       SeverityNumber = 15
	SeverityNumber_SEVERITY_NUMBER_WARN4       SeverityNumber = 16
	SeverityNumber_SEVERITY_NUMBER_ERROR       SeverityNumber = 17
	SeverityNumber_SEVERITY_NUMBER_ERROR2      SeverityNumber = 18
	SeverityNumber_SEVERITY_NUMBER_ERROR3      SeverityNumber = 19
	SeverityNumber_SEVERITY_NUMBER_ERROR4      SeverityNumber = 20
	SeverityNumber_SEVERITY_NUMBER_FATAL       SeverityNumber = 21
	SeverityNumber_SEVERITY_NUMBER_FATAL2      SeverityNumber = 22
	SeverityNumber_SEVERITY_NUMBER_FATAL3      SeverityNumber = 23
	SeverityNumber_SEVERITY_NUMBER_FATAL4      SeverityNumber = 24
)

// Enum value maps for SeverityNumber.
var (
	SeverityNumber_name = map[int32]string{
		0:  "SEVERITY_NUMBER_UNSPECIFIED",
		1:  "SEVERITY_NUMBER_TRACE",
		2:  "SEVERITY_NUMBER_TRACE2",
		3:  "SEVERITY_NUMBER_TRACE3",
		4:  "SEVERITY_NUMBER_TRACE4",
		5:  "SEVERITY_NUMBER_DEBUG",
		6:  "SEVERITY_NUMBER_DEBUG2",
		7:  "SEVERITY_NUMBER_DEBUG3",
		8:  "SEVERITY_NUMBER_DEBUG4",
		9:  "SEVERITY_NUMBER_INFO",
		10: "SEVERITY_NUMBER_INFO2",
		11: "SEVERITY_NUMBER_INFO3",
		12: "SEVERITY_NUMBER_INFO4",
		13: "SEVERITY_NUMBER_WARN",
		14: "SEVERITY_NUMBER_WARN2",
		15: "SEVERITY_NUMBER_WARN3",
		16: "SEVERITY_NUMBER_WARN4",
		17: "SEVERITY_NUMBER_ERROR",
		18: "SEVERITY_NUMBER_ERROR2",
		19: "SEVERITY_NUMBER_ERROR3",
		20: "SEVERITY_NUMBER_ERROR4",
		21: "SEVERITY_NUMBER_FATAL",
		22: "SEVERITY_NUMBER_FATAL2",
		23: "SEVERITY_NUMBER_FATAL3",
		24: "SEVERITY_NUMBER_FATAL4",
	}
	SeverityNumber_value = map[string]int32{
		"SEVERITY_NUMBER_UNSPECIFIED": 0,
		"SEVERITY_NUMBER_TRACE":       1,
		"SEVERITY_NUMBER_TRACE2":      2,
		"SEVERITY_NUMBER_TRACE3":      3,
		"SEVERITY_NUMBER_TRACE4":      4,
		"SEVERITY_NUMBER_DEBUG":       5,
		"SEVERITY_NUMBER_DEBUG2":      6,
		"SEVERITY_NUMBER_DEBUG3":      7,
		"SEVERITY_NUMBER_DEBUG4":      8,
		"SEVERITY_NUMBER_INFO":        9,
		"SEVERITY_NUMBER_INFO2":       10,
		"SEVERITY_NUMBER_INFO3":       11,
		"SEVERITY_NUMBER_INFO4":       12,
		"SEVERITY_NUMBER_WARN":        13,
		"SEVERITY_NUMBER_WARN2":       14,
		"SEVERITY_NUMBER_WARN3":       15,
		"SEVERITY_NUMBER_WARN4":       16,
		"SEVERITY_NUMBER_ERROR":       17,
		"SEVERITY_NUMBER_ERROR2":      18,
		"SEVERITY_NUMBER_ERROR3":      19,
		"SEVERITY_NUMBER_ERROR4":      20,
		"SEVERITY_NUMBER_FATAL":       21,
		"SEVERITY_NUMBER_FATAL2":      22,
		"SEVERITY_NUMBER_FATAL3":      23,
		"SEVERITY_NUMBER_FATAL4":      24,
	}
)

func (x SeverityNumber) Enum() *SeverityNumber {
	p := new(SeverityNumber)
	*p = x
	return p
}

func (x SeverityNumber) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SeverityNumber) Descriptor() protoreflect.EnumDescriptor {
	return file_logs_proto_enumTypes[0].Descriptor()
}

func (SeverityNumber) Type() protoreflect.EnumType {
	return &file_logs_proto_enumTypes[0]
}

func (x SeverityNumber) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use SeverityNumber.Descriptor instead.
func (SeverityNumber) EnumDescriptor() ([]byte, []int) {
	return file_logs_proto_rawDescGZIP(), []int{0}
}

// LogsData contains all log data
type LogsData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of ResourceLogs
	ResourceLogs []*ResourceLogs `protobuf:"bytes,1,rep,name=resource_logs,json=resourceLogs,proto3" json:"resource_logs,omitempty"`
}

func (x *LogsData) Reset() {
	*x = LogsData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logs_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogsData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogsData) ProtoMessage() {}

func (x *LogsData) ProtoReflect() protoreflect.Message {
	mi := &file_logs_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogsData.ProtoReflect.Descriptor instead.
func (*LogsData) Descriptor() ([]byte, []int) {
	return file_logs_proto_rawDescGZIP(), []int{0}
}

func (x *LogsData) GetResourceLogs() []*ResourceLogs {
	if x != nil {
		return x.ResourceLogs
	}
	return nil
}

type ResourceLogs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The Resource
	Resource *Resource `protobuf:"bytes,1,opt,name=resource,proto3" json:"resource,omitempty"`
	// List of ScopeLogs
	ScopeLogs []*ScopeLogs `protobuf:"bytes,2,rep,name=scope_logs,json=scopeLogs,proto3" json:"scope_logs,omitempty"`
	SchemaUrl string       `protobuf:"bytes,3,opt,name=schema_url,json=schemaUrl,proto3" json:"schema_url,omitempty"`
}

func (x *ResourceLogs) Reset() {
	*x = ResourceLogs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logs_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResourceLogs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResourceLogs) ProtoMessage() {}

func (x *ResourceLogs) ProtoReflect() protoreflect.Message {
	mi := &file_logs_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResourceLogs.ProtoReflect.Descriptor instead.
func (*ResourceLogs) Descriptor() ([]byte, []int) {
	return file_logs_proto_rawDescGZIP(), []int{1}
}

func (x *ResourceLogs) GetResource() *Resource {
	if x != nil {
		return x.Resource
	}
	return nil
}

func (x *ResourceLogs) GetScopeLogs() []*ScopeLogs {
	if x != nil {
		return x.ScopeLogs
	}
	return nil
}

func (x *ResourceLogs) GetSchemaUrl() string {
	if x != nil {
		return x.SchemaUrl
	}
	return ""
}

type Resource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Attributes             []*KeyValue `protobuf:"bytes,1,rep,name=attributes,proto3" json:"attributes,omitempty"`
	DroppedAttributesCount uint32      `protobuf:"varint,2,opt,name=dropped_attributes_count,json=droppedAttributesCount,proto3" json:"dropped_attributes_count,omitempty"`
}

func (x *Resource) Reset() {
	*x = Resource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logs_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_logs_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource.ProtoReflect.Descriptor instead.
func (*Resource) Descriptor() ([]byte, []int) {
	return file_logs_proto_rawDescGZIP(), []int{2}
}

func (x *Resource) GetAttributes() []*KeyValue {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *Resource) GetDroppedAttributesCount() uint32 {
	if x != nil {
		return x.DroppedAttributesCount
	}
	return 0
}

// A collection of Logs produced by a Scope.
type ScopeLogs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Scope *InstrumentationScope `protobuf:"bytes,1,opt,name=scope,proto3" json:"scope,omitempty"`
	// A list of log records.
	LogRecords []*LogRecord `protobuf:"bytes,2,rep,name=log_records,json=logRecords,proto3" json:"log_records,omitempty"`
	// This schema_url applies to all logs in the "logs" field.
	SchemaUrl string `protobuf:"bytes,3,opt,name=schema_url,json=schemaUrl,proto3" json:"schema_url,omitempty"`
}

func (x *ScopeLogs) Reset() {
	*x = ScopeLogs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logs_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScopeLogs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScopeLogs) ProtoMessage() {}

func (x *ScopeLogs) ProtoReflect() protoreflect.Message {
	mi := &file_logs_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScopeLogs.ProtoReflect.Descriptor instead.
func (*ScopeLogs) Descriptor() ([]byte, []int) {
	return file_logs_proto_rawDescGZIP(), []int{3}
}

func (x *ScopeLogs) GetScope() *InstrumentationScope {
	if x != nil {
		return x.Scope
	}
	return nil
}

func (x *ScopeLogs) GetLogRecords() []*LogRecord {
	if x != nil {
		return x.LogRecords
	}
	return nil
}

func (x *ScopeLogs) GetSchemaUrl() string {
	if x != nil {
		return x.SchemaUrl
	}
	return ""
}

type InstrumentationScope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name                   string      `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Version                string      `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Attributes             []*KeyValue `protobuf:"bytes,3,rep,name=attributes,proto3" json:"attributes,omitempty"`
	DroppedAttributesCount uint32      `protobuf:"varint,4,opt,name=dropped_attributes_count,json=droppedAttributesCount,proto3" json:"dropped_attributes_count,omitempty"`
}

func (x *InstrumentationScope) Reset() {
	*x = InstrumentationScope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logs_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstrumentationScope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstrumentationScope) ProtoMessage() {}

func (x *InstrumentationScope) ProtoReflect() protoreflect.Message {
	mi := &file_logs_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstrumentationScope.ProtoReflect.Descriptor instead.
func (*InstrumentationScope) Descriptor() ([]byte, []int) {
	return file_logs_proto_rawDescGZIP(), []int{4}
}

func (x *InstrumentationScope) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *InstrumentationScope) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *InstrumentationScope) GetAttributes() []*KeyValue {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *InstrumentationScope) GetDroppedAttributesCount() uint32 {
	if x != nil {
		return x.DroppedAttributesCount
	}
	return 0
}

type LogRecord struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TimeUnixNano           uint64         `protobuf:"fixed64,1,opt,name=time_unix_nano,json=timeUnixNano,proto3" json:"time_unix_nano,omitempty"`
	ObservedTimeUnixNano   uint64         `protobuf:"fixed64,11,opt,name=observed_time_unix_nano,json=observedTimeUnixNano,proto3" json:"observed_time_unix_nano,omitempty"`
	SeverityNumber         SeverityNumber `protobuf:"varint,2,opt,name=severity_number,json=severityNumber,proto3,enum=simple.SeverityNumber" json:"severity_number,omitempty"`
	SeverityText           string         `protobuf:"bytes,3,opt,name=severity_text,json=severityText,proto3" json:"severity_text,omitempty"`
	Attributes             []*KeyValue    `protobuf:"bytes,6,rep,name=attributes,proto3" json:"attributes,omitempty"`
	DroppedAttributesCount uint32         `protobuf:"varint,7,opt,name=dropped_attributes_count,json=droppedAttributesCount,proto3" json:"dropped_attributes_count,omitempty"`
	Flags                  uint32         `protobuf:"fixed32,8,opt,name=flags,proto3" json:"flags,omitempty"`
	TraceId                []byte         `protobuf:"bytes,9,opt,name=trace_id,json=traceId,proto3" json:"trace_id,omitempty"`
	SpanId                 []byte         `protobuf:"bytes,10,opt,name=span_id,json=spanId,proto3" json:"span_id,omitempty"`
}

func (x *LogRecord) Reset() {
	*x = LogRecord{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logs_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogRecord) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogRecord) ProtoMessage() {}

func (x *LogRecord) ProtoReflect() protoreflect.Message {
	mi := &file_logs_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogRecord.ProtoReflect.Descriptor instead.
func (*LogRecord) Descriptor() ([]byte, []int) {
	return file_logs_proto_rawDescGZIP(), []int{5}
}

func (x *LogRecord) GetTimeUnixNano() uint64 {
	if x != nil {
		return x.TimeUnixNano
	}
	return 0
}

func (x *LogRecord) GetObservedTimeUnixNano() uint64 {
	if x != nil {
		return x.ObservedTimeUnixNano
	}
	return 0
}

func (x *LogRecord) GetSeverityNumber() SeverityNumber {
	if x != nil {
		return x.SeverityNumber
	}
	return SeverityNumber_SEVERITY_NUMBER_UNSPECIFIED
}

func (x *LogRecord) GetSeverityText() string {
	if x != nil {
		return x.SeverityText
	}
	return ""
}

func (x *LogRecord) GetAttributes() []*KeyValue {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *LogRecord) GetDroppedAttributesCount() uint32 {
	if x != nil {
		return x.DroppedAttributesCount
	}
	return 0
}

func (x *LogRecord) GetFlags() uint32 {
	if x != nil {
		return x.Flags
	}
	return 0
}

func (x *LogRecord) GetTraceId() []byte {
	if x != nil {
		return x.TraceId
	}
	return nil
}

func (x *LogRecord) GetSpanId() []byte {
	if x != nil {
		return x.SpanId
	}
	return nil
}

type KeyValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string    `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value *AnyValue `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *KeyValue) Reset() {
	*x = KeyValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logs_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyValue) ProtoMessage() {}

func (x *KeyValue) ProtoReflect() protoreflect.Message {
	mi := &file_logs_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyValue.ProtoReflect.Descriptor instead.
func (*KeyValue) Descriptor() ([]byte, []int) {
	return file_logs_proto_rawDescGZIP(), []int{6}
}

func (x *KeyValue) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *KeyValue) GetValue() *AnyValue {
	if x != nil {
		return x.Value
	}
	return nil
}

type AnyValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Value:
	//	*AnyValue_StringValue
	//	*AnyValue_BoolValue
	//	*AnyValue_IntValue
	//	*AnyValue_DoubleValue
	//	*AnyValue_BytesValue
	Value isAnyValue_Value `protobuf_oneof:"value"`
}

func (x *AnyValue) Reset() {
	*x = AnyValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logs_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AnyValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AnyValue) ProtoMessage() {}

func (x *AnyValue) ProtoReflect() protoreflect.Message {
	mi := &file_logs_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AnyValue.ProtoReflect.Descriptor instead.
func (*AnyValue) Descriptor() ([]byte, []int) {
	return file_logs_proto_rawDescGZIP(), []int{7}
}

func (m *AnyValue) GetValue() isAnyValue_Value {
	if m != nil {
		return m.Value
	}
	return nil
}

func (x *AnyValue) GetStringValue() string {
	if x, ok := x.GetValue().(*AnyValue_StringValue); ok {
		return x.StringValue
	}
	return ""
}

func (x *AnyValue) GetBoolValue() bool {
	if x, ok := x.GetValue().(*AnyValue_BoolValue); ok {
		return x.BoolValue
	}
	return false
}

func (x *AnyValue) GetIntValue() int64 {
	if x, ok := x.GetValue().(*AnyValue_IntValue); ok {
		return x.IntValue
	}
	return 0
}

func (x *AnyValue) GetDoubleValue() float64 {
	if x, ok := x.GetValue().(*AnyValue_DoubleValue); ok {
		return x.DoubleValue
	}
	return 0
}

func (x *AnyValue) GetBytesValue() []byte {
	if x, ok := x.GetValue().(*AnyValue_BytesValue); ok {
		return x.BytesValue
	}
	return nil
}

type isAnyValue_Value interface {
	isAnyValue_Value()
}

type AnyValue_StringValue struct {
	StringValue string `protobuf:"bytes,1,opt,name=string_value,json=stringValue,proto3,oneof"`
}

type AnyValue_BoolValue struct {
	BoolValue bool `protobuf:"varint,2,opt,name=bool_value,json=boolValue,proto3,oneof"`
}

type AnyValue_IntValue struct {
	IntValue int64 `protobuf:"varint,3,opt,name=int_value,json=intValue,proto3,oneof"`
}

type AnyValue_DoubleValue struct {
	DoubleValue float64 `protobuf:"fixed64,4,opt,name=double_value,json=doubleValue,proto3,oneof"`
}

type AnyValue_BytesValue struct {
	BytesValue []byte `protobuf:"bytes,7,opt,name=bytes_value,json=bytesValue,proto3,oneof"`
}

func (*AnyValue_StringValue) isAnyValue_Value() {}

func (*AnyValue_BoolValue) isAnyValue_Value() {}

func (*AnyValue_IntValue) isAnyValue_Value() {}

func (*AnyValue_DoubleValue) isAnyValue_Value() {}

func (*AnyValue_BytesValue) isAnyValue_Value() {}

var File_logs_proto protoreflect.FileDescriptor

var file_logs_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x6c, 0x6f, 0x67, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x73, 0x69,
	0x6d, 0x70, 0x6c, 0x65, 0x22, 0x45, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x73, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x39, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6c, 0x6f, 0x67,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65,
	0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4c, 0x6f, 0x67, 0x73, 0x52, 0x0c, 0x72,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4c, 0x6f, 0x67, 0x73, 0x22, 0x8d, 0x01, 0x0a, 0x0c,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4c, 0x6f, 0x67, 0x73, 0x12, 0x2c, 0x0a, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x30, 0x0a, 0x0a, 0x73, 0x63,
	0x6f, 0x70, 0x65, 0x5f, 0x6c, 0x6f, 0x67, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x4c, 0x6f, 0x67,
	0x73, 0x52, 0x09, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x4c, 0x6f, 0x67, 0x73, 0x12, 0x1d, 0x0a, 0x0a,
	0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x55, 0x72, 0x6c, 0x22, 0x76, 0x0a, 0x08, 0x52,
	0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x30, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x69,
	0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x61,
	0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x38, 0x0a, 0x18, 0x64, 0x72, 0x6f,
	0x70, 0x70, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x5f,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x16, 0x64, 0x72, 0x6f,
	0x70, 0x70, 0x65, 0x64, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x22, 0x92, 0x01, 0x0a, 0x09, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x4c, 0x6f, 0x67,
	0x73, 0x12, 0x32, 0x0a, 0x05, 0x73, 0x63, 0x6f, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x72, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x63, 0x6f, 0x70, 0x65, 0x52, 0x05,
	0x73, 0x63, 0x6f, 0x70, 0x65, 0x12, 0x32, 0x0a, 0x0b, 0x6c, 0x6f, 0x67, 0x5f, 0x72, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73, 0x69, 0x6d,
	0x70, 0x6c, 0x65, 0x2e, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x52, 0x0a, 0x6c,
	0x6f, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x55, 0x72, 0x6c, 0x22, 0xb0, 0x01, 0x0a, 0x14, 0x49, 0x6e, 0x73,
	0x74, 0x72, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x63, 0x6f, 0x70,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x30, 0x0a, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x4b, 0x65, 0x79,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x73, 0x12, 0x38, 0x0a, 0x18, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x16, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x41, 0x74, 0x74, 0x72,
	0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x84, 0x03, 0x0a, 0x09,
	0x4c, 0x6f, 0x67, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x24, 0x0a, 0x0e, 0x74, 0x69, 0x6d,
	0x65, 0x5f, 0x75, 0x6e, 0x69, 0x78, 0x5f, 0x6e, 0x61, 0x6e, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x06, 0x52, 0x0c, 0x74, 0x69, 0x6d, 0x65, 0x55, 0x6e, 0x69, 0x78, 0x4e, 0x61, 0x6e, 0x6f, 0x12,
	0x35, 0x0a, 0x17, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x5f, 0x75, 0x6e, 0x69, 0x78, 0x5f, 0x6e, 0x61, 0x6e, 0x6f, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x06,
	0x52, 0x14, 0x6f, 0x62, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x55, 0x6e,
	0x69, 0x78, 0x4e, 0x61, 0x6e, 0x6f, 0x12, 0x3f, 0x0a, 0x0f, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69,
	0x74, 0x79, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x16, 0x2e, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x53, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74,
	0x79, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x52, 0x0e, 0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74,
	0x79, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x65, 0x76, 0x65, 0x72,
	0x69, 0x74, 0x79, 0x5f, 0x74, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x73, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x54, 0x65, 0x78, 0x74, 0x12, 0x30, 0x0a, 0x0a,
	0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x38,
	0x0a, 0x18, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x73, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x16, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75,
	0x74, 0x65, 0x73, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6c, 0x61, 0x67,
	0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x07, 0x52, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x19,
	0x0a, 0x08, 0x74, 0x72, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x07, 0x74, 0x72, 0x61, 0x63, 0x65, 0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x73, 0x70, 0x61,
	0x6e, 0x5f, 0x69, 0x64, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x73, 0x70, 0x61, 0x6e,
	0x49, 0x64, 0x22, 0x44, 0x0a, 0x08, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x26, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x41, 0x6e, 0x79, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xc0, 0x01, 0x0a, 0x08, 0x41, 0x6e, 0x79,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x5f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x0b, 0x73,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1f, 0x0a, 0x0a, 0x62, 0x6f,
	0x6f, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00,
	0x52, 0x09, 0x62, 0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1d, 0x0a, 0x09, 0x69,
	0x6e, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00,
	0x52, 0x08, 0x69, 0x6e, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x64, 0x6f,
	0x75, 0x62, 0x6c, 0x65, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01,
	0x48, 0x00, 0x52, 0x0b, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x21, 0x0a, 0x0b, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x0a, 0x62, 0x79, 0x74, 0x65, 0x73, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x42, 0x07, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2a, 0xc3, 0x05, 0x0a, 0x0e,
	0x53, 0x65, 0x76, 0x65, 0x72, 0x69, 0x74, 0x79, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x1f,
	0x0a, 0x1b, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45,
	0x52, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x19, 0x0a, 0x15, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42,
	0x45, 0x52, 0x5f, 0x54, 0x52, 0x41, 0x43, 0x45, 0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45,
	0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x54, 0x52,
	0x41, 0x43, 0x45, 0x32, 0x10, 0x02, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49,
	0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x54, 0x52, 0x41, 0x43, 0x45, 0x33,
	0x10, 0x03, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e,
	0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x54, 0x52, 0x41, 0x43, 0x45, 0x34, 0x10, 0x04, 0x12, 0x19,
	0x0a, 0x15, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45,
	0x52, 0x5f, 0x44, 0x45, 0x42, 0x55, 0x47, 0x10, 0x05, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x56,
	0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x44, 0x45, 0x42,
	0x55, 0x47, 0x32, 0x10, 0x06, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54,
	0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x44, 0x45, 0x42, 0x55, 0x47, 0x33, 0x10,
	0x07, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55,
	0x4d, 0x42, 0x45, 0x52, 0x5f, 0x44, 0x45, 0x42, 0x55, 0x47, 0x34, 0x10, 0x08, 0x12, 0x18, 0x0a,
	0x14, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52,
	0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x09, 0x12, 0x19, 0x0a, 0x15, 0x53, 0x45, 0x56, 0x45, 0x52,
	0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x32,
	0x10, 0x0a, 0x12, 0x19, 0x0a, 0x15, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e,
	0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x33, 0x10, 0x0b, 0x12, 0x19, 0x0a,
	0x15, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52,
	0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x34, 0x10, 0x0c, 0x12, 0x18, 0x0a, 0x14, 0x53, 0x45, 0x56, 0x45,
	0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x57, 0x41, 0x52, 0x4e,
	0x10, 0x0d, 0x12, 0x19, 0x0a, 0x15, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e,
	0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x57, 0x41, 0x52, 0x4e, 0x32, 0x10, 0x0e, 0x12, 0x19, 0x0a,
	0x15, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52,
	0x5f, 0x57, 0x41, 0x52, 0x4e, 0x33, 0x10, 0x0f, 0x12, 0x19, 0x0a, 0x15, 0x53, 0x45, 0x56, 0x45,
	0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x57, 0x41, 0x52, 0x4e,
	0x34, 0x10, 0x10, 0x12, 0x19, 0x0a, 0x15, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f,
	0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x11, 0x12, 0x1a,
	0x0a, 0x16, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45,
	0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x32, 0x10, 0x12, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45,
	0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x33, 0x10, 0x13, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49,
	0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x34,
	0x10, 0x14, 0x12, 0x19, 0x0a, 0x15, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e,
	0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x46, 0x41, 0x54, 0x41, 0x4c, 0x10, 0x15, 0x12, 0x1a, 0x0a,
	0x16, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52,
	0x5f, 0x46, 0x41, 0x54, 0x41, 0x4c, 0x32, 0x10, 0x16, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x56,
	0x45, 0x52, 0x49, 0x54, 0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x46, 0x41, 0x54,
	0x41, 0x4c, 0x33, 0x10, 0x17, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x45, 0x56, 0x45, 0x52, 0x49, 0x54,
	0x59, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x46, 0x41, 0x54, 0x41, 0x4c, 0x34, 0x10,
	0x18, 0x42, 0x0a, 0x5a, 0x08, 0x67, 0x65, 0x6e, 0x2f, 0x6c, 0x6f, 0x67, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_logs_proto_rawDescOnce sync.Once
	file_logs_proto_rawDescData = file_logs_proto_rawDesc
)

func file_logs_proto_rawDescGZIP() []byte {
	file_logs_proto_rawDescOnce.Do(func() {
		file_logs_proto_rawDescData = protoimpl.X.CompressGZIP(file_logs_proto_rawDescData)
	})
	return file_logs_proto_rawDescData
}

var file_logs_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_logs_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_logs_proto_goTypes = []interface{}{
	(SeverityNumber)(0),          // 0: simple.SeverityNumber
	(*LogsData)(nil),             // 1: simple.LogsData
	(*ResourceLogs)(nil),         // 2: simple.ResourceLogs
	(*Resource)(nil),             // 3: simple.Resource
	(*ScopeLogs)(nil),            // 4: simple.ScopeLogs
	(*InstrumentationScope)(nil), // 5: simple.InstrumentationScope
	(*LogRecord)(nil),            // 6: simple.LogRecord
	(*KeyValue)(nil),             // 7: simple.KeyValue
	(*AnyValue)(nil),             // 8: simple.AnyValue
}
var file_logs_proto_depIdxs = []int32{
	2,  // 0: simple.LogsData.resource_logs:type_name -> simple.ResourceLogs
	3,  // 1: simple.ResourceLogs.resource:type_name -> simple.Resource
	4,  // 2: simple.ResourceLogs.scope_logs:type_name -> simple.ScopeLogs
	7,  // 3: simple.Resource.attributes:type_name -> simple.KeyValue
	5,  // 4: simple.ScopeLogs.scope:type_name -> simple.InstrumentationScope
	6,  // 5: simple.ScopeLogs.log_records:type_name -> simple.LogRecord
	7,  // 6: simple.InstrumentationScope.attributes:type_name -> simple.KeyValue
	0,  // 7: simple.LogRecord.severity_number:type_name -> simple.SeverityNumber
	7,  // 8: simple.LogRecord.attributes:type_name -> simple.KeyValue
	8,  // 9: simple.KeyValue.value:type_name -> simple.AnyValue
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_logs_proto_init() }
func file_logs_proto_init() {
	if File_logs_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_logs_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogsData); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logs_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResourceLogs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logs_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Resource); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logs_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScopeLogs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logs_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstrumentationScope); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logs_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogRecord); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logs_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_logs_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AnyValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_logs_proto_msgTypes[7].OneofWrappers = []interface{}{
		(*AnyValue_StringValue)(nil),
		(*AnyValue_BoolValue)(nil),
		(*AnyValue_IntValue)(nil),
		(*AnyValue_DoubleValue)(nil),
		(*AnyValue_BytesValue)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_logs_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_logs_proto_goTypes,
		DependencyIndexes: file_logs_proto_depIdxs,
		EnumInfos:         file_logs_proto_enumTypes,
		MessageInfos:      file_logs_proto_msgTypes,
	}.Build()
	File_logs_proto = out.File
	file_logs_proto_rawDesc = nil
	file_logs_proto_goTypes = nil
	file_logs_proto_depIdxs = nil
}
