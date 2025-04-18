package store

import "sync"

type GarnetStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func New() *GarnetStore {
	return &GarnetStore{
		data: make(map[string]string),
	}
}

func (s *GarnetStore) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *GarnetStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func (s *GarnetStore) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.data[key]; exists {
		delete(s.data, key)
		return true
	}
	return false
}

func (s *GarnetStore) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[key]
	return exists
}

func (s *GarnetStore) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *GarnetStore) FlushAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]string)
}
