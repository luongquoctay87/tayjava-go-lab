package seatledger

import "sync"

// SeatLedger tracks remaining seats per course with a Mutex.
type SeatLedger struct {
	mu    sync.Mutex
	seats map[int64]int
}

// New creates a ledger seeded with the given course → seat counts.
func New(seats map[int64]int) *SeatLedger {
	copied := make(map[int64]int, len(seats))
	for id, n := range seats {
		copied[id] = n
	}
	return &SeatLedger{seats: copied}
}

// Reserve decrements one seat when available. Returns false when none left.
func (l *SeatLedger) Reserve(courseID int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	n, ok := l.seats[courseID]
	if !ok || n <= 0 {
		return false
	}
	l.seats[courseID] = n - 1
	return true
}

// Remaining returns seats left for a course (safe after WaitGroup joins).
func (l *SeatLedger) Remaining(courseID int64) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.seats[courseID]
}
