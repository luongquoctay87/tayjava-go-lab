package resilience

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetrySucceedsAfterTemporary(t *testing.T) {
	var n atomic.Int32
	err := Retry(context.Background(), 4, time.Millisecond, 5*time.Millisecond, func(context.Context) error {
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
	err := Retry(context.Background(), 4, time.Millisecond, 5*time.Millisecond, func(context.Context) error {
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

func TestRetryHonorsRetryAfter(t *testing.T) {
	var n atomic.Int32
	start := time.Now()
	err := Retry(context.Background(), 2, time.Millisecond, 5*time.Millisecond, func(context.Context) error {
		if n.Add(1) == 1 {
			return &RetryAfterError{Status: http.StatusTooManyRequests, After: 40 * time.Millisecond}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if n.Load() != 2 {
		t.Fatalf("attempts=%d", n.Load())
	}
	if time.Since(start) < 30*time.Millisecond {
		t.Fatal("429 Retry-After was not honored")
	}
}

func TestBackoffStaysUnderCap(t *testing.T) {
	base := time.Millisecond
	maxDelay := 50 * time.Millisecond
	for _, attempt := range []int{0, 10, 64} {
		for i := 0; i < 30; i++ {
			d := backoff(attempt, base, maxDelay)
			if d < 0 || d >= maxDelay {
				t.Fatalf("attempt=%d delay=%s outside [0, %s)", attempt, d, maxDelay)
			}
		}
	}
	if backoff(0, 0, time.Second) != 0 {
		t.Fatal("non-positive base must return 0")
	}
}

func TestClassifyRetryAfterAndConflict(t *testing.T) {
	tooMany := &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Header:     http.Header{"Retry-After": []string{"2"}},
	}
	err := classify(tooMany)
	var ra *RetryAfterError
	if !errors.As(err, &ra) {
		t.Fatalf("err=%v", err)
	}
	if ra.Status != http.StatusTooManyRequests || ra.After != 2*time.Second {
		t.Fatalf("got status=%d after=%s", ra.Status, ra.After)
	}
	if !IsRetryable(err) {
		t.Fatal("429 must be retryable")
	}

	when := time.Now().Add(3 * time.Second).UTC()
	got := parseRetryAfter(when.Format(http.TimeFormat))
	if got < 2*time.Second || got > 4*time.Second {
		t.Fatalf("http-date delay=%s", got)
	}

	conflict := classify(&http.Response{StatusCode: http.StatusConflict, Header: make(http.Header)})
	if IsRetryable(conflict) {
		t.Fatal("409 must not be retryable")
	}
}
