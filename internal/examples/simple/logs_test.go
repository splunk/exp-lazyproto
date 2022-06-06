package simple

import (
	"testing"

	gogolib "github.com/gogo/protobuf/proto"
	gogomsg "github.com/tigrannajaryan/exp-lazyproto/internal/examples/simple/gogo/gen/logs"
	googlemsg "github.com/tigrannajaryan/exp-lazyproto/internal/examples/simple/google/gen/logs"
	lazymsg "github.com/tigrannajaryan/exp-lazyproto/internal/examples/simple/lazy"
	googlelib "google.golang.org/protobuf/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tigrannajaryan/molecule"
)

func createAttr(k, v string) *gogomsg.KeyValue {
	return &gogomsg.KeyValue{
		Key: k,
		Value: &gogomsg.AnyValue{
			Value: &gogomsg.AnyValue_StringValue{
				StringValue: v,
			},
		},
	}
}

func createArrayValue(v string) *gogomsg.AnyValue {
	return &gogomsg.AnyValue{
		Value: &gogomsg.AnyValue_ArrayValue{
			ArrayValue: &gogomsg.ArrayValue{
				Values: []*gogomsg.AnyValue{
					{
						Value: &gogomsg.AnyValue_StringValue{
							StringValue: v,
						},
					},
				},
			},
		},
	}
}

func createKVList() *gogomsg.AnyValue {
	return &gogomsg.AnyValue{
		Value: &gogomsg.AnyValue_KvlistValue{
			KvlistValue: &gogomsg.KeyValueList{
				Values: []*gogomsg.KeyValue{
					createAttr("x", "10"),
					createAttr("y", "20"),
				},
			},
		},
	}
}

func createLogRecord(n int) *gogomsg.LogRecord {
	sl := &gogomsg.LogRecord{
		TimeUnixNano:   uint64(n * 10000),
		SeverityNumber: gogomsg.SeverityNumber(n % 25),
		Attributes: []*gogomsg.KeyValue{
			createAttr("http.method", "GET"),
			createAttr("http.url", "/checkout"),
			createAttr("http.server", "example.com"),
			createAttr("db.name", "postgres"),
			createAttr("host.name", "localhost"),
			{
				Key:   "multivalue",
				Value: createArrayValue("1.2.3.4"),
			},
		},
		DroppedAttributesCount: uint32(n),
	}

	if n == 0 {
		sl.SeverityText = "ERROR"
		sl.Flags = 1
		sl.ObservedTimeUnixNano = sl.TimeUnixNano + 200
		sl.SpanId = []byte{1, 2, 3, 4, 5}
		sl.TraceId = []byte{6, 7, 8, 9}
	}

	return sl
}

func createScopedLogs(n int) *gogomsg.ScopeLogs {
	sl := &gogomsg.ScopeLogs{
		Scope: &gogomsg.InstrumentationScope{
			Name:    "library",
			Version: "2.5",
			Attributes: []*gogomsg.KeyValue{
				createAttr("otel.profiling", "true"),
			},
		},
	}

	if n%2 == 0 {
		// Give half of the scopes a different attribute value.
		sl.Scope.Attributes[0] = createAttr("otel.profiling", "false")
	}

	for i := 0; i < 10; i++ {
		sl.LogRecords = append(sl.LogRecords, createLogRecord(i))
	}

	return sl
}

func createLogsData() *gogomsg.LogsData {
	src := &gogomsg.LogsData{}

	for i := 0; i < 10; i++ {
		rl := &gogomsg.ResourceLogs{
			Resource: &gogomsg.Resource{
				Attributes: []*gogomsg.KeyValue{
					createAttr("service.name", "checkout"),
					{
						Key:   "nested",
						Value: createKVList(),
					},
				},
				DroppedAttributesCount: 12,
			},
		}

		for j := 0; j < 10; j++ {
			rl.ScopeLogs = append(rl.ScopeLogs, createScopedLogs(j))
		}

		src.ResourceLogs = append(src.ResourceLogs, rl)
	}

	return src
}

func TestDecode(t *testing.T) {
	src := &gogomsg.LogsData{
		ResourceLogs: []*gogomsg.ResourceLogs{
			{
				Resource: &gogomsg.Resource{
					Attributes: []*gogomsg.KeyValue{
						createAttr("key1", "value1"),
						{
							Key:   "multivalue",
							Value: createArrayValue("1.2.3.4"),
						},
						{
							Key:   "nested",
							Value: createKVList(),
						},
					},
					DroppedAttributesCount: 12,
				},
				ScopeLogs: []*gogomsg.ScopeLogs{
					{
						SchemaUrl: "https://opentelemetry.io/schemas/1.0.0",
						Scope: &gogomsg.InstrumentationScope{
							Name:    "library",
							Version: "2.5",
							Attributes: []*gogomsg.KeyValue{
								createAttr("otel.profiling", "true"),
							},
						},
						LogRecords: []*gogomsg.LogRecord{
							{
								TimeUnixNano: 123,
								Attributes: []*gogomsg.KeyValue{
									createAttr("key2", "value2"),
								},
								DroppedAttributesCount: 234,
								SpanId:                 []byte{1, 2, 3, 4, 5},
								TraceId:                []byte{6, 7, 8, 9},
							},
						},
					},
				},
			},
		},
	}
	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(t, err)

	lazy, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
	require.NoError(t, err)

	rl := lazy.ResourceLogs()
	require.Len(t, rl, 1)

	resource := *rl[0].Resource()
	assert.EqualValues(t, 12, resource.DroppedAttributesCount())
	require.NotNil(t, resource)

	resAttrs := resource.Attributes()
	require.Len(t, resAttrs, 3)

	kvr := resAttrs[0]
	require.EqualValues(t, "key1", kvr.Key())
	require.EqualValues(t, lazymsg.AnyValueStringValue, kvr.Value().ValueType())
	require.EqualValues(t, "value1", kvr.Value().StringValue())

	kvr = resAttrs[1]
	require.EqualValues(t, "multivalue", kvr.Key())
	require.EqualValues(t, lazymsg.AnyValueArrayValue, kvr.Value().ValueType())
	arrayVals := kvr.Value().ArrayValue().Values()
	require.Len(t, arrayVals, 1)
	require.EqualValues(t, "1.2.3.4", arrayVals[0].StringValue())

	kvr = resAttrs[2]
	require.EqualValues(t, "nested", kvr.Key())
	require.EqualValues(t, lazymsg.AnyValueKvlistValue, kvr.Value().ValueType())
	kvVals := kvr.Value().KvlistValue().Values()
	require.Len(t, kvVals, 2)
	require.EqualValues(t, "x", kvVals[0].Key())
	require.EqualValues(t, "10", kvVals[0].Value().StringValue())
	require.EqualValues(t, "y", kvVals[1].Key())
	require.EqualValues(t, "20", kvVals[1].Value().StringValue())

	sls := rl[0].ScopeLogs()
	require.Len(t, sls, 1)

	sl := sls[0]
	logRecords := sl.LogRecords()
	require.Len(t, logRecords, 1)

	logRecord := logRecords[0]
	assert.EqualValues(t, 123, logRecord.TimeUnixNano())
	assert.EqualValues(t, 234, logRecord.DroppedAttributesCount())
	assert.EqualValues(t, []byte{1, 2, 3, 4, 5}, logRecord.SpanId())
	assert.EqualValues(t, []byte{6, 7, 8, 9}, logRecord.TraceId())

	attrs2 := logRecord.Attributes()
	require.Len(t, attrs2, 1)

	kv2 := attrs2[0]
	require.EqualValues(t, "key2", kv2.Key())
	require.EqualValues(t, lazymsg.AnyValueStringValue, kv2.Value().ValueType())
	require.EqualValues(t, "value2", kv2.Value().StringValue())

	ps := molecule.NewProtoStream()
	err = lazy.Marshal(ps)
	assert.NoError(t, err)

	lazyBytes, err := ps.BufferBytes()
	assert.NoError(t, err)
	assert.EqualValues(t, goldenWireBytes, lazyBytes)

	lazy.Free()
}

func TestLazyPassthrough(t *testing.T) {
	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(t, err)
	require.NotNil(t, goldenWireBytes)

	ps := molecule.NewProtoStream()
	lazy, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
	require.NoError(t, err)

	ps.Reset()
	err = lazy.Marshal(ps)
	require.NoError(t, err)

	lazyBytes, err := ps.BufferBytes()
	assert.NoError(t, err)
	assert.EqualValues(t, goldenWireBytes, lazyBytes)

	lazy.Free()
}

func BenchmarkGoogleMarshal(b *testing.B) {
	src := createLogsData()
	bytes, err := gogolib.Marshal(src)
	var ld googlemsg.LogsData
	err = googlelib.Unmarshal(bytes, &ld)
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bytes, err := googlelib.Marshal(&ld)
		require.NoError(b, err)
		require.NotNil(b, bytes)
	}
}

func BenchmarkGogoMarshal(b *testing.B) {
	src := createLogsData()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bytes, err := gogolib.Marshal(src)
		require.NoError(b, err)
		require.NotNil(b, bytes)
	}
}

func countAttrsLazy(lazy *lazymsg.LogsData) int {
	attrCount := 0
	rls := lazy.ResourceLogs()
	for _, rl := range rls {
		resource := *rl.Resource()

		attrs := resource.Attributes()
		attrCount += len(attrs)

		sls := rl.ScopeLogs()
		for _, sl := range sls {
			logRecords := sl.LogRecords()

			for _, logRecord := range logRecords {
				attrs2 := logRecord.Attributes()
				attrCount += len(attrs2)
			}
		}
	}
	return attrCount
}

func touchAll(lazy *lazymsg.LogsData) {
	rls := lazy.ResourceLogs()
	for _, rl := range rls {
		resource := *rl.Resource()

		attrs := resource.Attributes()
		for i := range attrs {
			attrs[i].SetKey(attrs[i].Key())
		}

		sls := rl.ScopeLogs()
		for _, sl := range sls {
			logRecords := sl.LogRecords()

			for _, logRecord := range logRecords {
				attrs2 := logRecord.Attributes()
				for i := range attrs2 {
					attrs2[i].SetKey(attrs2[i].Key())
				}
			}
		}
	}
}

func BenchmarkLazyMarshalUnchanged(b *testing.B) {
	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenWireBytes)

	lazy, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
	require.NoError(b, err)

	countAttrsLazy(lazy)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(b, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(b, err)
		assert.EqualValues(b, goldenWireBytes, lazyBytes)
	}
}

func BenchmarkLazyMarshalFullModified(b *testing.B) {
	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenWireBytes)

	lazy, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
	require.NoError(b, err)

	countAttrsLazy(lazy)
	touchAll(lazy)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(b, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(b, err)
		assert.EqualValues(b, goldenWireBytes, lazyBytes)
	}
}

func BenchmarkGoogleUnmarshal(b *testing.B) {
	src := createLogsData()

	bytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, bytes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var ld googlemsg.LogsData
		err := googlelib.Unmarshal(bytes, &ld)
		require.NoError(b, err)

		attrCount := 0
		for _, rl := range ld.ResourceLogs {
			attrCount += len(rl.Resource.Attributes)
			for _, sl := range rl.ScopeLogs {
				for _, lr := range sl.LogRecords {
					attrCount += len(lr.Attributes)
				}
			}
		}

		//require.EqualValues(b, 2010, attrCount)
	}
}

func BenchmarkGogoUnmarshal(b *testing.B) {
	src := createLogsData()

	bytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, bytes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var ld gogomsg.LogsData
		err := gogolib.Unmarshal(bytes, &ld)
		require.NoError(b, err)

		attrCount := 0
		for _, rl := range ld.ResourceLogs {
			attrCount += len(rl.Resource.Attributes)
			for _, sl := range rl.ScopeLogs {
				for _, lr := range sl.LogRecords {
					attrCount += len(lr.Attributes)
				}
			}
		}

		//require.EqualValues(b, 2010, attrCount)
	}
}

func BenchmarkLazyUnmarshalAndReadAll(b *testing.B) {
	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenWireBytes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lazy, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
		require.NoError(b, err)

		// Traverse all data to get it loaded. This is the worst case.
		countAttrsLazy(lazy)

		lazy.Free()
	}
}

func TestLazyUnmarshalAndReadAll(t *testing.T) {
	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(t, err)
	require.NotNil(t, goldenWireBytes)

	for i := 0; i < 2; i++ {
		lazy, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
		require.NoError(t, err)

		// Traverse all data to get it loaded. This is the worst case.
		countAttrsLazy(lazy)

		lazy.Free()
	}
}

func BenchmarkGooglePasssthrough(b *testing.B) {
	src := createLogsData()

	bytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, bytes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var ld googlemsg.LogsData
		err := googlelib.Unmarshal(bytes, &ld)
		require.NoError(b, err)

		destBytes, err := googlelib.Marshal(&ld)
		require.NoError(b, err)
		require.NotNil(b, destBytes)
	}
}

func BenchmarkGogoPasssthrough(b *testing.B) {
	src := createLogsData()

	bytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, bytes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var ld gogomsg.LogsData
		err := gogolib.Unmarshal(bytes, &ld)
		require.NoError(b, err)

		destBytes, err := gogolib.Marshal(src)
		require.NoError(b, err)
		require.NotNil(b, destBytes)
	}
}

func BenchmarkLazyPassthroughNoReadOrModify(b *testing.B) {
	// This is the best case scenario for passthrough. We don't read or modify any
	// data, just unmarshal and marshal it exactly as it is.

	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenWireBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		lazy, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
		require.NoError(b, err)

		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(b, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(b, err)
		assert.EqualValues(b, goldenWireBytes, lazyBytes)
	}
}

func BenchmarkLazyPassthroughFullReadNoModify(b *testing.B) {
	// This is the best case scenario for passthrough. We don't read or modify any
	// data, just unmarshal and marshal it exactly as it is.

	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenWireBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		lazy, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
		require.NoError(b, err)

		ps.Reset()
		err = lazy.Marshal(ps)
		countAttrsLazy(lazy)
		require.NoError(b, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(b, err)
		assert.EqualValues(b, goldenWireBytes, lazyBytes)

		lazy.Free()
	}
}

func BenchmarkLazyPassthroughFullModified(b *testing.B) {
	// This is the worst case scenario. We read of data, so lazy loading has no
	// performance benefit. We also modify all data, so we have to do full marshaling.

	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenWireBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		lazy, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
		require.NoError(b, err)

		// Touch all attrs
		touchAll(lazy)

		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(b, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(b, err)
		assert.EqualValues(b, goldenWireBytes, lazyBytes)

		lazy.Free()
	}
}

func TestLazyPassthroughFullModified(t *testing.T) {
	// This is the worst case scenario. We read of data, so lazy loading has no
	// performance benefit. We also modify all data, so we have to do full marshaling.

	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(t, err)
	require.NotNil(t, goldenWireBytes)

	ps := molecule.NewProtoStream()
	for i := 0; i < 3; i++ {
		lazy, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
		require.NoError(t, err)

		// Touch all attrs
		touchAll(lazy)

		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(t, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(t, err)
		assert.EqualValues(t, goldenWireBytes, lazyBytes)

		lazy.Free()
	}
}

func BenchmarkGogoInspectScopeAttr(b *testing.B) {
	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenWireBytes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var lazy gogomsg.LogsData
		err := gogolib.Unmarshal(goldenWireBytes, &lazy)
		require.NoError(b, err)

		foundCount := 0
		for _, rl := range lazy.ResourceLogs {
			for _, sl := range rl.ScopeLogs {
				if sl.Scope == nil {
					continue
				}
				for _, attr := range sl.Scope.Attributes {
					if attr.Key == "otel.profiling" &&
						attr.GetValue().GetStringValue() == "true" {
						foundCount++
					}
				}
			}
		}
		assert.Equal(b, 50, foundCount)

		destBytes, err := gogolib.Marshal(&lazy)
		require.NoError(b, err)
		assert.EqualValues(b, goldenWireBytes, destBytes)
	}
}

func BenchmarkLazyInspectScopeAttr(b *testing.B) {
	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenWireBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		inputMsg, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
		require.NoError(b, err)

		foundCount := 0
		for _, rl := range inputMsg.ResourceLogs() {
			for _, sl := range rl.ScopeLogs() {
				if sl.Scope() == nil {
					continue
				}
				for _, attr := range sl.Scope().Attributes() {
					if attr.Key() == "otel.profiling" &&
						attr.Value().ValueType() == lazymsg.AnyValueStringValue &&
						attr.Value().StringValue() == "true" {
						foundCount++
					}
				}
			}
		}
		assert.Equal(b, 50, foundCount)

		ps.Reset()
		err = inputMsg.Marshal(ps)
		require.NoError(b, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(b, err)
		assert.EqualValues(b, goldenWireBytes, lazyBytes)

		inputMsg.Free()
	}
}

func BenchmarkGogoFilterScopeAttr(b *testing.B) {
	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenWireBytes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var lazy gogomsg.LogsData
		err := gogolib.Unmarshal(goldenWireBytes, &lazy)
		require.NoError(b, err)

		foundCount := 0
		for _, rl := range lazy.ResourceLogs {
			for j := 0; j < len(rl.ScopeLogs); j++ {
				sl := rl.ScopeLogs[j]
				if sl.Scope == nil {
					continue
				}
				found := false
				for _, attr := range sl.Scope.Attributes {
					if attr.Key == "otel.profiling" && attr.GetValue().GetStringValue() == "true" {
						foundCount++
						found = true
						break
					}
				}
				if found {
					rl.ScopeLogs = append(rl.ScopeLogs[:j], rl.ScopeLogs[j+1:]...)
					j--
				}
			}
		}
		assert.Equal(b, 50, foundCount)

		destBytes, err := gogolib.Marshal(&lazy)
		require.NoError(b, err)
		assert.NotNil(b, destBytes)
	}
}

func BenchmarkLazyFilterScopeAttr(b *testing.B) {
	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenWireBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		inputMsg, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
		require.NoError(b, err)

		foundCount := 0
		for _, rl := range inputMsg.ResourceLogs() {
			rl.ScopeLogsRemoveIf(
				func(sl *lazymsg.ScopeLogs) bool {
					if sl.Scope() == nil {
						return false
					}
					for _, attr := range sl.Scope().Attributes() {
						if attr.Key() == "otel.profiling" &&
							attr.Value().ValueType() == lazymsg.AnyValueStringValue &&
							attr.Value().StringValue() == "true" {
							foundCount++
							return true
						}
					}
					return false
				},
			)
		}
		assert.Equal(b, 50, foundCount)

		ps.Reset()
		err = inputMsg.Marshal(ps)
		require.NoError(b, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(b, err)
		assert.NotNil(b, lazyBytes)

		inputMsg.Free()
	}
}

func gogoBatch(b *testing.B, inputWireBytes []byte) (batchedWireBytes []byte) {
	var inputMsg [10]gogomsg.LogsData
	var outputMsg gogomsg.LogsData

	for j := 0; j < 10; j++ {
		err := gogolib.Unmarshal(inputWireBytes, &inputMsg[j])
		require.NoError(b, err)

		outputMsg.ResourceLogs = append(
			outputMsg.ResourceLogs, inputMsg[j].ResourceLogs...,
		)
	}

	batchedWireBytes, err := gogolib.Marshal(&outputMsg)
	require.NoError(b, err)
	assert.NotNil(b, batchedWireBytes)

	return batchedWireBytes
}

func BenchmarkGogoBatch(b *testing.B) {
	src := createLogsData()

	inputWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, inputWireBytes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gogoBatch(b, inputWireBytes)
	}
}

func BenchmarkLazyBatch(b *testing.B) {
	src := createLogsData()

	inputWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, inputWireBytes)

	goldenBatchedBytes := gogoBatch(b, inputWireBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()

	for i := 0; i < b.N; i++ {
		var inputMsg [10]*lazymsg.LogsData
		outputMsg := lazymsg.LogsData{}
		var resourceLogs []*lazymsg.ResourceLogs

		for j := 0; j < 10; j++ {
			inputMsg[j], err = lazymsg.UnmarshalLogsData(inputWireBytes)
			require.NoError(b, err)

			resourceLogs = append(
				resourceLogs, inputMsg[j].ResourceLogs()...,
			)
		}

		outputMsg.SetResourceLogs(resourceLogs)

		ps.Reset()
		err := outputMsg.Marshal(ps)
		require.NoError(b, err)

		destBytes, err := ps.BufferBytes()
		require.NoError(b, err)
		assert.NotNil(b, destBytes)

		if i == 0 {
			assert.EqualValues(b, goldenBatchedBytes, destBytes)
		}
	}
}

func BenchmarkLazyTouchAll(b *testing.B) {
	src := createLogsData()

	goldenWireBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenWireBytes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lazy, err := lazymsg.UnmarshalLogsData(goldenWireBytes)
		require.NoError(b, err)

		countAttrsLazy(lazy)
		touchAll(lazy)
	}
}
