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
