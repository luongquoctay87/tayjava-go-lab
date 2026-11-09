package main

import (
	"fmt"
)

type Course struct {
	ID    int64
	Title string
	Seats int
}

type CourseStore interface {
	FindByID(id int64) (Course, error)
	Save(c Course) error
}

type MemoryCourseStore struct {
	items map[int64]Course
}

func NewMemoryCourseStore() *MemoryCourseStore {
	return &MemoryCourseStore{items: make(map[int64]Course)}
}

func (s *MemoryCourseStore) FindByID(id int64) (Course, error) {
	c, ok := s.items[id]
	if !ok {
		return Course{}, fmt.Errorf("course %d not found", id)
	}
	return c, nil
}

func (s *MemoryCourseStore) Save(c Course) error {
	s.items[c.ID] = c
	return nil
}

// StubCourseStore hard-codes one course for interface swapping demos.
type StubCourseStore struct{}

func (StubCourseStore) FindByID(id int64) (Course, error) {
	if id != 99 {
		return Course{}, fmt.Errorf("course %d not found", id)
	}
	return Course{ID: 99, Title: "Stub Golang", Seats: 5}, nil
}

func (StubCourseStore) Save(c Course) error {
	return nil
}

type EnrollmentService struct {
	store CourseStore
}

func (s EnrollmentService) EnrollTitle(id int64) (string, error) {
	c, err := s.store.FindByID(id)
	if err != nil {
		return "", err
	}
	return c.Title, nil
}

func main() {
	mem := NewMemoryCourseStore()
	_ = mem.Save(Course{ID: 1, Title: "Golang", Seats: 20})

	svc := EnrollmentService{store: mem}
	title, err := svc.EnrollTitle(1)
	if err != nil {
		fmt.Println("memory error:", err)
	} else {
		fmt.Println(title)
	}

	// Swap store without changing EnrollmentService.
	svc.store = StubCourseStore{}
	title, err = svc.EnrollTitle(99)
	if err != nil {
		fmt.Println("stub error:", err)
	} else {
		fmt.Println(title)
	}
}
