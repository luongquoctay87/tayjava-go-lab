package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const userIDKey ctxKey = "userID"

const (
	issuer   = "tayjava-course-api"
	audience = "tayjava-web"
)

// ErrInvalidToken means the bearer token failed validation.
var ErrInvalidToken = errors.New("invalid token")

// Auth issues and checks HS256 lab tokens. The HMAC secret comes from the caller (JWT_SECRET).
type Auth struct {
	secret []byte
	parser *jwt.Parser
}

// New rejects a secret shorter than 32 bytes (RFC 7518: HMAC key >= hash output).
func New(secret []byte) (*Auth, error) {
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET must be at least 32 bytes")
	}
	return &Auth{
		secret: secret,
		parser: jwt.NewParser(
			jwt.WithValidMethods([]string{"HS256"}),
			jwt.WithExpirationRequired(),
			jwt.WithLeeway(30*time.Second),
			jwt.WithIssuer(issuer),
			jwt.WithAudience(audience),
		),
	}, nil
}

// Issue creates an HS256 token with sub, exp, iss, and aud.
func (a *Auth) Issue(userID int64, ttl time.Duration) (string, error) {
	if userID <= 0 {
		return "", fmt.Errorf("%w: bad sub", ErrInvalidToken)
	}
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		Issuer:    issuer,
		Audience:  jwt.ClaimStrings{audience},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(a.secret)
}

// Parse validates the token and returns a positive user id from sub.
func (a *Auth) Parse(tokenStr string) (int64, error) {
	var claims jwt.RegisteredClaims
	tok, err := a.parser.ParseWithClaims(tokenStr, &claims, func(*jwt.Token) (any, error) {
		return a.secret, nil
	})
	if err != nil {
		return 0, fmt.Errorf("%w: %w", ErrInvalidToken, err)
	}
	if !tok.Valid {
		return 0, ErrInvalidToken
	}
	if claims.Subject == "" {
		return 0, fmt.Errorf("%w: missing sub", ErrInvalidToken)
	}
	userID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil || userID <= 0 {
		return 0, fmt.Errorf("%w: bad sub %q", ErrInvalidToken, claims.Subject)
	}
	return userID, nil
}

// UserIDFrom returns the authenticated user id from context.
func UserIDFrom(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(userIDKey).(int64)
	return v, ok
}

// Bearer requires Authorization: Bearer <jwt> (scheme is case-insensitive) and stores user id in context.
func (a *Auth) Bearer(next http.Handler) http.Handler {
	const prefix = "Bearer "
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if len(h) <= len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
			unauthorized(w)
			return
		}
		userID, err := a.Parse(h[len(prefix):])
		if err != nil {
			unauthorized(w)
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Bearer realm="tayjava"`)
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}
