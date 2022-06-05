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

func createLogRecord(n int) *gogomsg.LogRecord {
	sl := &gogomsg.LogRecord{
		Attributes: []*gogomsg.KeyValue{
			{
				Key:   "http.method",
				Value: "GET",
			},
			{
				Key:   "http.url",
				Value: "/checkout",
			},
			{
				Key:   "http.server",
				Value: "example.com",
			},
			{
				Key:   "db.name",
				Value: "postgres",
			},
			{
				Key:   "host.name",
				Value: "localhost",
			},
		},
		DroppedAttributesCount: 12,
	}

	return sl
}

func createScopedLogs(n int) *gogomsg.ScopeLogs {
	sl := &gogomsg.ScopeLogs{
		Scope: &gogomsg.InstrumentationScope{
			Name:    "library",
			Version: "2.5",
			Attributes: []*gogomsg.KeyValue{
				{
					Key:   "otel.profiling",
					Value: "true",
				},
			},
		},
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
					{
						Key:   "service.name",
						Value: "checkout",
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
						{
							Key:   "key1",
							Value: "value1",
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
								{
									Key:   "otel.profiling",
									Value: "true",
								},
							},
						},
						LogRecords: []*gogomsg.LogRecord{
							{
								TimeUnixNano: 123,
								Attributes: []*gogomsg.KeyValue{
									{
										Key:   "key2",
										Value: "value2",
									},
								},
								DroppedAttributesCount: 234,
							},
						},
					},
				},
			},
		},
	}
	goldenMarshalBytes, err := gogolib.Marshal(src)
	require.NoError(t, err)

	lazy := lazymsg.NewLogsData(goldenMarshalBytes)

	rl := lazy.ResourceLogs()
	require.Len(t, rl, 1)

	resource := *rl[0].Resource()
	assert.EqualValues(t, 12, resource.DroppedAttributesCount())
	require.NotNil(t, resource)

	attrs := resource.Attributes()
	require.Len(t, attrs, 1)

	kv1 := attrs[0]
	require.EqualValues(t, "key1", kv1.Key())
	require.EqualValues(t, "value1", kv1.Value())

	sls := rl[0].ScopeLogs()
	require.Len(t, sls, 1)

	sl := sls[0]
	logRecords := sl.LogRecords()
	require.Len(t, logRecords, 1)

	logRecord := logRecords[0]
	assert.EqualValues(t, 123, logRecord.TimeUnixNano())
	assert.EqualValues(t, 234, logRecord.DroppedAttributesCount())
	attrs2 := logRecord.Attributes()
	require.Len(t, attrs, 1)

	kv2 := attrs2[0]
	require.EqualValues(t, "key2", kv2.Key())
	require.EqualValues(t, "value2", kv2.Value())

	ps := molecule.NewProtoStream()
	lazy.Marshal(ps)

	lazyBytes, err := ps.BufferBytes()
	assert.NoError(t, err)
	assert.EqualValues(t, goldenMarshalBytes, lazyBytes)

	lazy.Free()
}

func TestLazyPassthrough(t *testing.T) {
	src := createLogsData()

	goldenMarshalBytes, err := gogolib.Marshal(src)
	require.NoError(t, err)
	require.NotNil(t, goldenMarshalBytes)

	ps := molecule.NewProtoStream()
	lazy := lazymsg.NewLogsData(goldenMarshalBytes)
	ps.Reset()
	err = lazy.Marshal(ps)
	require.NoError(t, err)

	lazyBytes, err := ps.BufferBytes()
	assert.NoError(t, err)
	assert.EqualValues(t, goldenMarshalBytes, lazyBytes)

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

	goldenMarshalBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenMarshalBytes)

	lazy := lazymsg.NewLogsData(goldenMarshalBytes)
	countAttrsLazy(lazy)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(b, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(b, err)
		assert.EqualValues(b, goldenMarshalBytes, lazyBytes)
	}
}

func BenchmarkLazyMarshalFullModified(b *testing.B) {
	src := createLogsData()

	goldenMarshalBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenMarshalBytes)

	lazy := lazymsg.NewLogsData(goldenMarshalBytes)
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
		assert.EqualValues(b, goldenMarshalBytes, lazyBytes)
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

	bytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, bytes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lazy := lazymsg.NewLogsData(bytes)

		// Traverse all data to get it loaded. This is the worst case.
		countAttrsLazy(lazy)

		lazy.Free()
	}
}

func TestLazyUnmarshalAndReadAll(t *testing.T) {
	src := createLogsData()

	bytes, err := gogolib.Marshal(src)
	require.NoError(t, err)
	require.NotNil(t, bytes)

	for i := 0; i < 2; i++ {
		lazy := lazymsg.NewLogsData(bytes)

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

	goldenMarshalBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenMarshalBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		lazy := lazymsg.NewLogsData(goldenMarshalBytes)
		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(b, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(b, err)
		assert.EqualValues(b, goldenMarshalBytes, lazyBytes)
	}
}

func BenchmarkLazyPassthroughFullReadNoModify(b *testing.B) {
	// This is the best case scenario for passthrough. We don't read or modify any
	// data, just unmarshal and marshal it exactly as it is.

	src := createLogsData()

	goldenMarshalBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenMarshalBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		lazy := lazymsg.NewLogsData(goldenMarshalBytes)
		ps.Reset()
		err = lazy.Marshal(ps)
		countAttrsLazy(lazy)
		require.NoError(b, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(b, err)
		assert.EqualValues(b, goldenMarshalBytes, lazyBytes)

		lazy.Free()
	}
}

func BenchmarkLazyPassthroughFullModified(b *testing.B) {
	// This is the worst case scenario. We read of data, so lazy loading has no
	// performance benefit. We also modify all data, so we have to do full marshaling.

	src := createLogsData()

	goldenMarshalBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenMarshalBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		lazy := lazymsg.NewLogsData(goldenMarshalBytes)

		// Touch all attrs
		touchAll(lazy)

		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(b, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(b, err)
		assert.EqualValues(b, goldenMarshalBytes, lazyBytes)

		lazy.Free()
	}
}

func TestLazyPassthroughFullModified(t *testing.T) {
	// This is the worst case scenario. We read of data, so lazy loading has no
	// performance benefit. We also modify all data, so we have to do full marshaling.

	src := createLogsData()

	goldenMarshalBytes, err := gogolib.Marshal(src)
	require.NoError(t, err)
	require.NotNil(t, goldenMarshalBytes)

	ps := molecule.NewProtoStream()
	for i := 0; i < 3; i++ {
		lazy := lazymsg.NewLogsData(goldenMarshalBytes)

		// Touch all attrs
		touchAll(lazy)

		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(t, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(t, err)
		assert.EqualValues(t, goldenMarshalBytes, lazyBytes)

		lazy.Free()
	}
}

func BenchmarkGogoInspectScopeAttr(b *testing.B) {
	src := createLogsData()

	goldenMarshalBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenMarshalBytes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var lazy gogomsg.LogsData
		err := gogolib.Unmarshal(goldenMarshalBytes, &lazy)
		require.NoError(b, err)

		foundCount := 0
		for _, rl := range lazy.ResourceLogs {
			for _, sl := range rl.ScopeLogs {
				if sl.Scope == nil {
					continue
				}
				for _, attr := range sl.Scope.Attributes {
					if attr.Key == "otel.profiling" && attr.Value == "true" {
						foundCount++
					}
				}
			}
		}
		assert.Equal(b, 100, foundCount)

		destBytes, err := gogolib.Marshal(&lazy)
		require.NoError(b, err)
		assert.EqualValues(b, goldenMarshalBytes, destBytes)
	}
}

func BenchmarkLazyInspectScopeAttr(b *testing.B) {
	src := createLogsData()

	goldenMarshalBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, goldenMarshalBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		lazy := lazymsg.NewLogsData(goldenMarshalBytes)

		foundCount := 0
		for _, rl := range lazy.ResourceLogs() {
			for _, sl := range rl.ScopeLogs() {
				if sl.Scope() == nil {
					continue
				}
				for _, attr := range sl.Scope().Attributes() {
					if attr.Key() == "otel.profiling" && attr.Value() == "true" {
						foundCount++
					}
				}
			}
		}
		assert.Equal(b, 100, foundCount)

		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(b, err)

		lazyBytes, err := ps.BufferBytes()
		assert.NoError(b, err)
		assert.EqualValues(b, goldenMarshalBytes, lazyBytes)

		lazy.Free()
	}
}
