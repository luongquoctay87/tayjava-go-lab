package resilience

import "sync"

// IdempotencyStore maps a client key to a prior enrollment id (in-memory lab store).
type IdempotencyStore struct {
	mu   sync.Mutex
	done map[string]int64
}

// NewIdempotencyStore creates an empty store.
func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{done: make(map[string]int64)}
}

// Get returns a prior enrollment id when the key was already processed.
func (s *IdempotencyStore) Get(key string) (int64, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, ok := s.done[key]
	return id, ok
}

// Put records the enrollment id for key (first writer wins).
func (s *IdempotencyStore) Put(key string, enrollmentID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.done[key]; ok {
		return
	}
	s.done[key] = enrollmentID
}
