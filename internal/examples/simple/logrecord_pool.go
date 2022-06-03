package simple

import "sync"

type LogRecordPool struct {
	freed []*LogRecord
	mux   sync.Mutex
}

var logRecordPool = LogRecordPool{}

func (p *LogRecordPool) Get(count int) []*LogRecord {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.freed) >= count {
		r := p.freed[len(p.freed)-count:]
		p.freed = p.freed[:len(p.freed)-count]
		return r
	}

	r := make([]*LogRecord, count)
	i := 0
	for ; i < len(p.freed); i++ {
		r[i] = p.freed[i]
	}
	p.freed = nil
	if i < count {
		storage := make([]LogRecord, count-i)
		j := 0
		for ; i < count; i++ {
			r[i] = &storage[j]
			j++
		}
	}

	return r
}

func (p *LogRecordPool) ReleaseSlice(records []*LogRecord) {
	for _, logRecord := range records {
		poolKeyValue.ReleaseSlice(logRecord.attributes)
	}

	p.mux.Lock()
	defer p.mux.Unlock()
	p.freed = append(p.freed, records...)
}
