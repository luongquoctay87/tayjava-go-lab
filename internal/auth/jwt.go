package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const userIDKey ctxKey = "userID"

// ErrMissingSecret means JWT_SECRET is not configured.
var ErrMissingSecret = errors.New("JWT_SECRET is required")

// Secret reads the lab HMAC secret from the environment.
func Secret() ([]byte, error) {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		return nil, ErrMissingSecret
	}
	return []byte(s), nil
}

// Issue creates an HS256 token with sub=userID and exp=now+ttl.
func Issue(userID string, ttl time.Duration) (string, error) {
	secret, err := Secret()
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(ttl).Unix(),
		"iat": time.Now().Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(secret)
}

// Parse validates HS256 and returns the subject (user id).
func Parse(tokenStr string) (string, error) {
	secret, err := Secret()
	if err != nil {
		return "", err
	}
	tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		return "", jwt.ErrTokenInvalidClaims
	}
	return claims.GetSubject()
}

// UserIDFrom returns the authenticated user id from context.
func UserIDFrom(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userIDKey).(string)
	return v, ok
}

// Bearer requires Authorization: Bearer <jwt> and injects user_id into context.
func Bearer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if !strings.HasPrefix(h, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		userID, err := Parse(strings.TrimPrefix(h, "Bearer "))
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
