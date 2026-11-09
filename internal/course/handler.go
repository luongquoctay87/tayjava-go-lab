package course

import (
	"fmt"
	"log"
	"net/http"
)

// GetCourse handles GET /api/courses/{id} (plain text skeleton for lesson 19).
func GetCourse(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_ = ctx // pass to loaders later (DB / service)

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	log.Printf("get course method=%s path=%s id=%s", r.Method, r.URL.Path, id)
	fmt.Fprintf(w, "course id=%s\n", id)
}

// CreateEnrollment handles POST /api/enrollments via query params (JSON in lesson 20).
func CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_ = ctx

	courseID := r.URL.Query().Get("courseId")
	userID := r.URL.Query().Get("userId")
	if courseID == "" || userID == "" {
		http.Error(w, "courseId and userId required", http.StatusBadRequest)
		return
	}
	log.Printf("enroll method=%s path=%s courseId=%s userId=%s", r.Method, r.URL.Path, courseID, userID)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "enrolled user=%s course=%s\n", userID, courseID)
}

// Healthz handles GET /healthz.
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
