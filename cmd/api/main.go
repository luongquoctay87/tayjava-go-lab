package main

import (
	"log"
	"net/http"

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

	db, err := course.OpenDB(cfg.DatabaseURL)
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
	mux.HandleFunc("GET /api/courses/{id}", courses.Get)
	mux.HandleFunc("POST /api/courses", courses.Create)
	mux.HandleFunc("POST /api/enrollments", enrollments.Create)

	var h http.Handler = mux
	h = api.Recovery(h)
	h = api.Logging(h)

	log.Printf("listening on %s", cfg.HTTPAddr)
	log.Fatal(http.ListenAndServe(cfg.HTTPAddr, h))
}
