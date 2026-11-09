package main

import (
	"log"
	"net/http"

	"github.com/tayjava/go-lab/internal/course"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", course.Healthz)
	mux.HandleFunc("GET /api/courses/{id}", course.GetCourse)
	mux.HandleFunc("POST /api/enrollments", course.CreateEnrollment)

	addr := ":8080"
	log.Println("listening", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
