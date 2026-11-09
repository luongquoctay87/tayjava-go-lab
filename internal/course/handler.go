package course

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// Healthz handles GET /healthz.
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

// GetCourse handles GET /api/courses/{id} and returns JSON.
func GetCourse(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}
	log.Printf("get course method=%s path=%s id=%d", r.Method, r.URL.Path, id)
	// Lab stub until memory/Postgres repository (lessons 21–22).
	writeJSON(w, http.StatusOK, Course{
		ID:    id,
		Title: "Go Context",
		Seats: 30,
		Level: "intermediate",
	})
}

// CreateCourse handles POST /api/courses with a JSON body.
func CreateCourse(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB
	var in Course
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if in.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title required"})
		return
	}
	if in.Level == "" {
		in.Level = "beginner"
	}
	in.ID = 101 // lab: fake assigned id
	writeJSON(w, http.StatusCreated, in)
}

// CreateEnrollment handles POST /api/enrollments (query params; JSON body in lesson 21).
func CreateEnrollment(w http.ResponseWriter, r *http.Request) {
	_ = r.Context()
	courseID := r.URL.Query().Get("courseId")
	userID := r.URL.Query().Get("userId")
	if courseID == "" || userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "courseId and userId required"})
		return
	}
	log.Printf("enroll method=%s path=%s courseId=%s userId=%s", r.Method, r.URL.Path, courseID, userID)
	writeJSON(w, http.StatusCreated, map[string]string{
		"userId":   userID,
		"courseId": courseID,
		"status":   "enrolled",
	})
}
