package simple

import "sync"

type KeyValuePool struct {
	freedKeyValue []KeyValue
}

var keyValuePool = KeyValuePool{}

var poolMux = sync.Mutex{}

func (p *KeyValuePool) GetKeyValues(count int) []KeyValue {
	poolMux.Lock()
	defer poolMux.Unlock()

	if len(p.freedKeyValue) >= count {
		r := p.freedKeyValue[len(p.freedKeyValue)-count:]
		p.freedKeyValue = p.freedKeyValue[:len(p.freedKeyValue)-count]
		return r
	}

	r := make([]KeyValue, count)
	i := 0
	for ; i < len(p.freedKeyValue); i++ {
		r[i] = p.freedKeyValue[i]
	}
	p.freedKeyValue = nil

	return r
}

func (p *KeyValuePool) Release(l *LogsData) {
	poolMux.Lock()
	defer poolMux.Unlock()

	for _, rl := range l.resourceLogs {
		scopeLogsPool.freedScopeLogs = append(
			scopeLogsPool.freedScopeLogs, rl.scopeLogs...,
		)
		for _, sl := range rl.scopeLogs {
			logRecordPool.freedLogRecord = append(
				logRecordPool.freedLogRecord, sl.logRecords...,
			)
			for _, logRecord := range sl.logRecords {
				p.freedKeyValue = append(p.freedKeyValue, logRecord.attributes...)
			}
		}
	}
}
