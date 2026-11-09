package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	"github.com/tayjava/go-lab/internal/api"
	"github.com/tayjava/go-lab/internal/config"
	"github.com/tayjava/go-lab/internal/course"
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

	repo := course.NewPostgresRepo(db)
	svc := course.NewService(repo)
	courses := api.NewCourseHandler(svc, cfg.EnrollmentCapDefault)
	enrollments := api.NewEnrollmentHandler(svc)

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
