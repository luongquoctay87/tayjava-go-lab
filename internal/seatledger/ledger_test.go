package seatledger_test

import (
	"sync"
	"testing"

	"github.com/tayjava/go-lab/internal/seatledger"
)

func TestReserveConcurrent(t *testing.T) {
	const seats = 3
	const goroutines = 10

	ledger := seatledger.New(map[int64]int{101: seats})
	var wg sync.WaitGroup
	results := make([]bool, goroutines)

	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		i := i
		go func() {
			defer wg.Done()
			results[i] = ledger.Reserve(101)
		}()
	}
	wg.Wait()

	okCount := 0
	for _, ok := range results {
		if ok {
			okCount++
		}
	}
	if okCount != seats {
		t.Fatalf("true reserves: got %d want %d", okCount, seats)
	}
	if left := ledger.Remaining(101); left != 0 {
		t.Fatalf("remaining: got %d want 0", left)
	}
}

func TestReserveEmpty(t *testing.T) {
	ledger := seatledger.New(map[int64]int{202: 0})
	if ledger.Reserve(202) {
		t.Fatal("expected false when no seats")
	}
	if ledger.Reserve(999) {
		t.Fatal("expected false for unknown course")
	}
}
