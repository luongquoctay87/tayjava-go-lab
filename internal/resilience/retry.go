package resilience

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

// ErrTemporary marks a retryable failure.
var ErrTemporary = errors.New("temporary")

// ErrSoldOut is a business error — do not retry.
var ErrSoldOut = errors.New("sold out")

// Retry runs fn with exponential backoff + jitter until success, non-retryable error, or attempts exhausted.
func Retry(ctx context.Context, attempts int, base time.Duration, fn func(context.Context) error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = ctx.Err(); err != nil {
			return err
		}
		err = fn(ctx)
		if err == nil {
			return nil
		}
		if !IsRetryable(err) {
			return err
		}
		if i == attempts-1 {
			break
		}
		backoff := base * time.Duration(1<<i)
		jitter := time.Duration(0)
		if base > 0 {
			jitter = time.Duration(rand.Int63n(int64(base)))
		}
		timer := time.NewTimer(backoff + jitter)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return err
}

// IsRetryable reports whether err should be retried.
func IsRetryable(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrTemporary)
}
