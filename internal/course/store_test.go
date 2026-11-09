package course_test

import (
	"path/filepath"
	"testing"

	"github.com/tayjava/go-lab/internal/course"
)

func TestAddDuplicate(t *testing.T) {
	dir := t.TempDir()
	store := course.Store{Path: filepath.Join(dir, "courses.json")}

	c := course.Course{ID: 1, Title: "Go", Seats: 10, Level: "beginner"}
	if err := store.Add(c); err != nil {
		t.Fatalf("first add: %v", err)
	}
	if err := store.Add(c); err == nil {
		t.Fatal("expected duplicate id error")
	}
}

func TestLoadMissingFile(t *testing.T) {
	store := course.Store{Path: filepath.Join(t.TempDir(), "missing.json")}
	list, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("want empty list, got %d", len(list))
	}
}
