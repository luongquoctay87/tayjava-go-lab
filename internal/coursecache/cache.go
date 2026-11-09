package coursecache

import "sync"

// Course is a tiny cache value for the race-detector lab.
type Course struct {
	ID    int64
	Title string
	Seats int
}

// UnsafeCache is a deliberate lesson-30 bug. Concurrent map reads alone are
// not a data race. The race is an unsynchronized write: Put writes the map,
// and ReserveSeat writes the Seats field.
type UnsafeCache struct {
	byID map[int64]*Course
}

// NewUnsafe builds an empty unsafe cache.
func NewUnsafe() *UnsafeCache {
	return &UnsafeCache{byID: make(map[int64]*Course)}
}

// Put stores c without synchronization. A concurrent Put is a map write race.
// Concurrent Get calls with no writer are not a race by themselves.
func (c *UnsafeCache) Put(course Course) {
	c.byID[course.ID] = &course
}

// ReserveSeat decrements Seats with no lock.
// Concurrent map reads alone are safe; the race is this field write
// (and any map write, such as Put).
func (c *UnsafeCache) ReserveSeat(id int64) {
	c.byID[id].Seats--
}

// Get loads by id without synchronization.
func (c *UnsafeCache) Get(id int64) (Course, bool) {
	v, ok := c.byID[id]
	if !ok {
		return Course{}, false
	}
	return *v, true
}

// SafeCache protects the map with a mutex.
type SafeCache struct {
	mu   sync.RWMutex
	byID map[int64]*Course
}

// NewSafe builds an empty mutex-backed cache.
func NewSafe() *SafeCache {
	return &SafeCache{byID: make(map[int64]*Course)}
}

// Put stores c under a write lock.
func (c *SafeCache) Put(course Course) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cp := course
	c.byID[course.ID] = &cp
}

// ReserveSeat decrements Seats under the write lock.
func (c *SafeCache) ReserveSeat(id int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.byID[id].Seats--
}

// Get loads by id under a read lock.
func (c *SafeCache) Get(id int64) (Course, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.byID[id]
	if !ok {
		return Course{}, false
	}
	return *v, true
}
