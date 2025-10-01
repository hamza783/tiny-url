package repository

import (
	"context"
	"errors"
	"log"
	"time"

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

func (r *RedisRepository) SaveURLBatch(ctx context.Context, batchId, longUrl, shortUrl string) error {
	key := "batch:" + batchId
	err := r.client.HSet(ctx, key, longUrl, shortUrl).Err()
	if err != nil {
		log.Println("Error saving to redis", err)
		return err
	}

	// set expire time so batch_id mapping doesn't stay forever
	err = r.client.Expire(ctx, key, 5*time.Minute).Err()
	if err != nil {
		log.Println("Error setting redis expire time", err)
		return err
	}

	return nil
}

func (r *RedisRepository) GetURLByBatchId(ctx context.Context, batchId string) (map[string]string, error) {
	key := "batch:" + batchId
	results, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return results, nil
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
