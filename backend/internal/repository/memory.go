package repository

import "sync"

// Memory is a thread-safe in-memory key-value store.
type Memory struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewMemory returns an empty in-memory store.
func NewMemory() *Memory {
	return &Memory{data: make(map[string]string)}
}

// Set stores a value under key.
func (m *Memory) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// Get returns the value for key and whether it existed.
func (m *Memory) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[key]
	return v, ok
}
