package models

type User struct {
	ID       int
	Username string
	Email    string
	Password string // hashed password
}
