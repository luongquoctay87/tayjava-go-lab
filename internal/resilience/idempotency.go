package resilience

import (
	"crypto/sha256"
	"errors"
	"strconv"
	"sync"
)

var (
	// ErrKeyReused means the same user sent the same key with a different body.
	ErrKeyReused = errors.New("idempotency key reused with a different body")
	// ErrInProgress means a request with this key is still running.
	ErrInProgress = errors.New("request with this key is still processing")
)

type entry struct {
	fingerprint [32]byte
	done        bool
	status      int
	body        []byte
}

// Idempotency is an in-memory lab store. Keys are scoped per user.
// The map is not shared across processes and is lost on restart.
type Idempotency struct {
	mu      sync.Mutex
	entries map[string]entry // key = userID + ":" + Idempotency-Key
}

// NewIdempotency creates an empty store.
func NewIdempotency() *Idempotency {
	return &Idempotency{entries: make(map[string]entry)}
}

// Do runs fn at most once per (userID, key) and returns the first status and body.
// The lock only guards the map: fn (seat reservation and other I/O) runs outside it.
// A 5xx result is not stored, so the next retry may run fn again.
func (s *Idempotency) Do(userID int64, key string, reqBody []byte, fn func() (int, []byte)) (int, []byte, error) {
	k := strconv.FormatInt(userID, 10) + ":" + key
	fp := sha256.Sum256(reqBody)

	s.mu.Lock()
	e, seen := s.entries[k]
	if !seen {
		s.entries[k] = entry{fingerprint: fp} // reserve the slot, state "processing"
	}
	s.mu.Unlock()

	if seen {
		switch {
		case e.fingerprint != fp:
			return 0, nil, ErrKeyReused
		case !e.done:
			return 0, nil, ErrInProgress
		}
		return e.status, e.body, nil
	}

	status, body := fn()

	s.mu.Lock()
	if status >= 500 {
		delete(s.entries, k)
	} else {
		stored := append([]byte(nil), body...)
		s.entries[k] = entry{fingerprint: fp, done: true, status: status, body: stored}
	}
	s.mu.Unlock()
	return status, body, nil
}
