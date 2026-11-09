package main

import (
	"fmt"
	"sync"
	"time"
)

// enrollNotify simulates sending an enrollment email (I/O delay).
func enrollNotify(courseID int64, email string) {
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("enrolled course=%d email=%s\n", courseID, email)
}

func notifySequential(courseID int64, emails []string) {
	start := time.Now()
	for _, email := range emails {
		enrollNotify(courseID, email)
	}
	fmt.Printf("sequential done in %s\n", time.Since(start).Round(time.Millisecond))
}

func notifyConcurrent(courseID int64, emails []string) {
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(len(emails))
	for _, email := range emails {
		// This module declares go 1.22, so each iteration already has its own email.
		// email := email is only for a go.mod that declares a version below 1.22.
		go func() {
			defer wg.Done()
			enrollNotify(courseID, email)
		}()
	}
	wg.Wait() // join point — do not fire-and-forget
	fmt.Printf("concurrent done in %s\n", time.Since(start).Round(time.Millisecond))
}

func main() {
	const courseID int64 = 101
	emails := []string{
		"an@tayjava.com",
		"binh@tayjava.com",
		"chi@tayjava.com",
		"dung@tayjava.com",
		"em@tayjava.com",
	}

	fmt.Println("--- sequential ---")
	notifySequential(courseID, emails)

	fmt.Println("--- concurrent (WaitGroup) ---")
	notifyConcurrent(courseID, emails)

	fmt.Println("all notified")
}
