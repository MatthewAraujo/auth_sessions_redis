package limit

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ErrLimitRequest = fmt.Errorf("LIMIT REQUEST")

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
