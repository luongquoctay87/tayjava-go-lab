package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/tayjava/go-lab/internal/api"
	"github.com/tayjava/go-lab/internal/course"
)

func main() {
	repo, cleanup := mustRepo()
	if cleanup != nil {
		defer cleanup()
	}

	svc := course.NewService(repo)
	courses := api.NewCourseHandler(svc)
	enrollments := api.NewEnrollmentHandler(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", api.Healthz)
	mux.HandleFunc("GET /api/courses/{id}", courses.Get)
	mux.HandleFunc("POST /api/courses", courses.Create)
	mux.HandleFunc("POST /api/enrollments", enrollments.Create)

	var h http.Handler = mux
	h = api.Recovery(h)
	h = api.Logging(h)

	addr := ":8080"
	log.Println("listening", addr)
	log.Fatal(http.ListenAndServe(addr, h))
}

// mustRepo prefers Postgres when DATABASE_URL is set; otherwise MemoryRepo for local demos.
func mustRepo() (course.Repository, func()) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("DATABASE_URL unset — using MemoryRepo (data lost on restart)")
		return course.NewMemoryRepo(), nil
	}
	db, err := course.OpenDB(dsn)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	log.Println("postgres ok")
	return course.NewPostgresRepo(db), func() { _ = db.Close() }
}
