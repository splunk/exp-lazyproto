package simple

import (
	"testing"

	gogolib "github.com/gogo/protobuf/proto"
	"github.com/richardartoul/molecule"
	"github.com/richardartoul/molecule/src/protowire"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	googlemsg "github.com/tigrannajaryan/exp-lazyproto/internal/examples/simple/google/gen/logs"
	googlelib "google.golang.org/protobuf/proto"

	gogomsg "github.com/tigrannajaryan/exp-lazyproto/internal/examples/simple/gogo/gen/logs"
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
	marshalledBytes, err := gogolib.Marshal(src)
	require.NoError(t, err)

	lazy := NewLogsData(marshalledBytes)

	rl := *lazy.GetResourceLogs()
	require.Len(t, rl, 1)

	resource := *rl[0].GetResource()
	assert.EqualValues(t, 12, resource.DroppedAttributesCount)
	require.NotNil(t, resource)

	attrs := *resource.GetAttributes()
	require.Len(t, attrs, 1)

	kv1 := attrs[0]
	require.EqualValues(t, "key1", kv1.key)
	require.EqualValues(t, "value1", kv1.value)

	sls := *rl[0].GetScopeLogs()
	require.Len(t, sls, 1)

	sl := sls[0]
	logRecords := *sl.GetLogRecords()
	require.Len(t, logRecords, 1)

	logRecord := logRecords[0]
	assert.EqualValues(t, 123, logRecord.timeUnixNano)
	assert.EqualValues(t, 234, logRecord.droppedAttributesCount)
	attrs2 := *logRecord.GetAttributes()
	require.Len(t, attrs, 1)

	kv2 := attrs2[0]
	require.EqualValues(t, "key2", kv2.key)
	require.EqualValues(t, "value2", kv2.value)

	ps := molecule.NewProtoStream()
	lazy.Marshal(ps)

	assert.EqualValues(t, marshalledBytes, ps.BufferBytes())
}

func TestLazyPassthrough(t *testing.T) {
	src := createLogsData()

	marshalBytes, err := gogolib.Marshal(src)
	require.NoError(t, err)
	require.NotNil(t, marshalBytes)

	ps := molecule.NewProtoStream()
	lazy := NewLogsData(marshalBytes)
	ps.Reset()
	err = lazy.Marshal(ps)
	require.NoError(t, err)
	assert.EqualValues(t, marshalBytes, ps.BufferBytes())
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
	rls := *lazy.GetResourceLogs()
	for _, rl := range rls {
		resource := *rl.GetResource()

		attrs := *resource.GetAttributes()
		attrCount += len(attrs)

		sls := *rl.GetScopeLogs()
		for _, sl := range sls {
			logRecords := *sl.GetLogRecords()

			for _, logRecord := range logRecords {
				attrs2 := *logRecord.GetAttributes()
				attrCount += len(attrs2)
			}
		}
	}
	return attrCount
}

func touchAll(lazy *LogsData) {
	rls := *lazy.GetResourceLogs()
	for _, rl := range rls {
		resource := *rl.GetResource()

		attrs := *resource.GetAttributes()
		for i := range attrs {
			attrs[i].SetKey(attrs[i].Key())
		}

		sls := *rl.GetScopeLogs()
		for _, sl := range sls {
			logRecords := *sl.GetLogRecords()

			for _, logRecord := range logRecords {
				attrs2 := *logRecord.GetAttributes()
				for i := range attrs2 {
					attrs2[i].SetKey(attrs2[i].Key())
				}
			}
		}
	}
}

func BenchmarkLazyMarshalUnchanged(b *testing.B) {
	src := createLogsData()

	marshalBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, marshalBytes)

	lazy := NewLogsData(marshalBytes)
	countAttrs(lazy)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(b, err)
		require.NotNil(b, ps.BufferBytes())
	}
}

func BenchmarkLazyMarshalFullModified(b *testing.B) {
	src := createLogsData()

	marshalBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, marshalBytes)

	lazy := NewLogsData(marshalBytes)
	countAttrs(lazy)
	touchAll(lazy)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(b, err)
		require.NotNil(b, ps.BufferBytes())
	}
}

func BenchmarkLazyUnmarshalAndReadAll(b *testing.B) {
	src := createLogsData()

	bytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, bytes)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		lazy := NewLogsData(bytes)

		// Traverse all data to get it loaded. This is the worst case.
		countAttrs(lazy)

		keyValuePool.Release(lazy)
	}
}

func BenchmarkLazyPassthroughNoReadOrModify(b *testing.B) {
	// This is the best case scenario for passthrough. We don't read or modify any
	// data, just unmarshal and marshal it exactly as it is.

	src := createLogsData()

	marshalBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, marshalBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		lazy := NewLogsData(marshalBytes)
		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(b, err)
		require.NotNil(b, ps.BufferBytes())
	}
}

func BenchmarkLazyPassthroughFullReadNoModify(b *testing.B) {
	// This is the best case scenario for passthrough. We don't read or modify any
	// data, just unmarshal and marshal it exactly as it is.

	src := createLogsData()

	marshalBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, marshalBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		lazy := NewLogsData(marshalBytes)
		ps.Reset()
		err = lazy.Marshal(ps)
		countAttrs(lazy)
		require.NoError(b, err)
		require.NotNil(b, ps.BufferBytes())
		keyValuePool.Release(lazy)
	}
}

func BenchmarkLazyPassthroughFullModified(b *testing.B) {
	// This is the worst case scenario. We read of data, so lazy loading has no
	// performance benefit. We also modify all data, so we have to do full marshaling.

	src := createLogsData()

	marshalBytes, err := gogolib.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, marshalBytes)

	b.ResetTimer()

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		lazy := NewLogsData(marshalBytes)

		// Touch all attrs
		touchAll(lazy)

		ps.Reset()
		err = lazy.Marshal(ps)
		require.NoError(b, err)
		require.NotNil(b, ps.BufferBytes())

		keyValuePool.Release(lazy)
	}
}

func BenchmarkKeyValueMarshal(b *testing.B) {
	kv := KeyValue{
		key:   "key",
		value: "val",
	}

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		ps.Reset()
		kv.Marshal(ps)
		//require.Len(b, ps.BufferBytes(), 10)
	}
}

func BenchmarkResourceMarshal(b *testing.B) {
	kv := KeyValue{
		key:   "key",
		value: "val",
	}
	res := Resource{
		attributes:             []KeyValue{kv},
		DroppedAttributesCount: 0,
	}

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		ps.Reset()
		res.Marshal(ps)
	}
}

func BenchmarkLogRecordMarshal(b *testing.B) {
	kv := KeyValue{
		key:   "key",
		value: "val",
	}
	lr := LogRecord{
		attributes: []*KeyValue{&kv},
	}

	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		ps.Reset()
		lr.Marshal(ps)
	}
}

func BenchmarkAppendVarintProtowire(b *testing.B) {
	bts := make([]byte, 0, b.N*10)
	value := uint64(123)
	for i := 0; i < b.N; i++ {
		bts = protowire.AppendVarint(bts, value)
	}
}

//func BenchmarkAppendVarintMolecule(b *testing.B) {
//	bts := make([]byte, 0, b.N*10)
//	value := uint64(123)
//	for i := 0; i < b.N; i++ {
//		bts = molecule.AppendVarint(bts, value)
//	}
//}

func BenchmarkWriteUint32(b *testing.B) {
	val := uint32(123)
	ps := molecule.NewProtoStream()
	for i := 0; i < b.N; i++ {
		ps.Reset()
		ps.Uint32(1, val)
	}
}

func BenchmarkWriteUint32Prepared(b *testing.B) {
	val := uint32(123)
	ps := molecule.NewProtoStream()
	k := molecule.PrepareUint32Field(1)
	for i := 0; i < b.N; i++ {
		ps.Reset()
		ps.Uint32Prepared(k, val)
	}
}

func BenchmarkWriteUint32LongPrepared(b *testing.B) {
	val := uint32(123)
	ps := molecule.NewProtoStream()
	k := molecule.PrepareLongField(1, protowire.VarintType)
	for i := 0; i < b.N; i++ {
		ps.Reset()
		ps.Uint32LongPrepared(k, val)
	}
}
