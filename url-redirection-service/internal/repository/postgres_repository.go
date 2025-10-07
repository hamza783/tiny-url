package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBClient struct {
	db *pgxpool.Pool
}

func NewDBClient(ctx context.Context, dbUrl string) (*DBClient, error) {
	pool, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		return nil, err
	}

	return &DBClient{db: pool}, nil
}

func (c *DBClient) LookupURL(ctx context.Context, shortUrl string) (string, error) {
	var longUrl string
	err := c.db.QueryRow(ctx, `select long_url from urls where short_url = $1`, shortUrl).Scan(&longUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Printf("No Long url exists in redis for %v", shortUrl)
			return "", errors.New("url not found")
		}
		return "", err
	}

	return longUrl, nil
}

func (c *DBClient) Close() {
	if c.db != nil {
		c.db.Close()
	}
}
