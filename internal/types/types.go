package types

import "time"

type User struct {
	Username   string
	Password   string
	Token      string
	LoginTime  time.Time
	Expiration time.Time
}

type TokenData struct {
	LoginTime  time.Time `json:"login_time"`
	Expiration time.Time `json:"expiration"`
}

type Login struct {
	Username string
	Password string
}

type LoginStore interface {
	CreateUser(username, password string) error
	LoginUser(login *Login) (string, error)
}

type LimitStore interface {
	IncrementTokenCount(token string) error
	GetTokenCount(token string) (int64, error)
	TokenIsExpired(token string) (bool, error)
}
