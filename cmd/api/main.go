package main

import (
	"context"
	"errors"
	"expvar"
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
	"github.com/tayjava/go-lab/internal/auth"
	"github.com/tayjava/go-lab/internal/config"
	"github.com/tayjava/go-lab/internal/course"
	"github.com/tayjava/go-lab/internal/enrollment"
	"github.com/tayjava/go-lab/internal/metrics"
	"github.com/tayjava/go-lab/internal/notify"
	"github.com/tayjava/go-lab/internal/resilience"
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

	pool := notify.NewPool(context.Background(), 4, 64, notify.FakeSMTP)

	repo := course.NewPostgresRepo(db)
	svc := course.NewService(repo)
	counters := &metrics.Counters{}
	idem := resilience.NewIdempotencyStore()
	courses := api.NewCourseHandler(svc, cfg.EnrollmentCapDefault)
	enrollments := api.NewEnrollmentHandler(svc, pool, counters, idem)

	advMetrics := &enrollment.Metrics{}
	enroller := enrollment.NewEnroller(svc, idem, advMetrics, 64)
	asyncEnroll := enrollment.NewHandler(enroller)
	expvar.Publish("enroll_accepted_total", expvar.Func(func() any { return advMetrics.Accepted.Load() }))
	expvar.Publish("enroll_rejected_total", expvar.Func(func() any { return advMetrics.Rejected.Load() }))
	expvar.Publish("enroll_processed_total", expvar.Func(func() any { return advMetrics.Processed.Load() }))
	expvar.Publish("enroll_failed_total", expvar.Func(func() any { return advMetrics.Failed.Load() }))
	expvar.Publish("enroll_queue_depth", expvar.Func(func() any { return enroller.Depth() }))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", api.Healthz)
	mux.Handle("GET /metrics", expvar.Handler())
	mux.HandleFunc("POST /login", auth.LoginHandler)
	mux.HandleFunc("GET /courses/{id}", courses.Get)
	mux.HandleFunc("POST /courses", courses.Create)
	mux.Handle("POST /courses/{id}/enrollments", auth.Bearer(http.HandlerFunc(enrollments.CreateNested)))
	mux.HandleFunc("GET /courses/{id}/enrollments", enrollments.List)
	mux.HandleFunc("GET /api/courses/{id}", courses.Get)
	mux.HandleFunc("POST /api/courses", courses.Create)
	// Advanced async enroll (lesson 39) — JWT + 202/503
	mux.Handle("POST /api/enrollments", auth.Bearer(http.HandlerFunc(asyncEnroll.CreateAsync)))

	var h http.Handler = mux
	h = api.Recovery(h)
	h = api.Logging(h)

	srv := &http.Server{Addr: cfg.HTTPAddr, Handler: h}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	enroller.Run(4)

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
		slog.Warn("http shutdown incomplete", "err", err)
	}
	// Stop new enqueues, then close the job channel. Workers drain with for range
	// on their own context; this deadline bounds how long we wait for them.
	enroller.Close()
	drained := make(chan struct{})
	go func() {
		enroller.Wait()
		close(drained)
	}()
	select {
	case <-drained:
		slog.Info("workers drained")
	case <-shutdownCtx.Done():
		slog.Error("shutdown deadline hit, jobs dropped", "jobs_left", enroller.Depth())
	}
	pool.Close()
	log.Println("bye")
}
