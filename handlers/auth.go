package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
	"sms/db"
)

// represents an active user login session
type Session struct {
	UserID    int       `json:"user_id"`
	Role      string    `json:"role"`
	FullName  string    `json:"full_name"`
	ExpiresAt time.Time `json:"expires_at"`
}

var (
	sessions     = make(map[string]Session)
	sessionMutex sync.RWMutex
)

// generate 32 character hex token
func generateToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// auth user and start session
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token := generateToken()
	expiresAt := time.Now().Add(24 * time.Hour)
	session := Session{
		UserID:    int(user.ID),
		Role:      user.Role,
		FullName:  user.FullName,
		ExpiresAt: expiresAt,
	}

	sessionMutex.Lock()
	sessions[token] = session
	sessionMutex.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
	})

	writeJSON(w, http.StatusOK, map[string]any{
		"id":        user.ID,
		"full_name": user.FullName,
		"role":      user.Role,
	})
}

// end session and clear cookie
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		sessionMutex.Lock()
		delete(sessions, cookie.Value)
		sessionMutex.Unlock()
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "logged out",
	})
}

// return current user session
func HandleMe(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	sessionMutex.RLock()
	session, exists := sessions[cookie.Value]
	sessionMutex.RUnlock()

	if !exists {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if time.Now().After(session.ExpiresAt) {
		sessionMutex.Lock()
		delete(sessions, cookie.Value)
		sessionMutex.Unlock()
		writeError(w, http.StatusUnauthorized, "session expired")
		return
	}

	writeJSON(w, http.StatusOK, session)
}
