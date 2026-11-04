package db

import (
	"database/sql"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// creates a default admin user if one does not exist
func SeedAdmin(database *sql.DB) error {
	var count int
	query := "SELECT COUNT(*) FROM users WHERE role = 'admin'"
	err := database.QueryRow(query).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing admin: %w", err)
	}

	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash default admin password: %w", err)
	}

	insertSQL := `INSERT INTO users (full_name, email, password_hash, role) VALUES (?, ?, ?, ?)`
	_, err = database.Exec(insertSQL, "System Administrator", "admin@school.edu", string(hash), "admin")
	if err != nil {
		return fmt.Errorf("failed to insert default admin: %w", err)
	}

	log.Println("Default admin account created: admin@school.edu / admin123")
	return nil
}
