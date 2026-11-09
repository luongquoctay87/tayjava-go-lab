package course

import "fmt"

// Reserve decrements Seats by one. Requires a non-nil course with seats left.
func Reserve(c *Course) error {
	if c == nil {
		return fmt.Errorf("course is nil")
	}
	if c.Seats <= 0 {
		return fmt.Errorf("no seats left")
	}
	c.Seats--
	return nil
}

// ReserveCopy takes a Course by value — the caller's Seats stay unchanged.
func ReserveCopy(c Course) {
	c.Seats--
}
