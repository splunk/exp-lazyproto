package simple

import "sync"

type ResourcePool struct {
	freed []*Resource
	mux   sync.Mutex
}

var resourcePool = ResourcePool{}

func (p *ResourcePool) Get() *Resource {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.freed) >= 1 {
		r := p.freed[len(p.freed)-1]
		p.freed = p.freed[:len(p.freed)-1]
		return r
	}

	return &Resource{}
}

func (p *ResourcePool) GetSlice(count int) []*Resource {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.freed) >= count {
		r := p.freed[len(p.freed)-count:]
		p.freed = p.freed[:len(p.freed)-count]
		return r
	}

	r := make([]*Resource, count)
	i := 0
	for ; i < len(p.freed); i++ {
		r[i] = p.freed[i]
	}
	p.freed = nil
	for ; i < count; i++ {
		r[i] = &Resource{}
	}

	return r
}

func (p *ResourcePool) ReleaseSlice(records []*Resource) {
	for _, resource := range records {
		keyValuePool.ReleaseSlice(resource.attributes)
	}

	p.mux.Lock()
	defer p.mux.Unlock()
	p.freed = append(p.freed, records...)
}

func (p *ResourcePool) Release(resource *Resource) {
	keyValuePool.ReleaseSlice(resource.attributes)

	p.mux.Lock()
	defer p.mux.Unlock()

	p.freed = append(p.freed, resource)
}
