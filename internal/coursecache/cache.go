package coursecache

import "sync"

// Course is a tiny cache value for the race-detector lab.
type Course struct {
	ID    int64
	Title string
	Seats int
}

// UnsafeCache races on concurrent Get/Put (deliberate bug for lesson 30).
type UnsafeCache struct {
	byID map[int64]*Course
}

// NewUnsafe builds an empty unsafe cache.
func NewUnsafe() *UnsafeCache {
	return &UnsafeCache{byID: make(map[int64]*Course)}
}

// Put stores c without synchronization.
func (c *UnsafeCache) Put(course Course) {
	c.byID[course.ID] = &course
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
