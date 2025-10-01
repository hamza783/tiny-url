package repository

import (
	"context"
	"errors"
	"fmt"
)

type MemRepository struct {
	store map[string]string
}

// in memory storage just for testing. Create and get url
func NewMemRepository() *MemRepository {
	return &MemRepository{
		store: make(map[string]string),
	}
}

func (r *MemRepository) SaveURL(ctx context.Context, shortUrl string, longUrl string) error {
	fmt.Println("Saving url ====>")
	r.store[shortUrl] = longUrl
	return nil
}

func (r *MemRepository) SaveURLBatch(ctx context.Context, batchId, longUrl, shortUrl string) error {
	// TODO
	return nil
}

func (r *MemRepository) GetURLByBatchId(ctx context.Context, batchId string) (map[string]string, error) {
	// TODO
	return nil, nil
}

// get long url give short url. just for testing
func (r *MemRepository) GetURLByShortURL(ctx context.Context, shortUrl string) (string, error) {
	fmt.Println("Getting long url for short url ====>", shortUrl)
	longUrl := r.store[shortUrl]
	fmt.Println("Long url ====>", longUrl)
	if longUrl == "" {
		return "", errors.New("no long url found with this short url")
	}
	return longUrl, nil
}
