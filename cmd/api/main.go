package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/tayjava/go-lab/internal/api"
	"github.com/tayjava/go-lab/internal/config"
	"github.com/tayjava/go-lab/internal/course"
)

func main() {
	healthcheck := flag.Bool("healthcheck", false, "probe /healthz then exit 0 (healthy) or 1")
	flag.Parse()
	if *healthcheck {
		os.Exit(probe())
	}

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

// probe runs inside the container, so it reaches the server via 127.0.0.1.
// It runs before config.Load and does not require DATABASE_URL.
func probe() int {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return 1
	}
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://" + net.JoinHostPort("127.0.0.1", port) + "/healthz")
	if err != nil {
		return 1
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}
