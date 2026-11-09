package main

import (
	"context"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/tayjava/go-lab/internal/api"
	"github.com/tayjava/go-lab/internal/config"
	"github.com/tayjava/go-lab/internal/course"
	"github.com/tayjava/go-lab/internal/notify"
)

func main() {
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

	// This process has no signal context yet. Background stays active until exit,
	// so Close drains the queue. Canceling the pool context would drop queued jobs.
	pool := notify.NewPool(context.Background(), 4, 64, notify.FakeSMTP)
	defer pool.Close()

	repo := course.NewPostgresRepo(db)
	svc := course.NewService(repo)
	courses := api.NewCourseHandler(svc, cfg.EnrollmentCapDefault)
	enrollments := api.NewEnrollmentHandler(svc, pool)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", api.Healthz)

	// Intermediate lab routes (lesson 27)
	mux.HandleFunc("GET /courses/{id}", courses.Get)
	mux.HandleFunc("POST /courses", courses.Create)
	mux.HandleFunc("POST /courses/{id}/enrollments", enrollments.CreateNested)
	mux.HandleFunc("GET /courses/{id}/enrollments", enrollments.List)

	// Aliases kept from earlier lessons
	mux.HandleFunc("GET /api/courses/{id}", courses.Get)
	mux.HandleFunc("POST /api/courses", courses.Create)
	mux.HandleFunc("POST /api/enrollments", enrollments.Create)

	var h http.Handler = mux
	h = api.Recovery(h)
	h = api.Logging(h)

	log.Printf("listening on %s", cfg.HTTPAddr)
	log.Fatal(http.ListenAndServe(cfg.HTTPAddr, h))
}
