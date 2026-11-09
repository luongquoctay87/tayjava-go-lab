package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds process settings loaded from the environment (twelve-factor light).
type Config struct {
	HTTPAddr             string
	DatabaseURL          string
	EnrollmentCapDefault int
	ShutdownTimeout      time.Duration
}

// Load reads env vars. DATABASE_URL is required; other fields have safe defaults.
func Load() (Config, error) {
	cfg := Config{
		HTTPAddr:             getenv("HTTP_ADDR", ":8080"),
		DatabaseURL:          os.Getenv("DATABASE_URL"),
		EnrollmentCapDefault: 30,
		ShutdownTimeout:      10 * time.Second,
	}

	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	if v := os.Getenv("ENROLLMENT_CAP_DEFAULT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			return Config{}, fmt.Errorf("invalid ENROLLMENT_CAP_DEFAULT: %q", v)
		}
		cfg.EnrollmentCapDefault = n
	}

	if v := os.Getenv("SHUTDOWN_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %w", err)
		}
		cfg.ShutdownTimeout = d
	}

	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
