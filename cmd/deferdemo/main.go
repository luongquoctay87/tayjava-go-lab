package main

import (
	"fmt"
	"os"
)

func readCourses(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 64)
	n, err := f.Read(buf)
	if err != nil {
		return err
	}
	fmt.Println("read bytes:", n)
	return nil
}

// enrollOrPanic panics on negative seats (invariant violation demo).
func enrollOrPanic(seats int) {
	if seats < 0 {
		panic("seats must not be negative")
	}
	fmt.Println("enrolled with seats budget:", seats)
}

// enroll reports negative seats as an error instead of panicking.
func enroll(seats int) error {
	if seats < 0 {
		return fmt.Errorf("seats must not be negative: %d", seats)
	}
	fmt.Println("enrolled with seats budget:", seats)
	return nil
}

func safePublish(title string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("publish panicked: %v", r)
		}
	}()
	if title == "" {
		panic("empty title")
	}
	fmt.Println("published:", title)
	return nil
}

func demoDeferOrder() {
	defer fmt.Println("first defer")
	defer fmt.Println("second defer")
	fmt.Println("body")
}

func main() {
	demoDeferOrder()

	if err := readCourses("go.mod"); err != nil {
		fmt.Println("readCourses go.mod:", err)
	}
	if err := readCourses("does-not-exist.json"); err != nil {
		fmt.Println("readCourses missing:", err)
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered enrollOrPanic:", r)
			}
		}()
		enrollOrPanic(-1)
	}()

	if err := enroll(-1); err != nil {
		fmt.Println("enroll error (no panic):", err)
	}

	if err := safePublish(""); err != nil {
		fmt.Println(err)
	}
	_ = safePublish("Go Beginner")
}
