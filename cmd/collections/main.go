package main

import "fmt"

type Course struct {
	ID    int64
	Title string
	Seats int
}

type Enrollment struct {
	UserID   int64
	CourseID int64
	Seats    int
}

func demoNilMapPanic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("nil map write panicked (expected):", r)
		}
	}()
	var bad map[int64]Course
	bad[1] = Course{ID: 1, Title: "panic"} // assignment to entry in nil map
}

func main() {
	titles := []string{"Golang", "Java", "SQL"}
	fmt.Println(len(titles), titles[0])

	ids := make([]int64, 0, 8)
	ids = append(ids, 101, 102)
	fmt.Println(ids, len(ids), cap(ids))

	var a []int
	b := []int{}
	fmt.Println("nil slice:", a == nil, "empty slice:", b == nil)

	courses := []Course{
		{ID: 1, Title: "Go", Seats: 10},
		{ID: 2, Title: "Java", Seats: 15},
	}
	courses = append(courses, Course{ID: 3, Title: "SQL", Seats: 20})
	fmt.Println("len:", len(courses), "cap:", cap(courses))

	index := map[int64]Course{}
	for _, c := range courses {
		index[c.ID] = c
	}
	if c, ok := index[2]; ok {
		fmt.Printf("lookup 2: %+v ok=%v\n", c, ok)
	}

	demoNilMapPanic()
	fixed := make(map[int64]Course)
	fixed[1] = Course{ID: 1, Title: "Go", Seats: 10}
	fmt.Println("fixed map len:", len(fixed))

	byUser := map[int64]Enrollment{
		7: {UserID: 7, CourseID: 101, Seats: 1},
	}
	if e, ok := byUser[7]; ok {
		fmt.Println("enrolled course:", e.CourseID)
	}
	if _, ok := byUser[8]; !ok {
		fmt.Println("not enrolled")
	}
	delete(byUser, 7)
}
