package handlers

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	_ "modernc.org/sqlite"
	"sms/db"
)

// temp in-memory db for testing
func setupTestDB(t *testing.T) *sql.DB {
	testDB, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	schema, err := os.ReadFile("../schema.sql")
	if err != nil {
		t.Fatalf("failed to read schema.sql: %v", err)
	}

	if _, err = testDB.Exec(string(schema)); err != nil {
		t.Fatalf("failed to execute schema: %v", err)
	}

	if err = db.SeedAdmin(testDB); err != nil {
		t.Fatalf("failed to seed admin: %v", err)
	}

	db.DB = testDB
	return testDB
}

func TestLoginSuccess(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	body := bytes.NewBufferString(`{"email":"admin@school.edu","password":"admin123"}`)
	req := httptest.NewRequest("POST", "/api/login", body)
	rr := httptest.NewRecorder()

	HandleLogin(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var hasSessionCookie bool
	for _, cookie := range rr.Result().Cookies() {
		if cookie.Name == "session_token" && cookie.Value != "" {
			hasSessionCookie = true
			break
		}
	}

	if !hasSessionCookie {
		t.Fatalf("expected session_token cookie to be set")
	}
}

func TestLoginWrongPassword(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	body := bytes.NewBufferString(`{"email":"admin@school.edu","password":"wrongpassword"}`)
	req := httptest.NewRequest("POST", "/api/login", body)
	rr := httptest.NewRecorder()

	HandleLogin(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", rr.Code)
	}
}

func TestRequireRoleMiddleware(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
	protected := RequireRole(dummyHandler, "admin")

	// test without session cookie
	reqNoCookie := httptest.NewRequest("GET", "/api/admin/only", nil)
	rrNoCookie := httptest.NewRecorder()
	protected(rrNoCookie, reqNoCookie)

	if rrNoCookie.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 without cookie, got %d", rrNoCookie.Code)
	}

	// create admin session in map
	adminToken := "test-admin-session-token"
	sessionMutex.Lock()
	sessions[adminToken] = Session{
		UserID:    1,
		Role:      "admin",
		FullName:  "System Administrator",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	sessionMutex.Unlock()

	// test with valid admin session cookie
	reqWithCookie := httptest.NewRequest("GET", "/api/admin/only", nil)
	reqWithCookie.AddCookie(&http.Cookie{Name: "session_token", Value: adminToken})
	rrWithCookie := httptest.NewRecorder()
	protected(rrWithCookie, reqWithCookie)

	if rrWithCookie.Code != http.StatusOK {
		t.Fatalf("expected status 200 with admin cookie, got %d", rrWithCookie.Code)
	}
}
