package repository

import "context"

type Repository interface {
	SaveURL(ctx context.Context, shortUrl string, longUrl string) error
	GetURLByShortURL(ctx context.Context, shortUrl string) (string, error)
}
