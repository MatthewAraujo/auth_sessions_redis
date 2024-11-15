package limit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/MatthewAraujo/auth-limit-redis/internal/types"
	"github.com/redis/go-redis/v9"
)

var ErrLimitRequest = fmt.Errorf("LIMIT REQUEST")
var ChargePrice = 1.4

type TokenStore struct {
	db *redis.Client
}

func NewTokenStore(db *redis.Client) *TokenStore {
	return &TokenStore{
		db: db,
	}
}

// IncrementTokenCount increments the count for a given token
func (store *TokenStore) IncrementTokenCount(token string) error {
	ctx := context.Background()
	key := "token:count:" + token

	_, err := store.db.Incr(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to increment token count: %w", err)
	}

	tokenValue, err := store.GetTokenCount(token)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if tokenValue > 8 {
		return ErrLimitRequest
	}

	return nil
}

// GetTokenCount retrieves the count for a given token
func (store *TokenStore) GetTokenCount(token string) (int64, error) {
	ctx := context.Background()
	key := "token:count:" + token

	count, err := store.db.Get(ctx, key).Int64()
	if err != nil && err != redis.Nil {
		return 0, fmt.Errorf("failed to get token count: %w", err)
	}

	return count, nil
}

func (store *TokenStore) TokenIsExpired(token string) (bool, error) {
	ctx := context.Background()
	key := "token:" + token

	tokenDataStr, err := store.db.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return false, fmt.Errorf("error fetching token data: %v", err)
	}

	if err == redis.Nil {
		log.Print("token expirado?")
		return true, nil
	}

	var tokenData types.TokenData
	err = json.Unmarshal([]byte(tokenDataStr), &tokenData)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling token data: %v", err)
	}

	if time.Now().After(tokenData.LoginTime.Add(48 * time.Hour)) {
		log.Printf("token expirou")
		return true, nil
	}

	return false, nil
}
