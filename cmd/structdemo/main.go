package main

import "fmt"

// Course is the e-learning domain type for this lesson.
type Course struct {
	ID        int64
	Title     string
	Seats     int
	Published bool
}

func (c Course) RemainingSeats(enrolled int) int {
	return c.Seats - enrolled
}

func (c Course) WithTitle(title string) Course {
	c.Title = title
	return c
}

func (c *Course) Publish() {
	c.Published = true
}

func (c *Course) EnrollOne() {
	if c.Seats > 0 {
		c.Seats--
	}
}

// Remaining reports free seats after enrolled students.
func (c *Course) Remaining(enrolled int) int {
	return c.Seats - enrolled
}

// Label returns "ID:Title" for logging / display.
func (c Course) Label() string {
	return fmt.Sprintf("%d:%s", c.ID, c.Title)
}

// publishByValue is the deliberate mistake: Published on the copy only.
func (c Course) publishByValue() {
	c.Published = true
}

func main() {
	c := Course{ID: 1, Title: "Go", Seats: 30, Published: false}

	// Value-receiver Publish does not mutate the original.
	c.publishByValue()
	fmt.Println("after value Publish:", c.Published)

	c.Publish()
	fmt.Printf("after pointer Publish: %+v\n", c)
	fmt.Println("remaining(5):", c.Remaining(5))
	fmt.Println("label:", c.Label())

	demo := Course{Title: "Go", Seats: 2}
	demo.Publish()
	demo.EnrollOne()
	fmt.Printf("%+v\n", demo)
	fmt.Println(demo.WithTitle("Go Lab").Title, demo.Title)
}
