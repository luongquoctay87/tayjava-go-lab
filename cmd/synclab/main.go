package main

import (
	"fmt"
	"sync"

	"github.com/tayjava/go-lab/internal/seatledger"
)

// notifyAll waits for every email send with WaitGroup only (no shared mutable data).
func notifyAll(emails []string) {
	var wg sync.WaitGroup
	wg.Add(len(emails))
	for _, email := range emails {
		email := email
		go func() {
			defer wg.Done()
			fmt.Println("sent:", email)
		}()
	}
	wg.Wait()
	fmt.Println("all emails sent")
}

func main() {
	fmt.Println("--- WaitGroup notify ---")
	notifyAll([]string{"an@tayjava.com", "binh@tayjava.com", "chi@tayjava.com"})

	fmt.Println("--- Mutex seat ledger ---")
	ledger := seatledger.New(map[int64]int{101: 2})
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ok := ledger.Reserve(101)
			fmt.Println("reserved", ok)
		}()
	}
	wg.Wait()
	fmt.Println("left", ledger.Remaining(101))
}
