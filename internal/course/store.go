package course

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Store loads and saves a Course list as indented JSON.
type Store struct {
	Path string
}

// Load returns an empty list when the file is missing or empty.
func (s Store) Load() ([]Course, error) {
	data, err := os.ReadFile(s.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Course{}, nil
		}
		return nil, err
	}
	if len(data) == 0 {
		return []Course{}, nil
	}
	var list []Course
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("parse courses: %w", err)
	}
	return list, nil
}

// Save writes the course list with indentation.
func (s Store) Save(list []Course) error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(s.Path, data, 0o644)
}

// Add appends a course or returns an error when the ID already exists.
func (s Store) Add(c Course) error {
	if c.Seats < 0 {
		return fmt.Errorf("seats must not be negative")
	}
	list, err := s.Load()
	if err != nil {
		return err
	}
	for _, existing := range list {
		if existing.ID == c.ID {
			return fmt.Errorf("course id %d already exists", c.ID)
		}
	}
	list = append(list, c)
	return s.Save(list)
}

// Get returns a pointer to the course with the given ID.
func (s Store) Get(id int64) (*Course, error) {
	list, err := s.Load()
	if err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].ID == id {
			return &list[i], nil
		}
	}
	return nil, fmt.Errorf("course %d not found", id)
}
