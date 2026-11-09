package main

import (
	"log"
	"net/http"

	"github.com/tayjava/go-lab/internal/api"
	"github.com/tayjava/go-lab/internal/course"
)

func main() {
	repo := course.NewMemoryRepo()
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
