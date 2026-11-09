package config_test

import (
	"testing"
	"time"

	"github.com/tayjava/go-lab/internal/config"
)

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	if _, err := config.Load(); err == nil {
		t.Fatal("expected error when DATABASE_URL is empty")
	}
}

func TestLoadDefaultsAndOverrides(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://tayjava:tayjava@localhost:5432/tayjava?sslmode=disable")
	t.Setenv("HTTP_ADDR", ":8081")
	t.Setenv("ENROLLMENT_CAP_DEFAULT", "50")
	t.Setenv("SHUTDOWN_TIMEOUT", "5s")

	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPAddr != ":8081" {
		t.Fatalf("HTTPAddr=%q", cfg.HTTPAddr)
	}
	if cfg.EnrollmentCapDefault != 50 {
		t.Fatalf("EnrollmentCapDefault=%d", cfg.EnrollmentCapDefault)
	}
	if cfg.ShutdownTimeout != 5*time.Second {
		t.Fatalf("ShutdownTimeout=%v", cfg.ShutdownTimeout)
	}
}
