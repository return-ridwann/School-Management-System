package db

import (
	"database/sql"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	schema, err := os.ReadFile("../schema.sql")
	if err != nil {
		t.Fatalf("failed to read schema.sql: %v", err)
	}

	if _, err = db.Exec(string(schema)); err != nil {
		t.Fatalf("failed to execute schema: %v", err)
	}

	return db
}

func TestSeedAdmin(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	// first seed: should create admin
	if err := SeedAdmin(testDB); err != nil {
		t.Fatalf("SeedAdmin failed: %v", err)
	}

	var fullName, email, passwordHash, role string
	err := testDB.QueryRow("SELECT full_name, email, password_hash, role FROM users WHERE role = 'admin'").
		Scan(&fullName, &email, &passwordHash, &role)
	if err != nil {
		t.Fatalf("failed to query admin user: %v", err)
	}

	if fullName != "System Administrator" {
		t.Errorf("expected full_name 'System Administrator', got '%s'", fullName)
	}
	if email != "admin@school.edu" {
		t.Errorf("expected email 'admin@school.edu', got '%s'", email)
	}
	if role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", role)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte("admin123")); err != nil {
		t.Errorf("password hash does not match 'admin123': %v", err)
	}

	// second seed: should be idempotent and do nothing
	if err := SeedAdmin(testDB); err != nil {
		t.Fatalf("second SeedAdmin call failed: %v", err)
	}

	var count int
	if err := testDB.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin'").Scan(&count); err != nil {
		t.Fatalf("failed to count admins: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 admin user, got %d", count)
	}
}
