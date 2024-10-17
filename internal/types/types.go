package types

import "time"

type User struct {
	Username   string
	Password   string
	Token      string
	LoginTime  time.Time
	Expiration time.Time
}

type Login struct {
	Username string
	Password string
}

// LoginStore defines methods for user management.
type LoginStore interface {
	CreateUser(username, password string) error // Create a new user
	LoginUser(login *Login) (string, error)     // Log in a user and return a token
}

type LimitStore interface {
	IncrementTokenCount(token string) error
	GetTokenCount(token string) (int64, error)
}
