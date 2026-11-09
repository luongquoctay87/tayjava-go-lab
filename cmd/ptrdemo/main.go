package main

import (
	"fmt"

	"github.com/tayjava/go-lab/course"
)

func main() {
	c := course.New(101, "Golang from zero", 30)
	fmt.Println("before reserve:", c.Seats)

	if err := course.Reserve(&c); err != nil {
		fmt.Println("reserve error:", err)
	} else {
		fmt.Println("after reserve:", c.Seats)
	}

	beforeCopy := c.Seats
	course.ReserveCopy(c)
	fmt.Println("after reserveCopy (unchanged):", c.Seats, "(was", beforeCopy, ")")

	if err := course.Reserve(nil); err != nil {
		fmt.Println("reserve nil:", err)
	}
}
