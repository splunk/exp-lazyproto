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
| Pass_ModifyAll | A call to Unmarshal(), then reading and changing all messages and nested messages (without introducing new messages or deleting existing messages). This an "inspect all, update and passthrough" scenario. |
| Inspect_ScopeAttr | A call to Unmarshal(), then read the Scope attributes to see if a specific attribute is found, followed by Marshal() call. This is a "inspect the Scopes and passthrough" scenario. |
| Inspect_LogAttr | A call to Unmarshal(), then read the LogRecord attributes to see if a specific attribute is found, followed by Marshal() call. This is a "inspect the LogRecords and passthrough" scenario. |
| Filter_ScopeAttr | A call to Unmarshal(), then read the Scope attributes and if a specific attribute is found drop that particular Scope and the embedded messages, followed by Marshal() call on the remaining data. This is a "filter based on the Scope attribute" scenario. |
| Batch | A call to Unmarshal() for a number of messages, then stitching the messages together into one message, followed by Marshal() called of the resulting message. |

