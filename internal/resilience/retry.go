package resilience

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

// ErrTemporary marks a retryable failure (timeouts, 502/503/504/429).
var ErrTemporary = errors.New("temporary")

// ErrSoldOut is a business error (no seats / 409). Do not retry it.
var ErrSoldOut = errors.New("sold out")

// RetryAfterError is a transient HTTP status that told the client how long to wait.
type RetryAfterError struct {
	Status int
	After  time.Duration
}

func (e *RetryAfterError) Error() string {
	return fmt.Sprintf("status %d, retry after %s", e.Status, e.After)
}

// Unwrap lets errors.Is(err, ErrTemporary) succeed for 429 and 503.
func (e *RetryAfterError) Unwrap() error { return ErrTemporary }

// Retry runs fn with full-jitter backoff until success, a non-retryable error, or attempts run out.
// The wait is a random duration in [0, min(maxDelay, base*2^attempt)).
// When err is a RetryAfterError whose After is longer than that draw, Retry waits After instead.
func Retry(ctx context.Context, attempts int, base, maxDelay time.Duration, fn func(context.Context) error) error {
	var err error
	for i := 0; i < attempts; i++ {
		if err = ctx.Err(); err != nil {
			return err
		}
		err = fn(ctx)
		if err == nil || !IsRetryable(err) {
			return err
		}
		if i == attempts-1 {
			break
		}
		delay := backoff(i, base, maxDelay)
		var ra *RetryAfterError
		if errors.As(err, &ra) && ra.After > delay {
			delay = ra.After
		}
		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}
	}
	return err
}

// backoff uses full jitter: a random wait in [0, min(maxDelay, base*2^attempt)).
// base <= maxDelay>>attempt is checked before shifting so base<<attempt cannot overflow int64.
func backoff(attempt int, base, maxDelay time.Duration) time.Duration {
	if base <= 0 || maxDelay <= 0 || attempt < 0 {
		return 0
	}
	ceiling := maxDelay
	// Shifts of 63 bits turn a positive duration into a negative int64; rand.N would panic.
	if attempt < 63 && base <= maxDelay>>attempt {
		ceiling = base << attempt
	}
	if ceiling <= 0 {
		return 0
	}
	return rand.N(ceiling)
}

// IsRetryable reports whether err should be retried.
// Timeouts and transient errors (502/503/504/429) are retryable. Sold out / 409 is not.
func IsRetryable(err error) bool {
	return errors.Is(err, context.DeadlineExceeded) || errors.Is(err, ErrTemporary)
}

// classify maps an HTTP response to a retryable or terminal error.
func classify(resp *http.Response) error {
	switch code := resp.StatusCode; {
	case code < 300:
		return nil
	case code == http.StatusTooManyRequests, code == http.StatusServiceUnavailable:
		return &RetryAfterError{Status: code, After: parseRetryAfter(resp.Header.Get("Retry-After"))}
	case code == http.StatusBadGateway, code == http.StatusGatewayTimeout:
		return fmt.Errorf("%w: status %d", ErrTemporary, code)
	default:
		// 400, 401, 403, 409 (no seats / already enrolled), 500: do not retry.
		return fmt.Errorf("enroll failed: status %d", code)
	}
}

// parseRetryAfter accepts delta-seconds or an HTTP-date. Unknown values yield 0.
func parseRetryAfter(v string) time.Duration {
	if s, err := strconv.Atoi(v); err == nil && s > 0 {
		return time.Duration(s) * time.Second
	}
	if t, err := http.ParseTime(v); err == nil {
		return time.Until(t)
	}
	return 0
}
