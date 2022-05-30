package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"

	"github.com/tigrannajaryan/lazyproto/internal/examples/simple/google"
)

func createLogRecord(n int) *google.LogRecord {
	sl := &google.LogRecord{
		Attributes: []*google.KeyValue{
			{
				Key:   "http.method",
				Value: "GET",
			},
			{
				Key:   "http.url",
				Value: "/checkout",
			},
		},
		DroppedAttributesCount: 12,
	}

	return sl

}

func createScopedLogs(n int) *google.ScopeLogs {
	sl := &google.ScopeLogs{}

	for i := 0; i < 10; i++ {
		sl.LogRecords = append(sl.LogRecords, createLogRecord(i))
	}

	return sl
}

func createLogsData() *google.LogsData {
	src := &google.LogsData{}

	for i := 0; i < 10; i++ {
		rl := &google.ResourceLogs{
			Resource: &google.Resource{
				Attributes: []*google.KeyValue{
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
	src := &google.LogsData{
		ResourceLogs: []*google.ResourceLogs{
			{
				Resource: &google.Resource{
					Attributes: []*google.KeyValue{
						{
							Key:   "key1",
							Value: "value1",
						},
					},
					DroppedAttributesCount: 12,
				},
				ScopeLogs: []*google.ScopeLogs{
					{
						LogRecords: []*google.LogRecord{
							{
								TimeUnixNano: 123,
								Attributes: []*google.KeyValue{
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
	bytes, err := proto.Marshal(src)
	require.NoError(t, err)

	lazy := NewLogsData(bytes)

	rl := *lazy.GetResourceLogs()
	require.Len(t, rl, 1)

	resource := *rl[0].GetResource()
	assert.EqualValues(t, 12, resource.DroppedAttributesCount)
	require.NotNil(t, resource)

	attrs := *resource.GetAttributes()
	require.Len(t, attrs, 1)

	kv1 := attrs[0]
	require.EqualValues(t, "key1", kv1.Key)
	require.EqualValues(t, "value1", kv1.Value)

	sls := *rl[0].GetScopeLogs()
	require.Len(t, sls, 1)

	sl := sls[0]
	logRecords := *sl.GetLogRecords()
	require.Len(t, logRecords, 1)

	logRecord := logRecords[0]
	assert.EqualValues(t, 123, logRecord.TimeUnixNano)
	assert.EqualValues(t, 234, logRecord.DroppedAttributesCount)
	attrs = *logRecord.GetAttributes()
	require.Len(t, attrs, 1)

	kv2 := attrs[0]
	require.EqualValues(t, "key2", kv2.Key)
	require.EqualValues(t, "value2", kv2.Value)
}

func BenchmarkGoogleMarshal(b *testing.B) {
	src := createLogsData()

	for i := 0; i < b.N; i++ {
		bytes, err := proto.Marshal(src)
		require.NoError(b, err)
		require.NotNil(b, bytes)
	}
}

func BenchmarkGoogleUnmarshal(b *testing.B) {
	src := createLogsData()

	bytes, err := proto.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, bytes)

	for i := 0; i < b.N; i++ {
		var ld google.LogsData
		err := proto.Unmarshal(bytes, &ld)
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

		require.EqualValues(b, 2010, attrCount)
	}
}

func BenchmarkLazyUnmarshal(b *testing.B) {
	src := createLogsData()

	bytes, err := proto.Marshal(src)
	require.NoError(b, err)
	require.NotNil(b, bytes)

	for i := 0; i < b.N; i++ {
		lazy := NewLogsData(bytes)

		// Traverse all data. This is the worst case.
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
					attrs = *logRecord.GetAttributes()
					attrCount += len(attrs)
				}
			}
		}
		require.EqualValues(b, 2010, attrCount)
	}
}
