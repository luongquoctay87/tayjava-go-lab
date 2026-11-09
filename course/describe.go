package course

import (
	"fmt"
	"strings"
)

// DescribeSlow builds a label with repeated string concat (allocation-heavy).
func DescribeSlow(c Course, enrolled int) string {
	s := "course="
	s += fmt.Sprintf("%d", c.ID)
	s += " title="
	s += c.Title
	s += " seats="
	s += fmt.Sprintf("%d", c.Seats)
	s += " enrolled="
	s += fmt.Sprintf("%d", enrolled)
	s += " remaining="
	s += fmt.Sprintf("%d", c.Seats-enrolled)
	return s
}

// DescribeFast builds the same label with a strings.Builder.
func DescribeFast(c Course, enrolled int) string {
	var b strings.Builder
	b.Grow(64)
	fmt.Fprintf(&b, "course=%d title=%s seats=%d enrolled=%d remaining=%d",
		c.ID, c.Title, c.Seats, enrolled, c.Seats-enrolled)
	return b.String()
}
