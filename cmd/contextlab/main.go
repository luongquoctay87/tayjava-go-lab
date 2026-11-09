package main

import (
	"context"
	"fmt"
	"time"
)

type ctxKey string

const requestIDKey ctxKey = "requestID"

func withRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

func requestIDFrom(ctx context.Context) string {
	v, _ := ctx.Value(requestIDKey).(string)
	return v
}

func loadCourse(ctx context.Context, id int64) (string, error) {
	select {
	case <-time.After(800 * time.Millisecond):
		return fmt.Sprintf("Course-%d", id), nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func checkSeat(ctx context.Context, courseID int64) error {
	select {
	case <-time.After(80 * time.Millisecond):
		_ = courseID
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func saveEnrollment(ctx context.Context, courseID, userID int64) error {
	select {
	case <-time.After(80 * time.Millisecond):
		fmt.Printf("enrolled user=%d course=%d requestID=%s\n", userID, courseID, requestIDFrom(ctx))
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Enroll walks a fake call chain that honors the same parent context.
func Enroll(ctx context.Context, courseID, userID int64) error {
	if err := checkSeat(ctx, courseID); err != nil {
		return err
	}
	return saveEnrollment(ctx, courseID, userID)
}

func main() {
	fmt.Println("--- timeout too short (expect deadline) ---")
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	title, err := loadCourse(ctx, 101)
	if err != nil {
		fmt.Println("failed:", err)
	} else {
		fmt.Println("ok:", title)
	}

	fmt.Println("--- enroll with request ID (expect success) ---")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel2()
	ctx2 = withRequestID(ctx2, "req-42")
	if err := Enroll(ctx2, 101, 7); err != nil {
		fmt.Println("enroll failed:", err)
	} else {
		fmt.Println("enroll ok")
	}

	fmt.Println("--- enroll deadline (80ms+80ms > 100ms) ---")
	ctx3, cancel3 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel3()
	if err := Enroll(ctx3, 202, 9); err != nil {
		fmt.Println("enroll failed:", err)
	}

	fmt.Println("--- manual cancel ---")
	ctx4, cancel4 := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel4()
	}()
	title, err = loadCourse(ctx4, 303)
	fmt.Println("result:", title, err)
}
