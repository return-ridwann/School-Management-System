package db

import (
	"testing"
)

func TestGetUserByEmail(t *testing.T) {
	testDB := setupTestDB(t)
	defer testDB.Close()

	// temporarily assign global DB to testDB
	oldDB := DB
	DB = testDB
	defer func() { DB = oldDB }()

	if err := SeedAdmin(testDB); err != nil {
		t.Fatalf("failed to seed admin: %v", err)
	}

	user, err := GetUserByEmail("admin@school.edu")
	if err != nil {
		t.Fatalf("GetUserByEmail returned error: %v", err)
	}

	if user.Email != "admin@school.edu" {
		t.Errorf("expected email 'admin@school.edu', got '%s'", user.Email)
	}
	if user.Role != "admin" {
		t.Errorf("expected role 'admin', got '%s'", user.Role)
	}
	if user.FullName != "System Administrator" {
		t.Errorf("expected full_name 'System Administrator', got '%s'", user.FullName)
	}

	// test non-existent user
	nonExistent, err := GetUserByEmail("nonexistent@school.edu")
	if err == nil {
		t.Fatalf("expected error for non-existent user, got user: %v", nonExistent)
	}
}
