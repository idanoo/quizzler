package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	ctx    = context.Background()
)

const DefaultTTL = 3 * time.Hour

func Init() error {
	host := getEnv("REDIS_HOST", "redis")
	port := getEnv("REDIS_PORT", "6379")
	db, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
		DB:   db,
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	slog.Info("Connected to Redis")
	return nil
}

func Close() {
	if client != nil {
		client.Close()
	}
}

func Set(key string, value any) error {
	return SetWithTTL(key, value, DefaultTTL)
}

func SetWithTTL(key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.Set(ctx, key, data, ttl).Err()
}

func Get(key string, dest any) error {
	data, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

func Delete(key string) error {
	return client.Del(ctx, key).Err()
}

func DeletePattern(pattern string) error {
	keys, err := client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return client.Del(ctx, keys...).Err()
	}
	return nil
}

func Exists(key string) bool {
	result, err := client.Exists(ctx, key).Result()
	return err == nil && result > 0
}

// IsAvailable returns true if Redis is connected and available
func IsAvailable() bool {
	if client == nil {
		return false
	}
	return client.Ping(ctx).Err() == nil
}

// Increment increments a key and sets expiry if it's a new key
// Returns the new count after incrementing
func Increment(key string, ttl time.Duration) (int64, error) {
	if client == nil {
		return 0, fmt.Errorf("redis client not initialized")
	}

	count, err := client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// If this is the first increment (count == 1), set the expiry
	if count == 1 {
		client.Expire(ctx, key, ttl)
	}

	return count, nil
}

func UserKey(userID int) string {
	return fmt.Sprintf("user:%d", userID)
}

func DeckKey(deckID int) string {
	return fmt.Sprintf("deck:%d", deckID)
}

func UserDecksKey(userID int) string {
	return fmt.Sprintf("user:%d:decks", userID)
}

func DeckCardsKey(deckID int) string {
	return fmt.Sprintf("deck:%d:cards", deckID)
}

func CardKey(cardID int) string {
	return fmt.Sprintf("card:%d", cardID)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
