package simple

import "sync"

type resourceLogsPoolType struct {
	freed []*ResourceLogs
	mux   sync.Mutex
}

var resourceLogsPool = resourceLogsPoolType{}

func (p *resourceLogsPoolType) Get() *ResourceLogs {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.freed) >= 1 {
		r := p.freed[len(p.freed)-1]
		p.freed = p.freed[:len(p.freed)-1]
		return r
	}

	return &ResourceLogs{}
}

func (p *resourceLogsPoolType) GetSlice(count int) []*ResourceLogs {
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

func (p *resourceLogsPoolType) ReleaseSlice(slice []*ResourceLogs) {
	for _, elem := range slice {
		if elem.resource != nil {
			resourcePool.Release(elem.resource)
		}
		scopeLogsPool.ReleaseSlice(elem.scopeLogs)
		*elem = ResourceLogs{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	p.freed = append(p.freed, slice...)
}

func (p *resourceLogsPoolType) Release(elem *ResourceLogs) {
	if elem.resource != nil {
		resourcePool.Release(elem.resource)
	}
	scopeLogsPool.ReleaseSlice(elem.scopeLogs)
	*elem = ResourceLogs{}

	p.mux.Lock()
	defer p.mux.Unlock()

	p.freed = append(p.freed, elem)
}
