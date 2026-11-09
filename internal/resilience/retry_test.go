package resilience

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetrySucceedsAfterTemporary(t *testing.T) {
	var n atomic.Int32
	err := Retry(context.Background(), 4, time.Millisecond, func(context.Context) error {
		if n.Add(1) < 3 {
			return ErrTemporary
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if n.Load() != 3 {
		t.Fatalf("attempts=%d", n.Load())
	}
}

func TestRetryDoesNotRetrySoldOut(t *testing.T) {
	var n atomic.Int32
	err := Retry(context.Background(), 4, time.Millisecond, func(context.Context) error {
		n.Add(1)
		return ErrSoldOut
	})
	if err != ErrSoldOut {
		t.Fatalf("err=%v", err)
	}
	if n.Load() != 1 {
		t.Fatalf("attempts=%d want 1", n.Load())
	}
}

func TestIdempotencySameKey(t *testing.T) {
	s := NewIdempotencyStore()
	s.Put("abc", 42)
	id, ok := s.Get("abc")
	if !ok || id != 42 {
		t.Fatalf("got=%d ok=%v", id, ok)
	}
	s.Put("abc", 99) // ignored
	id, _ = s.Get("abc")
	if id != 42 {
		t.Fatalf("id=%d", id)
	}
}
