package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

type loginIn struct {
	UserID   int64  `json:"userId"`
	Password string `json:"password"` // lab-only fake check
}

type loginOut struct {
	Token string `json:"token"`
}

// LoginHandler issues a 15-minute JWT for a positive userId when password is "lab".
func (a *Auth) LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var in loginIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil || in.UserID <= 0 {
		http.Error(w, "invalid login", http.StatusBadRequest)
		return
	}
	if in.Password != "lab" {
		unauthorized(w)
		return
	}
	tok, err := a.Issue(in.UserID, 15*time.Minute)
	if err != nil {
		http.Error(w, "auth misconfigured", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(loginOut{Token: tok})
}
