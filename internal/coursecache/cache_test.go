package coursecache

import (
	"sync"
	"testing"
)

func TestSafeCacheConcurrent(t *testing.T) {
	c := NewSafe()
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			c.Put(Course{ID: id, Title: "Go", Seats: 30})
			_, _ = c.Get(id)
		}(int64(i % 8))
	}
	wg.Wait()
}

func TestSafeReserveSeat(t *testing.T) {
	c := NewSafe()
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
	got, ok := c.Get(1)
	if !ok || got.Seats != 0 {
		t.Fatalf("seats=%d ok=%v want 0", got.Seats, ok)
	}
}
