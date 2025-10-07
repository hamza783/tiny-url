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

func NewDBClient(ctx context.Context, postgresUrl string) (*DBClient, error) {
	pool, err := pgxpool.New(ctx, postgresUrl)
	if err != nil {
		return nil, err
	}

	return &DBClient{db: pool}, nil
}

func (c *DBClient) SaveURL(ctx context.Context, shortUrl, longUrl string) error {
	_, err := c.db.Exec(ctx, `
        INSERT INTO urls (short_url, long_url)
        VALUES ($1, $2)
    `, shortUrl, longUrl)
	return err
}

func (c *DBClient) SaveURLBatch(ctx context.Context, batchId, longUrl, shortUrl string) error {
	_, err := c.db.Exec(ctx, `
        UPDATE urls
				SET batch_id = $1
				WHERE short_url = $2 and long_url = $3
    `, batchId, shortUrl, longUrl)
	return err
}

func (c *DBClient) GetURLByShortURL(ctx context.Context, shortUrl string) (string, error) {
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

func (c *DBClient) GetURLByBatchId(ctx context.Context, batchId string) (map[string]string, error) {
	results := make(map[string]string)

	rows, err := c.db.Query(ctx, `select short_url, long_url from urls where batch_id = $1`, batchId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var shortURL, longURL string
		if err := rows.Scan(&shortURL, &longURL); err != nil {
			return nil, err
		}
		// Keep consistency with RedisRepository: map[longUrl]shortUrl
		results[longURL] = shortURL
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		log.Printf("No URLs found for batch_id %s", batchId)
	}

	return results, nil
}

func (c *DBClient) Close() {
	if c.db != nil {
		c.db.Close()
	}
}
