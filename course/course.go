package course

// Course is the shared e-learning domain type.
// Unexported field level is only reachable via Level().
type Course struct {
	ID    int64
	Title string
	Seats int
	level string
}

// New constructs a Course with default level "beginner".
func New(id int64, title string, seats int) Course {
	return Course{ID: id, Title: title, Seats: seats, level: "beginner"}
}

// Level returns the unexported difficulty level.
func (c Course) Level() string {
	return c.level
}
