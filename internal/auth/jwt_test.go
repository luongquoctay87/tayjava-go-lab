package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "0123456789abcdef0123456789abcdef"

func TestIssueParseRoundTrip(t *testing.T) {
	a := testAuth(t)
	tok, err := a.Issue(42, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	sub, err := a.Parse(tok)
	if err != nil || sub != 42 {
		t.Fatalf("sub=%d err=%v", sub, err)
	}
}

func TestExpiredTokenRejected(t *testing.T) {
	a := testAuth(t)
	tok, err := a.Issue(7, -time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := a.Parse(tok); err == nil {
		t.Fatal("expected expiry error")
	}
}

func TestLeewayAllowsSmallClockDrift(t *testing.T) {
	a := testAuth(t)
	tok := signClaims(t, testSecret, jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "7",
		Issuer:    issuer,
		Audience:  jwt.ClaimStrings{audience},
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Minute)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-10 * time.Second)),
	})
	sub, err := a.Parse(tok)
	if err != nil || sub != 7 {
		t.Fatalf("sub=%d err=%v", sub, err)
	}
}

func TestRejectsMissingExpSubAndOtherAlg(t *testing.T) {
	a := testAuth(t)
	base := jwt.RegisteredClaims{
		Subject:  "7",
		Issuer:   issuer,
		Audience: jwt.ClaimStrings{audience},
		IssuedAt: jwt.NewNumericDate(time.Now()),
	}
	missingExp := signClaims(t, testSecret, jwt.SigningMethodHS256, base)
	if _, err := a.Parse(missingExp); err == nil {
		t.Fatal("expected missing exp to fail")
	}

	noSub := base
	noSub.Subject = ""
	noSub.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
	if _, err := a.Parse(signClaims(t, testSecret, jwt.SigningMethodHS256, noSub)); err == nil {
		t.Fatal("expected missing sub to fail")
	}

	badSub := noSub
	badSub.Subject = "nope"
	if _, err := a.Parse(signClaims(t, testSecret, jwt.SigningMethodHS256, badSub)); err == nil {
		t.Fatal("expected non-numeric sub to fail")
	}

	zero := badSub
	zero.Subject = "0"
	if _, err := a.Parse(signClaims(t, testSecret, jwt.SigningMethodHS256, zero)); err == nil {
		t.Fatal("expected non-positive sub to fail")
	}

	hs384Key := testSecret + testSecret
	otherAlg := signClaims(t, hs384Key, jwt.SigningMethodHS384, jwt.RegisteredClaims{
		Subject:   "7",
		Issuer:    issuer,
		Audience:  jwt.ClaimStrings{audience},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
	})
	if _, err := a.Parse(otherAlg); err == nil {
		t.Fatal("expected HS384 token to fail")
	}
}

func TestShortSecretRejected(t *testing.T) {
	if _, err := New([]byte("short")); err == nil {
		t.Fatal("expected short secret error")
	}
	if _, err := New(nil); err == nil {
		t.Fatal("expected empty secret error")
	}
}

func TestBearerUnauthorizedAndLowercaseScheme(t *testing.T) {
	a := testAuth(t)
	h := a.Bearer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, ok := UserIDFrom(r.Context())
		if !ok || id != 42 {
			t.Fatalf("id=%d ok=%v", id, ok)
		}
		w.WriteHeader(http.StatusOK)
	}))

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/", nil))
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status=%d", rr.Code)
	}
	if got := rr.Header().Get("WWW-Authenticate"); got != `Bearer realm="tayjava"` {
		t.Fatalf("WWW-Authenticate=%q", got)
	}

	tok, err := a.Issue(42, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "bearer "+tok)
	rr = httptest.NewRecorder()
	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("lowercase bearer status=%d body=%s", rr.Code, rr.Body.String())
	}
}

func testAuth(t *testing.T) *Auth {
	t.Helper()
	a, err := New([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func signClaims(t *testing.T, secret string, method jwt.SigningMethod, claims jwt.RegisteredClaims) string {
	t.Helper()
	s, err := jwt.NewWithClaims(method, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatal(err)
	}
	return s
}
