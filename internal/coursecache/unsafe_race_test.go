//go:build racedemo

package coursecache

import (
	"sync"
	"testing"
)

// Run: go test -race -tags=racedemo ./internal/coursecache -run TestUnsafeCacheRace -count=1
// The map is only written before the goroutines start. Concurrent map reads are
// not the race; every ReserveSeat writes Seats with no happens-before edge.
func TestUnsafeCacheRace(t *testing.T) {
	c := NewUnsafe()
	c.Put(Course{ID: 1, Title: "Go", Seats: 100})
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.ReserveSeat(1)
		}()
	}
	wg.Wait()
}
