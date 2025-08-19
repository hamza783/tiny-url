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

func NewRedisRepository(c *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: c,
	}
}

func (r *RedisRepository) LookupURL(ctx context.Context, shortUrl string) (string, error) {
	val, err := r.client.Get(ctx, shortUrl).Result()
	if err == redis.Nil {
		log.Printf("No Long url exists in redis for %v", shortUrl)
		return "", errors.New("url not found")
	}
	return val, err
}
