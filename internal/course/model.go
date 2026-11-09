package course

// Course is the JSON-persisted e-learning course record for the CLI lab.
type Course struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Seats int    `json:"seats"`
	Level string `json:"level"`
}
