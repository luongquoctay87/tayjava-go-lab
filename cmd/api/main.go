package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"github.com/tayjava/go-lab/internal/api"
	"github.com/tayjava/go-lab/internal/config"
	"github.com/tayjava/go-lab/internal/course"
	"github.com/tayjava/go-lab/internal/metrics"
	"github.com/tayjava/go-lab/internal/notify"
)

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
	log.Println("postgres ok")

	// Pool lifetime ≠ signal ctx: drain via Close after HTTP Shutdown.
	pool := notify.NewPool(context.Background(), 4, 64, notify.FakeSMTP)

	repo := course.NewPostgresRepo(db)
	svc := course.NewService(repo)
	counters := &metrics.Counters{}
	courses := api.NewCourseHandler(svc, cfg.EnrollmentCapDefault)
	enrollments := api.NewEnrollmentHandler(svc, pool, counters)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", api.Healthz)
	mux.HandleFunc("GET /metrics", counters.Handler())
	mux.HandleFunc("GET /courses/{id}", courses.Get)
	mux.HandleFunc("POST /courses", courses.Create)
	mux.HandleFunc("POST /courses/{id}/enrollments", enrollments.CreateNested)
	mux.HandleFunc("GET /courses/{id}/enrollments", enrollments.List)
	mux.HandleFunc("GET /api/courses/{id}", courses.Get)
	mux.HandleFunc("POST /api/courses", courses.Create)
	mux.HandleFunc("POST /api/enrollments", enrollments.Create)

	var h http.Handler = mux
	h = api.Recovery(h)
	h = api.Logging(h)

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: h}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// pprof on localhost only — never expose :6060 publicly in the lab.
	go func() {
		log.Println("pprof listening on 127.0.0.1:6060")
		if err := http.ListenAndServe("127.0.0.1:6060", nil); err != nil {
			log.Printf("pprof: %v", err)
		}
	}()

	go func() {
		log.Printf("listening on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown: %v", err)
	}
	pool.Close()
	log.Println("bye")
}
