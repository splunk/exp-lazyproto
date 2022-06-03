package simple

import "sync"

type ResourceLogsPool struct {
	freed []*ResourceLogs
	mux   sync.Mutex
}

var resourceLogsPool = ResourceLogsPool{}

func (p *ResourceLogsPool) Get(count int) []*ResourceLogs {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.freed) >= count {
		r := p.freed[len(p.freed)-count:]
		p.freed = p.freed[:len(p.freed)-count]
		return r
	}

	r := make([]*ResourceLogs, count)
	i := 0
	for ; i < len(p.freed); i++ {
		r[i] = p.freed[i]
	}
	p.freed = nil
	if i < count {
		storage := make([]ResourceLogs, count-i)
		j := 0
		for ; i < count; i++ {
			r[i] = &storage[j]
			j++
		}
	}

	return r
}

func (p *ResourceLogsPool) ReleaseSlice(rls []*ResourceLogs) {
	for _, rl := range rls {
		if rl.resource != nil {
			resourcePool.Release(rl.resource)
		}
		scopeLogsPool.ReleaseSlice(rl.scopeLogs)
	}

	p.mux.Lock()
	defer p.mux.Unlock()
	resourceLogsPool.freed = append(resourceLogsPool.freed, rls...)
}
