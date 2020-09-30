package kv

import "sync"

var (
	inMemory *InMemory
)

type InMemory struct {
	data map[string]string
	lock sync.RWMutex
}

func Memory() *InMemory {
	return inMemory
}

// SetString set a key value
func (i *InMemory) SetString(key, value string) {
	i.lock.Lock()
	defer i.lock.Unlock()

	i.data[key] = value
}

// Get get a key value
func (i *InMemory) Get(key string) (string, bool) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	r, ok := i.data[key]
	return r, ok
}

func init() {
	inMemory = new(InMemory)
	inMemory.data = make(map[string]string)
}
