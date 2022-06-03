package simple

import "sync"

type LogsDataPool struct {
	freed []*LogsData
	mux   sync.Mutex
}

var logsDataPool = LogsDataPool{}

func (p *LogsDataPool) Get() *LogsData {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.freed) >= 1 {
		r := p.freed[len(p.freed)-1]
		p.freed = p.freed[:len(p.freed)-1]
		return r
	}

	return &LogsData{}
}

func (p *LogsDataPool) GetSlice(count int) []*LogsData {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.freed) >= count {
		r := p.freed[len(p.freed)-count:]
		p.freed = p.freed[:len(p.freed)-count]
		return r
	}

	r := make([]*LogsData, count)
	i := 0
	for ; i < len(p.freed); i++ {
		r[i] = p.freed[i]
	}
	p.freed = nil
	if i < count {
		storage := make([]LogsData, count-i)
		j := 0
		for ; i < count; i++ {
			r[i] = &storage[j]
			j++
		}
	}

	return r
}

func (p *LogsDataPool) ReleaseSlice(rls []*LogsData) {
	for _, rl := range rls {
		resourceLogsPool.ReleaseSlice(rl.resourceLogs)
	}

	p.mux.Lock()
	defer p.mux.Unlock()
	logsDataPool.freed = append(logsDataPool.freed, rls...)
}

func (p *LogsDataPool) Release(rls *LogsData) {
	resourceLogsPool.ReleaseSlice(rls.resourceLogs)

	p.mux.Lock()
	defer p.mux.Unlock()
	logsDataPool.freed = append(logsDataPool.freed, rls)
}
