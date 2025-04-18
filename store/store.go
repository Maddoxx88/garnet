package store

import "sync"

type GarnetStore struct {
    mu    sync.RWMutex
    data  map[string]string
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
