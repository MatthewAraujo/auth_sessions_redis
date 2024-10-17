package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/MatthewAraujo/auth-limit-redis/internal/types"
	"github.com/redis/go-redis/v9"
)

type RedisUserStore struct {
	db *redis.Client
}

func NewRedisUserStore(db *redis.Client) *RedisUserStore {
	return &RedisUserStore{
		db: db,
	}
}

func (s *RedisUserStore) CreateUser(username, password string) error {
	userPrefix := "session:"
	var ctx = context.Background()

	userKey := userPrefix + username

	err := s.db.HSet(ctx, userKey, map[string]interface{}{
		"username":   username,
		"password":   password,
		"loginTime":  time.Now().Format(time.RFC3339),
		"expiration": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}).Err()

	if err != nil {
		log.Fatal("Error storing user in Redis:", err)

		return err
	}

	log.Printf("User created successfully!")
	return nil
}

func (s *RedisUserStore) LoginUser(login *types.Login) (string, error) {
	username := login.Username
	password := login.Password

	userPrefix := "session:"

	var ctx = context.Background()

	keys, err := s.db.Keys(ctx, userPrefix+"*").Result()
	if err != nil {
		log.Fatal("Error fetching keys:", err)
		return "", err
	}

	for _, key := range keys {
		userData, err := s.db.HGetAll(ctx, key).Result()
		if err != nil {
			log.Fatal("Error retrieving user from Redis:", err)
			return "", err
		}

		if userData["username"] == username && userData["password"] == password {
			token := GenerateSecureToken(32) // Generate a 32-byte token

			err = s.db.Set(ctx, token, key, 24*time.Hour).Err() // Store the token with a 24-hour expiration
			if err != nil {
				log.Fatal("Error storing token in Redis:", err)
				return "", err
			}

			log.Printf("User found and token generated!")
			return token, nil
		}
	}

	log.Fatal("User not found or invalid credentials.")
	return "", fmt.Errorf("user not found or invalid credentials")
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
