package database

import (
	"database/sql"
	"fmt"
	"ringer/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/lib/pq"
)

var jwtKey = []byte("my_secret_key")

// Claims represents the payload data inside the JWT token
type Claims struct {
	Username string `json:"username"` // You can include user information like username or user ID here
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token for the user.
func GenerateJWT(username string) (string, error) {
	// Set token expiration time (1 hour in this example)
	expirationTime := time.Now().Add(1 * time.Hour)

	// Create the JWT claims, which include the username and expiration time
	claims := &Claims{
		Username: username, // Add user-specific data here
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // Set token expiration time
		},
	}

	// Create the token using HS256 signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	// Return the generated token
	return tokenString, nil
}

func StoreToken(userID int, token string, expiresAt time.Time) error {
	query := `INSERT INTO user_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := DB.Exec(query, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("error storing token in database: %v", err)
	}
	return nil
}

func GetUserByToken(token string) (*models.User, error) {
	// Step 1: Validate the token exists and is not revoked
	var userID int
	query := `SELECT user_id FROM user_tokens WHERE token = $1 AND revoked = false AND expires_at > NOW()`
	err := DB.QueryRow(query, token).Scan(&userID)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("token not found or revoked/expired")
	} else if err != nil {
		return nil, fmt.Errorf("error querying token: %v", err)
	}

	// Step 2: Get the user details from the users table
	var user models.User
	query = `SELECT id, username, email FROM users WHERE id = $1`
	err = DB.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %v", err)
	}

	return &user, nil
}

func RevokeToken(token string) error {
	query := `UPDATE user_tokens SET revoked = true WHERE token = $1`
	_, err := DB.Exec(query, token)
	if err != nil {
		return fmt.Errorf("error revoking token: %v", err)
	}
	return nil
}

func ValidateToken(token string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM user_tokens WHERE token = $1 AND revoked = false AND expires_at > NOW())`
	err := DB.QueryRow(query, token).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error validating token: %v", err)
	}
	return exists, nil
}
