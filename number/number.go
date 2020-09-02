package number

import "sync"

var ID = &id{}

type id struct {
	index int64
	m     sync.Mutex
}

func (i *id) Get() int64 {
	i.m.Lock()
	defer i.m.Unlock()
	i.index++
	return i.index
}
