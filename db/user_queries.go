package db

import (
	"sms/models"
)

// user record by email
func GetUserByEmail(email string) (*models.User, error) {
	query := "SELECT id, full_name, email, password_hash, role, created_at FROM users WHERE email = ?"
	var user models.User
	err := DB.QueryRow(query, email).Scan(&user.ID, &user.FullName, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
