package main

import (
	"errors"
	"fmt"

	"github.com/tayjava/go-lab/course"
)

func main() {
	store := course.NewMemoryStore()
	_ = store.Save(course.New(1, "Golang from zero", 30))

	err := course.Enroll(store, 9, 42)
	if err != nil {
		fmt.Println("enroll error:", err)
		if errors.Is(err, course.ErrNotFound) {
			fmt.Println("404-ish")
		}
	}

	// Contrast: wrapping with %v loses the sentinel for errors.Is.
	bad := fmt.Errorf("enroll failed: %v", course.ErrNotFound)
	fmt.Println("Is with percent-v wrap:", errors.Is(bad, course.ErrNotFound))
	good := fmt.Errorf("enroll failed: %w", course.ErrNotFound)
	fmt.Println("Is with percent-w wrap:", errors.Is(good, course.ErrNotFound))
}
