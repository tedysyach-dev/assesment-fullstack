package marketplace

import (
	"context"
	"sync"
)

type SessionStore interface {
	Get(ctx context.Context, key string) (*Session, error)
	Set(ctx context.Context, key string, session *Session) error
	Delete(ctx context.Context, key string) error
}

type MemorySessionStore struct {
	data map[string]*Session
	mu   sync.RWMutex
}

func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{
		data: map[string]*Session{},
	}
}

func (m *MemorySessionStore) Get(ctx context.Context, key string) (*Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.data[key], nil
}

func (m *MemorySessionStore) Set(ctx context.Context, key string, session *Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = session
	return nil
}

func (m *MemorySessionStore) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.data, key)
	return nil
}
