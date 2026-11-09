package export

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestPipelineOK(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	rows, err := Pipeline(ctx, []int64{1, 2, 3}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 {
		t.Fatalf("rows=%d want 3", len(rows))
	}
}

func TestPipelineCancelsOnError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := Pipeline(ctx, []int64{10, 20, 30}, 20)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, context.Canceled) && err.Error() != "export course 20 failed" {
		// first error from failing stage should surface
		if err.Error() != "export course 20 failed" {
			t.Fatalf("unexpected err: %v", err)
		}
	}
}

func TestPipelineRespectsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := Pipeline(ctx, []int64{1, 2}, 0)
	if err == nil {
		t.Fatal("expected canceled error")
	}
}
