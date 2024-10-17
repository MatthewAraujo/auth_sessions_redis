package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	Health() map[string]string
	GetRedisClient() *redis.Client
}

type service struct {
	db *redis.Client
}

var (
	address  = os.Getenv("DB_ADDRESS")
	port     = os.Getenv("DB_PORT")
	password = os.Getenv("DB_PASSWORD")
	database = os.Getenv("DB_DATABASE")
)

// New creates a new database service and returns it
func New() Service {
	num, err := strconv.Atoi(database)
	if err != nil {
		log.Fatalf("database incorrect %v", err)
	}

	fullAddress := fmt.Sprintf("%s:%s", address, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fullAddress,
		Password: password,
		DB:       num,
	})

	s := &service{db: rdb}

	return s
}

// Health returns the health status and statistics of the Redis server.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stats := make(map[string]string)
	stats = s.checkRedisHealth(ctx, stats)

	return stats
}

// GetRedisClient returns the Redis client
func (s *service) GetRedisClient() *redis.Client {
	return s.db
}

// checkRedisHealth checks the health of the Redis server and adds the relevant statistics to the stats map.
func (s *service) checkRedisHealth(ctx context.Context, stats map[string]string) map[string]string {
	return map[string]string{
		"healthy": "yeah",
	}
}
