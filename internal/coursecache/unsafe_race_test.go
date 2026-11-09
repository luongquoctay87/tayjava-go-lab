//go:build racedemo

package coursecache

import (
	"sync"
	"testing"
)

// Run: go test -race -tags=racedemo ./internal/coursecache -run TestUnsafeCacheRace -count=1
func TestUnsafeCacheRace(t *testing.T) {
	c := NewUnsafe()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			c.Put(Course{ID: id % 4, Title: "Go", Seats: 10})
			_, _ = c.Get(id % 4)
		}(int64(i))
	}
	wg.Wait()
}
