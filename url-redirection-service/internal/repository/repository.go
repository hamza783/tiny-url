package repository

import "context"

type Repository interface {
	LookupURL(ctx context.Context, shortUrl string) (string, error)
}
