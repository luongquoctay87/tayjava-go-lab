package main

import "fmt"

func main() {
	// Exercise: declare with var and short declaration
	var courseTitle string
	seats := 25
	courseTitle = "Golang from zero"
	fmt.Println("course:", courseTitle, "seats:", seats)

	// Zero value of an unassigned float64
	var price float64
	fmt.Printf("price zero value: %v\n", price)

	// Explicit conversion between numeric types
	var enrolled int32 = 10
	var capacity int64 = 30
	remaining := capacity - int64(enrolled)
	fmt.Println("remaining seats:", remaining)

	// Zero values for common types
	var zSeats int
	var zTitle string
	var published bool
	var course *int
	fmt.Printf("seats=%d title=%q published=%v course=%v\n",
		zSeats, zTitle, published, course)
}
