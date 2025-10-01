package shorten

import (
	"context"
	"log"

	"github.com/hamza4253/tiny-url/shortener/internal/repository"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type ShortenService struct {
	repo repository.Repository
}

func NewShortenService(repo repository.Repository) *ShortenService {
	return &ShortenService{
		repo: repo,
	}
}

func (s *ShortenService) Shorten(ctx context.Context, longUrl string) (string, error) {
	shortUrl, err := createShortRandomUrl()
	if err != nil {
		log.Printf("Error creating a short URL for %v. Error: %v", longUrl, err)
		return "", err
	}

	err = s.repo.SaveURL(ctx, shortUrl, longUrl)
	if err != nil {
		log.Printf("Error saving short url for %v Error: %v", shortUrl, err)
		return "", err
	}
	return shortUrl, nil
}

func (s *ShortenService) GetFullURL(ctx context.Context, shortUrl string) (string, error) {
	longUrl, err := s.repo.GetURLByShortURL(ctx, shortUrl)
	if err != nil {
		log.Printf("Error generating long url for %v. Error: %v", shortUrl, err)
		return "", err
	}
	return longUrl, nil
}

func (s *ShortenService) GetUrlsByBatchId(ctx context.Context, batchId string) (map[string]string, error) {
	urlsMap, err := s.repo.GetURLByBatchId(ctx, batchId)
	if err != nil {
		log.Printf("Error generating urls for batch id: %v. Error: %v", batchId, err)
		return nil, err
	}
	return urlsMap, nil
}

func createShortRandomUrl() (string, error) {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	id, err := gonanoid.Generate(alphabet, 6)
	if err != nil {
		return "", err
	}

	return id, nil
}
