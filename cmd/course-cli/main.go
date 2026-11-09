package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tayjava/go-lab/internal/course"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	store := course.Store{Path: "data/courses.json"}
	cmd := os.Args[1]

	switch cmd {
	case "add":
		fs := flag.NewFlagSet("add", flag.ExitOnError)
		id := fs.Int64("id", 0, "course id")
		title := fs.String("title", "", "course title")
		seats := fs.Int("seats", 0, "seat capacity")
		level := fs.String("level", "beginner", "difficulty level")
		_ = fs.Parse(os.Args[2:])
		if *id == 0 || *title == "" {
			usage()
			os.Exit(2)
		}
		must(store.Add(course.Course{
			ID:    *id,
			Title: *title,
			Seats: *seats,
			Level: *level,
		}))
		fmt.Printf("added course %d\n", *id)

	case "list":
		list, err := store.Load()
		must(err)
		for _, c := range list {
			fmt.Printf("%d\t%s\tseats=%d\t%s\n", c.ID, c.Title, c.Seats, c.Level)
		}

	case "get":
		fs := flag.NewFlagSet("get", flag.ExitOnError)
		id := fs.Int64("id", 0, "course id")
		_ = fs.Parse(os.Args[2:])
		if *id == 0 {
			usage()
			os.Exit(2)
		}
		c, err := store.Get(*id)
		must(err)
		fmt.Printf("%+v\n", c)

	default:
		usage()
		os.Exit(2)
	}
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, `usage:
  course-cli add -id 1 -title "Go Beginner" -seats 30 -level beginner
  course-cli list
  course-cli get -id 1`)
}
