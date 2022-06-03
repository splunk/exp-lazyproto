package simple

import "sync"

type ScopeLogsPool struct {
	freed []*ScopeLogs
	mux   sync.Mutex
}

var scopeLogsPool = ScopeLogsPool{}

func (p *ScopeLogsPool) Get(count int) []*ScopeLogs {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.freed) >= count {
		r := p.freed[len(p.freed)-count:]
		p.freed = p.freed[:len(p.freed)-count]
		return r
	}

	r := make([]*ScopeLogs, count)
	i := 0
	for ; i < len(p.freed); i++ {
		r[i] = p.freed[i]
	}
	p.freed = nil
	if i < count {
		storage := make([]ScopeLogs, count-i)
		j := 0
		for ; i < count; i++ {
			r[i] = &storage[j]
			j++
		}
	}

	return r
}

func (p *ScopeLogsPool) ReleaseSlice(sls []*ScopeLogs) {
	for _, sl := range sls {
		logRecordPool.ReleaseSlice(sl.logRecords)
	}

	p.mux.Lock()
	defer p.mux.Unlock()
	scopeLogsPool.freed = append(scopeLogsPool.freed, sls...)
}
