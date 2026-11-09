package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestIssueParseRoundTrip(t *testing.T) {
	t.Setenv("JWT_SECRET", "lab-test-secret")
	tok, err := Issue("42", time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	sub, err := Parse(tok)
	if err != nil || sub != "42" {
		t.Fatalf("sub=%q err=%v", sub, err)
	}
}

func TestExpiredTokenRejected(t *testing.T) {
	t.Setenv("JWT_SECRET", "lab-test-secret")
	tok, err := Issue("7", -time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Parse(tok); err == nil {
		t.Fatal("expected expiry error")
	}
}

func TestBearerUnauthorized(t *testing.T) {
	t.Setenv("JWT_SECRET", "lab-test-secret")
	h := Bearer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", rr.Code)
	}
}

func TestSecretRequired(t *testing.T) {
	_ = os.Unsetenv("JWT_SECRET")
	if _, err := Issue("1", time.Minute); err == nil {
		t.Fatal("expected missing secret")
	}
}
