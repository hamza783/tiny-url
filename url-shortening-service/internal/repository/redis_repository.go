package repository

import (
	"context"
	"errors"
	"log"

	"github.com/redis/go-redis/v9"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

func (r *RedisRepository) SaveURL(ctx context.Context, shortUrl string, longURL string) error {
	return r.client.Set(ctx, shortUrl, longURL, 0).Err()
}

// GetLongURL retrieves the original URL for the given shortCode.
// Returns an error if not found.
func (r *RedisRepository) GetURLByShortURL(ctx context.Context, shortUrl string) (string, error) {
	val, err := r.client.Get(ctx, shortUrl).Result()
	if err == redis.Nil {
		log.Printf("No Long url exists in redis for %v", shortUrl)
		return "", errors.New("url not found")
	}
	return val, err
}
