package atomicutil

import (
	"sync"
	"sync/atomic"
)

// TODO: make other types generateable with https://github.com/scritchley/gen

// NewMapStringUint64 creates a MapStringUint64 with the map initialised
func NewMapStringUint64() MapStringUint64 {
	return MapStringUint64{m: make(map[string]*uint64)}
}

// MapStringUint64 implements a concurrent safe map[string]uint64
// zero value is not usable, this type must be initialised via NewMapUint64()
type MapStringUint64 struct {
	mu sync.RWMutex
	m  map[string]*uint64
}

// Get a key from the map, the second return value indicates whether or not a key was found
func (m *MapStringUint64) Get(k string) (uint64, bool) {
	m.mu.RLock()
	vPtr := m.m[k]
	if vPtr != nil {
		v := atomic.LoadUint64(vPtr)
		m.mu.RUnlock()
		return v, true
	}
	m.mu.RUnlock()
	return 0, false
}

// Set a key in the map
func (m *MapStringUint64) Set(k string, v uint64) {
	m.mu.RLock()
	if _, ok := m.m[k]; ok {
		atomic.StoreUint64(m.m[k], v)
	} else {
		m.mu.RUnlock()
		m.mu.Lock()
		m.m[k] = &v
		m.mu.Unlock()
		return
	}
	m.mu.RUnlock()
}

// IncN increments a key in the map by the given delta, creating and setting to 1 if it
// does not already exist, and returns the incremented value
func (m *MapStringUint64) IncN(k string, delta uint64) uint64 {
	m.mu.RLock()
	var v uint64
	if _, ok := m.m[k]; ok {
		v = atomic.AddUint64(m.m[k], delta)
	} else {
		m.mu.RUnlock()
		m.mu.Lock()
		v := uint64(1)
		m.m[k] = &v
		m.mu.Unlock()
		return v
	}
	m.mu.RUnlock()
	return v
}

// Inc is equivalent to IncN(k, 1)
func (m *MapStringUint64) Inc(k string) uint64 {
	return m.IncN(k, 1)
}

// Reset the map, reinitialising it
func (m *MapStringUint64) Reset() {
	m.mu.Lock()
	m.m = make(map[string]*uint64)
	m.mu.Unlock()
}

// Len returns the number of keys in the map
func (m *MapStringUint64) Len() int {
	m.mu.RLock()
	l := len(m.m)
	m.mu.RUnlock()
	return l
}
