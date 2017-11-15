package mockstorage

import (
	"sync"
)

type MockStorage struct {
	mu    sync.RWMutex
	store map[string]interface{}
}

func NewMockStorage() *MockStorage {
	return &MockStorage{store: make(map[string]interface{})}
}

func (m *MockStorage) Store(id string, data interface{}) {
	m.mu.Lock()
	m.store[id] = data
	m.mu.Unlock()
}

func (m *MockStorage) Load(id string) interface{} {
	m.mu.RLock()
	data, found := m.store[id]
	m.mu.RUnlock()
	if found {
		return data
	}

	return nil
}

func (m *MockStorage) List() []interface{} {
	list := []interface{}{}
	m.mu.RLock()
	for _, v := range m.store {
		list = append(list, v)
	}
	m.mu.RUnlock()

	return list
}

func (m *MockStorage) Clear() {
	m.mu.Lock()
	m.store = make(map[string]interface{})
	m.mu.Unlock()
}
