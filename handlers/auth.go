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

// get session from request cookie and check validity
func GetSessionFromRequest(r *http.Request) (Session, bool) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return Session{}, false
	}

	sessionMutex.RLock()
	sess, exists := sessions[cookie.Value]
	sessionMutex.RUnlock()

	if !exists {
		return Session{}, false
	}

	if time.Now().After(sess.ExpiresAt) {
		sessionMutex.Lock()
		delete(sessions, cookie.Value)
		sessionMutex.Unlock()
		return Session{}, false
	}

	return sess, true
}

// check session and role permissions
func RequireRole(next http.HandlerFunc, allowedRoles ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, ok := GetSessionFromRequest(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "Unauthorized: please log in")
			return
		}

		for _, role := range allowedRoles {
			if sess.Role == role {
				next(w, r)
				return
			}
		}

		writeError(w, http.StatusForbidden, "Forbidden: insufficient permissions")
	}
}

// test endpoint for admin role
func HandleAdminTest(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "Admin access granted"})
}

// test endpoint for teacher role
func HandleTeacherTest(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"message": "Teacher access granted"})
}

// return current user session
func HandleMe(w http.ResponseWriter, r *http.Request) {
	sess, ok := GetSessionFromRequest(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	writeJSON(w, http.StatusOK, sess)
}
