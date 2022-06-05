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
	"github.com/tigrannajaryan/molecule/src/protowire"
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
	sl := &gogomsg.ScopeLogs{}

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

	lazy := NewLogsData(goldenMarshalBytes)

	rl := lazy.ResourceLogs()
	require.Len(t, rl, 1)

	resource := *rl[0].Resource()
	assert.EqualValues(t, 12, resource.DroppedAttributesCount)
	require.NotNil(t, resource)

	attrs := resource.Attributes()
	require.Len(t, attrs, 1)

	kv1 := attrs[0]
	require.EqualValues(t, "key1", kv1.key)
	require.EqualValues(t, "value1", kv1.value)

	sls := rl[0].ScopeLogs()
	require.Len(t, sls, 1)

	sl := sls[0]
	logRecords := sl.LogRecords()
	require.Len(t, logRecords, 1)

	logRecord := logRecords[0]
	assert.EqualValues(t, 123, logRecord.timeUnixNano)
	assert.EqualValues(t, 234, logRecord.droppedAttributesCount)
	attrs2 := logRecord.Attributes()
	require.Len(t, attrs, 1)

	kv2 := attrs2[0]
	require.EqualValues(t, "key2", kv2.key)
	require.EqualValues(t, "value2", kv2.value)

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

func BenchmarkGogoMarshal(b *testing.B) {
	src := createLogsData()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		bytes, err := gogolib.Marshal(src)
		require.NoError(b, err)
		require.NotNil(b, bytes)
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

func countAttrs(lazy *LogsData) int {
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

//func BenchmarkLazyUnmarshalAndReadAll(b *testing.B) {
//	src := createLogsData()
//
//	bytes, err := gogolib.Marshal(src)
//	require.NoError(b, err)
//	require.NotNil(b, bytes)
//
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		lazy := NewLogsData(bytes)
//
//		// Traverse all data to get it loaded. This is the worst case.
//		countAttrs(lazy)
//
//		lazy.Free()
//	}
//}

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

func TestLazyGenUnmarshalAndReadAll(t *testing.T) {
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

//func BenchmarkLazyGenPassthroughNoReadOrModify(b *testing.B) {
//	// This is the best case scenario for passthrough. We don't read or modify any
//	// data, just unmarshal and marshal it exactly as it is.
//
//	src := createLogsData()
//
//	marshalBytes, err := gogolib.Marshal(src)
//	require.NoError(b, err)
//	require.NotNil(b, marshalBytes)
//
//	b.ResetTimer()
//
//	ps := molecule.NewProtoStream()
//	for i := 0; i < b.N; i++ {
//		lazy := lazymsg.NewLogsData(marshalBytes)
//		ps.Reset()
//		err = lazy.Marshal(ps)
//		require.NoError(b, err)
//		require.NotNil(b, ps.BufferBytes())
//	}
//}

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

//func BenchmarkKeyValueMarshal(b *testing.B) {
//	kv := KeyValue{
//		key:   "key",
//		value: "val",
//	}
//
//	ps := molecule.NewProtoStream()
//	for i := 0; i < b.N; i++ {
//		ps.Reset()
//		kv.Marshal(ps)
//		//require.Len(b, ps.BufferBytes(), 10)
//	}
//}
//
//func BenchmarkResourceMarshal(b *testing.B) {
//	kv := KeyValue{
//		key:   "key",
//		value: "val",
//	}
//	res := Resource{
//		attributes:             []*KeyValue{&kv},
//		DroppedAttributesCount: 0,
//	}
//
//	ps := molecule.NewProtoStream()
//	for i := 0; i < b.N; i++ {
//		ps.Reset()
//		res.Marshal(ps)
//	}
//}
//
//func BenchmarkLogRecordMarshal(b *testing.B) {
//	kv := KeyValue{
//		key:   "key",
//		value: "val",
//	}
//	lr := LogRecord{
//		attributes: []*KeyValue{&kv},
//	}
//
//	ps := molecule.NewProtoStream()
//	for i := 0; i < b.N; i++ {
//		ps.Reset()
//		lr.Marshal(ps)
//	}
//}

func BenchmarkAppendVarintProtowire(b *testing.B) {
	bts := make([]byte, 0, b.N*10)
	value := uint64(123)
	for i := 0; i < b.N; i++ {
		bts = protowire.AppendVarint(bts, value)
	}
}

//func BenchmarkWriteUint32(b *testing.B) {
//	val := uint32(123)
//	ps := molecule.NewProtoStream()
//	for i := 0; i < b.N; i++ {
//		ps.Reset()
//		ps.Uint32(1, val)
//	}
//}
//
//func BenchmarkWriteUint32Prepared(b *testing.B) {
//	val := uint32(123)
//	ps := molecule.NewProtoStream()
//	k := molecule.PrepareUint32Field(1)
//	for i := 0; i < b.N; i++ {
//		ps.Reset()
//		ps.Uint32Prepared(k, val)
//	}
//}
