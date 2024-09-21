package database

import (
	"database/sql"
	"fmt"
	"ringer/models"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AuthenticateUser(username string, password string) (*models.User, error) {

	query := `SELECT id, username, email, password FROM users WHERE name = $1`

	row := DB.QueryRow(query, username)
	var user models.User

	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {

			return nil, fmt.Errorf("user not found")
		}

		return nil, fmt.Errorf("error querying user: %v", err)

	}

	return &user, nil

}

func Login(username, password string) (string, error) {
	// Authenticate user
	user, err := AuthenticateUser(username, password)
	if err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := GenerateJWT(user.Username)
	if err != nil {
		return "", fmt.Errorf("error generating token: %v", err)
	}

	// Store the token in the database
	expiresAt := time.Now().Add(1 * time.Hour) // Set token expiration time
	err = StoreToken(user.ID, token, expiresAt)
	if err != nil {
		return "", fmt.Errorf("error storing token: %v", err)
	}

	return token, nil
}
