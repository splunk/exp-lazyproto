# LazyProto

This is an experimental implementation of Protobufs in Go that combines a few
interesting techniques to achieve significantly better performance than the official
Google [Go Protobuf library](https://pkg.go.dev/google.golang.org/protobuf) and better
than the [GoGo Protobuf library](https://github.com/gogo/protobuf).

Warning: not for production use. This is a research library.

## TLDR; Performance Comparison

Here is how this implementation compares to Google's
[Go Protobuf library](https://pkg.go.dev/google.golang.org/protobuf) and to 
[GoGo Protobuf library](https://github.com/gogo/protobuf).

![Comparison Benchmarks](internal/images/comparison.png)

Above we have a comparison between Google, GoGo and LazyProto implementations. LazyProto
is benchmarked twice: with and without unmarshalling validations (more on it later in
this document).

The bars show the time consumed when performing a specific benchmark. The benchmark
code can be found in [this file](internal/examples/simple/logs_test.go), here is a
rough explanation of each benchmark:

| Benchmark | Description                                           |
|--|------------------------------------------------------------------------------------------------------------------------------------------|
| Unmarshal | A call to Unmarshal() function that accepts a []byte slice and returns a message structure that allows the message fields to be accessed. |
| Unmarshal_AndReadAll | A call to Unmarshal(), followed by accessing every attribute of every message and embedded message. This forces LazyProto to decode every message and every field, essentially negating any benefits of lazy decoding. |
| Marshal_Unchanged | A call to Marshal(), assuming none of the fields of the message or of embedded messages where changed. |
| Marshal_ModifyAll | A call to Marshal(), assuming all the fields of the message and of all embedded messages where changed before the Marshal() call. The time to change the fields is not included in the benchmark. |
| Pass_NoReadNoModify | A call to Unmarshal(), immediately followed by a call to Marshal(). This is a passthrough scenario. None of the message data is read or changed. |
| Pass_ReadAllNoModify | A call to Unmarshal(), then reading all messages and nested messages, without changing any fields, followed by a call to Marshal(). This an "inspect all and passthrough" scenario. |
| Pass_ModifyAll | A call to Unmarshal(), then reading and changing all messages and nested messages (without introducing new messages or deleting existing messages). This an "inspect all, update and passthrough" scenario. This is the worst case for LazyProto since it prevents the lazy decoding and lazy decoding capabilities from improving performance (in fact these capabilities become a performance overhead) |
| Inspect_ScopeAttr | A call to Unmarshal(), then read the Scope attributes to see if a specific attribute is found, followed by Marshal() call. This is a "inspect the Scopes and passthrough" scenario. |
| Inspect_LogAttr | A call to Unmarshal(), then read the LogRecord attributes to see if a specific attribute is found, followed by Marshal() call. This is a "inspect the LogRecords and passthrough" scenario. |
| Filter_ScopeAttr | A call to Unmarshal(), then read the Scope attributes and if a specific attribute is found drop that particular Scope and the embedded messages, followed by Marshal() call on the remaining data. This is a "filter based on the Scope attribute" scenario. |
| Batch | A call to Unmarshal() for a number of messages, then stitching the messages together into one message, followed by Marshal() called of the resulting message. |

## How it Works

LazyProto uses a few techniques to improve the performance compared to other Protobuf
implementations in Go. Below we describe each of these techniques.

### Lazy Decoding

We utilize lazy decoding, which means that the `Unmarshal()` operation does not necessarily
decode all data from the wire representation into distinct memory structures. Instead,
LazyProto performs the decoding on-demand, when a particular piece of data is accessed.

For every message that contains embedded messages LazyProto tracks whether the particular
embedded message is decoded or no. Initially all embedded messages are undecoded. The
undecoded messages are represented as byte slices that point to a subslice of the original
slice that was passed to the `Unmarshal()` call. The subslice is the wire representation
of the message. When the message is accessed via a getter that needs to return the
message as a Go struct LazyProto checks if the message is decoded. If the message is
not yet decoded LazyProto will decode from the wire representation bytes into the
Go struct, populating the corresponding fields in the struct. Any nested embedded
messages will remain undecoded and will be only decoded when they in turn are accessed.

To track whether a particular embedded message is decoded or no LazyProto maintains
a `_flags` field in the parent message. For each field that represents an embedded message
the `_flags` field allocated one bit that indicates whether that particular embedded
message is decoded.

The getters of the field are implemented like this:

```go
func (m *LogsData) ResourceLogs() (r []*ResourceLogs) {
	if m._flags&flags_LogsData_ResourceLogs_Decoded == 0 {
		m.decodeResourceLogs()
	}
	return m.resourceLogs
}
```

Here the corresponding "decoded" bit in the `_flags` is checked. If it is set then the
`LogsData.resourceLogs` field is already decoded and is returned immediately. If the
bit is unset then the `decodeResourceLogs()` function is called to decode from the
wire representation  into the `LogsData.resourceLogs` struct field first.

This technique introduces a minimal overhead compared to the direct access to the
`resourceLogs` struct field. The `ResourceLogs()` is normally inlined, so any call
to it when the field is already decoded adds only a couple machine instructions to check
the `_flags` bit.

The `BenchmarkLazy_CountAttrs` benchmarks the access to the fields. 
`BenchmarkGogo_CountAttrs` performs the exact same benchmark by accessing the struct
fields generated by GoGo libray directly. Here is the comparison:

```go
name               time/op
Gogo_CountAttrs-8  26.3µs ± 1%
Lazy_CountAttrs-8  29.0µs ± 1%
```

The performance penalty due to using getters is negligible in most use cases.

### Lazy Encoding

We utilize lazy encoding, which means that the `Marshal()` operation does not necessarily
encode all data from the distinct memory structures into wire representation. Instead,
LazyProto will reuse the existing wire representation bytes for the message if the
message was previously unmarshalled from the wire representation and was not modified
since then.

For every message that LazyProto tracks whether the particular message has a
corresponding wire representation as a byte slice. This information is set when the
message is unmarshalled. When the message is modified the reference to the byte slice is
reset. In the `Marshal()` call we check if the byte slice is present and simply copy
the bytes to the output as is, otherwise we perform full encoding of the message
struct fields into wire bytes.

Lazy encoding is performed per message. This means that in the tree of messages
that point to each other we will use lazy encoding for unmodified messages and will
only perform full encoding for modified messages or messages that were created locally
and where not unmarshalled.

The setters of the field are implemented like this:

```go
// SetSchemaUrl sets the value of the schemaUrl.
func (m *ResourceLogs) SetSchemaUrl(v string) {
	m.schemaUrl = v

	// Mark this message modified, if not already.
	m._protoMessage.MarkModified()
}
```

Here the `m._protoMessage.MarkModified()` call will reset the reference to the byte
slice with the wire representation (if any), which will force the subsequent `Marshal()`
call to perform full encoding of the `ResourceLogs` message.

The combination of lazy decoding and lazy encoding results in very large savings
in passthrough scenarios when we unmarshal the data, perform no or minimal amount of
reading and modifying of the data and then marshal it. This is a common scenario
for intermediary services such as 
[OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector)

### Struct Pooling

Go Protobuf libraries allocate structs on the heap when unmarshalling. This typically
results in a large number of short-lived heap allocations that need to be
garbage-collected immediately after the incoming request is handled. These allocations
and subsequent garbage collection consume a significant portion of the entire processing
time and also result in a large number of total memory allocated until the garbage
collection reclaims the unused structs.

This is a very wasteful approach for most server use cases. There is an
[existing proposal](https://github.com/golang/go/issues/51317) to add arenas to Go
and one of the target use cases is the Protobuf request handling.

We could implement our own arenas, however we chose not to. The reason for that is
that arenas do not work well when the lifetime of the unmarshalled objects is not
precisely tied to the single incoming request. For example, in
[OpenTelemetry Collector](https://github.com/open-telemetry/opentelemetry-collector)
multiple incoming requests may be combined into a single batch and then processed
and marshalled together. The typical approach where the unmarshalled data from a single
incoming request is placed in one arena and the arena discarded immediately after the
incoming request processing is done doesn't work well in this case. There is no clear
moment in time where it is known that the particular arena can be safely discarded.

Instead of arenas we chose to implement message struct pools. Every struct that needs
to be allocated is picked from a pool of available structs. Every struct that is
no longer needed is returned to the pool of available structs. In intermediary
services like OpenTelemetry Collector we know when a particular message struct
processing is done (typically when the message is sent out from the Collector) so
it is easy to return the struct to the pool.

To return the message struct to the pool the user must call the `Free()` method of
the message struct:

```go
logsData.Free()
```

The `Free()` method returns the struct of the message and of all the embedded messages
reachable from the message to the struct pools.

Internally the pools implement an interface that allows to obtain one or more message
structs and to return them to the pool when the `Free()` method is called. The pool
interface looks like this:

```go
type Pool interface {
    Get() *T
	GetSlice(slice []*T)
    Release(elem *T)
    ReleaseSlice(slice []*T)
}
```

The pool implementations are concurrent-safe so the `Free()` method can be called
anytime without any additional synchronization, provided that the message, the embedded
messages or any of the message fields are no longer referenced from elsewhere.

