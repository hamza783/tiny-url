package repository

import "context"

type Repository interface {
	SaveURL(ctx context.Context, shortUrl, longUrl string) error
	SaveURLBatch(ctx context.Context, batchId, longUrl, shortUrl string) error
	GetURLByShortURL(ctx context.Context, shortUrl string) (string, error)
	GetURLByBatchId(ctx context.Context, batchId string) (map[string]string, error)
}
