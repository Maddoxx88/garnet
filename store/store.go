package store

import (
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time // zero = no expiry
}

type GarnetStore struct {
	mu   sync.RWMutex
	data map[string]entry
}

func New() *GarnetStore {
	return &GarnetStore{
		data: make(map[string]entry),
	}
}

func (s *GarnetStore) Set(key, value string, ttlSeconds int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expires time.Time
	if ttlSeconds > 0 {
		expires = time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	}

	s.data[key] = entry{
		value:     value,
		expiresAt: expires,
	}
}

func (s *GarnetStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	e, ok := s.data[key]
	if !ok {
		return "", false
	}

	if !e.expiresAt.IsZero() && time.Now().After(e.expiresAt) {
		return "", false // expired
	}

	return e.value, true
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
	s.data = make(map[string]entry)
}
