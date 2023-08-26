package csv

import (
	"sync"
)

type SyncMap[K comparable, V any] struct {
	sync.RWMutex
	internal map[K]V
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		internal: make(map[K]V),
	}
}

func (rm *SyncMap[K, V]) Load(key K) (value V, ok bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()
	return result, ok
}
func (rm *SyncMap[K, V]) LoadAll() (value []V) {
	rm.RLock()
	vals := []V{}
	for _, v := range rm.internal {
		vals = append(vals, v)
	}
	rm.RUnlock()
	return vals
}

func (rm *SyncMap[K, V]) Delete(key K) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *SyncMap[K, V]) Store(key K, value V) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}
