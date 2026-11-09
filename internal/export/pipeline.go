package export

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

// Pipeline runs a 3-stage fake export: list IDs → fetch rows → collect results.
// failID, when non-zero, makes that course ID return an error so WithContext cancels siblings.
func Pipeline(ctx context.Context, courseIDs []int64, failID int64) ([]string, error) {
	g, ctx := errgroup.WithContext(ctx)

	ids := make(chan int64)
	rows := make(chan string)

	g.Go(func() error {
		defer close(ids)
		for _, id := range courseIDs {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case ids <- id:
			}
		}
		return nil
	})

	var fetchGroup errgroup.Group
	fetchGroup.SetLimit(4)
	g.Go(func() error {
		defer close(rows)
		for id := range ids {
			id := id
			fetchGroup.Go(func() error {
				if failID != 0 && id == failID {
					return fmt.Errorf("export course %d failed", id)
				}
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(5 * time.Millisecond):
				}
				row := fmt.Sprintf("course=%d enrollments=ok", id)
				select {
				case <-ctx.Done():
					return ctx.Err()
				case rows <- row:
					return nil
				}
			})
		}
		return fetchGroup.Wait()
	})

	var out []string
	g.Go(func() error {
		for row := range rows {
			out = append(out, row)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return out, nil
}
