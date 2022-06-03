package simple

import "sync"

type poolKeyValueType struct {
	freed []*KeyValue
	mux   sync.Mutex
}

var poolKeyValue = poolKeyValueType{}

func (p *poolKeyValueType) GetSlice(count int) []*KeyValue {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.freed) >= count {
		r := p.freed[len(p.freed)-count:]
		p.freed = p.freed[:len(p.freed)-count]
		return r
	}

	r := make([]*KeyValue, count)
	i := 0
	for ; i < len(p.freed); i++ {
		r[i] = p.freed[i]
	}
	p.freed = nil

	if i < count {
		storage := make([]KeyValue, count-i)
		j := 0
		for ; i < count; i++ {
			r[i] = &storage[j]
			j++
		}
	}

	return r
}

func (p *poolKeyValueType) ReleaseSlice(attributes []*KeyValue) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.freed = append(p.freed, attributes...)
}
