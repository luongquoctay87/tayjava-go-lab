package main

import (
	"fmt"
	"time"
)

func demoBufferedJobs() {
	jobs := make(chan int64, 2)
	jobs <- 101
	jobs <- 102
	// Third send would block until a receiver drains — do it in a goroutine.
	go func() {
		jobs <- 103
		close(jobs) // producer owns close
	}()

	fmt.Println("--- drain buffered jobs ---")
	for id := range jobs {
		fmt.Println("process course", id)
	}
}

func demoSelectTimeout() {
	enroll := make(chan string, 1)
	cancel := make(chan struct{})

	go func() {
		time.Sleep(50 * time.Millisecond)
		enroll <- "an@tayjava.com"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		close(cancel)
	}()

	fmt.Println("--- select enroll / cancel / timeout ---")
	select {
	case email := <-enroll:
		fmt.Println("enrolled", email)
	case <-cancel:
		fmt.Println("canceled")
	case <-time.After(100 * time.Millisecond):
		fmt.Println("timeout waiting enroll")
	}
}

func demoSelectDrain() {
	out := make(chan string, 2)
	go func() {
		out <- "course-101"
		out <- "course-102"
		close(out)
	}()

	fmt.Println("--- select until closed ---")
	for {
		select {
		case id, ok := <-out:
			if !ok {
				fmt.Println("drained")
				return
			}
			fmt.Println("got", id)
		}
	}
}

func main() {
	demoBufferedJobs()
	demoSelectTimeout()
	demoSelectDrain()
}
