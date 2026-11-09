package main

import (
	"errors"
	"fmt"
)

func remainingSeats(capacity, enrolled int) int {
	return capacity - enrolled
}

func enroll(seats, request int) (int, error) {
	if request <= 0 {
		return 0, errors.New("request must be positive")
	}
	if request > seats {
		return seats, errors.New("not enough seats")
	}
	return seats - request, nil
}

// canEnroll reports whether a seat request can be fulfilled.
func canEnroll(seats, request int) (bool, error) {
	if request <= 0 {
		return false, errors.New("request must be positive")
	}
	if request > seats {
		return false, errors.New("not enough seats")
	}
	return true, nil
}

func bump(seats int) {
	seats++
}

func printEnrollment(userID, courseID int) {
	defer fmt.Println("done")
	fmt.Println("enroll user", userID, "into course", courseID)
}

func main() {
	fmt.Println("remaining:", remainingSeats(30, 10))

	cases := [][2]int{{30, 5}, {30, 40}, {30, 0}}
	for _, c := range cases {
		ok, err := canEnroll(c[0], c[1])
		fmt.Printf("canEnroll(%d, %d) => ok=%v err=%v\n", c[0], c[1], ok, err)
	}

	// Shadowing: the err inside if is a new variable for this block.
	seats := 30
	if left, err := enroll(seats, 5); err != nil {
		fmt.Println("enroll failed:", err)
	} else {
		fmt.Println("seats left after enroll:", left)
	}

	n := 10
	bump(n)
	fmt.Println("after bump (pass by value):", n)

	printEnrollment(7, 101)
}
