package simple

import "sync"

// Pool of ResourceLogs structs.
type resourceLogsPoolType struct {
	pool []*ResourceLogs
	mux  sync.Mutex
}

var resourceLogsPool = resourceLogsPoolType{}

// Get one element from the pool. Creates a new element if the pool is empty.
func (p *resourceLogsPoolType) Get() *ResourceLogs {
	p.mux.Lock()
	defer p.mux.Unlock()

	// Have elements in the pool?
	if len(p.pool) >= 1 {
		// Get the last element.
		r := p.pool[len(p.pool)-1]
		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-1]
		return r
	}

	// Pool is empty, create a new element.
	return &ResourceLogs{}
}

func (p *resourceLogsPoolType) GetSlice(count int) []*ResourceLogs {
	p.mux.Lock()
	defer p.mux.Unlock()

	// Have enough elements in the pool?
	if len(p.pool) >= count {
		// Cut the required slice from the end of the pool.
		r := p.pool[len(p.pool)-count:]
		// Shrink the pool.
		p.pool = p.pool[:len(p.pool)-count]
		return r
	}

	// Create a new slice.
	r := make([]*ResourceLogs, count)

	// Initialize with what remains in the pool.
	i := 0
	for ; i < len(p.pool); i++ {
		r[i] = p.pool[i]
	}
	p.pool = nil

	if i < count {
		// Create remaining elements.
		storage := make([]ResourceLogs, count-i)
		j := 0
		for ; i < count; i++ {
			r[i] = &storage[j]
			j++
		}
	}

	return r
}

// ReleaseSlice releases a slice of elements back to the pool.
func (p *resourceLogsPoolType) ReleaseSlice(slice []*ResourceLogs) {
	// Release nested messages recursively to their respective pools.
	for _, elem := range slice {
		if elem.resource != nil {
			resourcePool.Release(elem.resource)
		}
		scopeLogsPool.ReleaseSlice(elem.scopeLogs)

		// Zero-initialize the released element.
		*elem = ResourceLogs{}
	}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the slice to the end of the pool.
	p.pool = append(p.pool, slice...)
}

// Release an element back to the pool.
func (p *resourceLogsPoolType) Release(elem *ResourceLogs) {
	// Release nested messages recursively to their respective pools.
	if elem.resource != nil {
		resourcePool.Release(elem.resource)
	}
	scopeLogsPool.ReleaseSlice(elem.scopeLogs)

	// Zero-initialize the released element.
	*elem = ResourceLogs{}

	p.mux.Lock()
	defer p.mux.Unlock()

	// Add the element to the end of the pool.
	p.pool = append(p.pool, elem)
}
