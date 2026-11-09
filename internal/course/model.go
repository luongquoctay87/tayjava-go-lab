package course

// Course is the e-learning course record.
// Seats is total capacity (not seats_left). Enrollment count is checked against seats.
type Course struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Subtitle string `json:"subtitle,omitempty"` // additive JSON field (lesson 40)
	Seats    int    `json:"seats"`
	Level    string `json:"level"`
}

// Enrollment is a user registered in a course.
type Enrollment struct {
	ID       int64 `json:"id"`
	CourseID int64 `json:"courseId"`
	UserID   int64 `json:"userId"`
}
