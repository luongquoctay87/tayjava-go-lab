package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/tayjava/go-lab/internal/config"
	"github.com/tayjava/go-lab/internal/course"
	"github.com/tayjava/go-lab/internal/enrollment"
	"github.com/tayjava/go-lab/internal/resilience"
)

// Standalone enrollment worker binary (lesson 42): Go for the queue/worker
// while a JVM Course CRUD API can remain elsewhere.
func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := course.OpenDBWithRetry(cfg.DatabaseURL, 30, time.Second)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer db.Close()

	svc := course.NewService(course.NewPostgresRepo(db))
	m := &enrollment.Metrics{}
	enroller := enrollment.NewEnroller(svc, resilience.NewIdempotencyStore(), m, 64)
	enroller.Run(4)

	// Demo: enqueue a few jobs then drain on SIGTERM.
	_ = enroller.Enqueue(context.Background(), enrollment.Job{
		IdempotencyKey: "worker-demo-1",
		CourseID:       1,
		UserID:         1001,
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	slog.Info("worker_running", "queue", 64, "note", "Go worker beside optional JVM Course API")
	<-ctx.Done()
	slog.Info("worker_shutting_down")
	enroller.Close()
}
