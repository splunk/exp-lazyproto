package simple

type LogRecordPool struct {
	freedLogRecord []LogRecord
}

var logRecordPool = LogRecordPool{}

func (p *LogRecordPool) GetLogRecords(count int) []LogRecord {
	poolMux.Lock()
	defer poolMux.Unlock()

	if len(p.freedLogRecord) >= count {
		r := p.freedLogRecord[len(p.freedLogRecord)-count:]
		p.freedLogRecord = p.freedLogRecord[:len(p.freedLogRecord)-count]
		return r
	}

	r := make([]LogRecord, count)
	i := 0
	for ; i < len(p.freedLogRecord); i++ {
		r[i] = p.freedLogRecord[i]
	}
	p.freedLogRecord = nil

	return r
}

//func (p *LogRecordPool) Release(l *LogsData) {
//	for _, rl := range l.resourceLogs {
//		for _, sl := range rl.scopeLogs {
//			p.freedLogRecord = append(p.freedLogRecord, sl.logRecords...)
//		}
//	}
//}
