package main

import (
	"fmt"

	"github.com/tayjava/go-lab/course"
)

func main() {
	c := course.New(101, "Golang from zero", 30)
	fmt.Println(c.Title, c.Level())
	// c.level is unexported — accessing it here would not compile.
}
