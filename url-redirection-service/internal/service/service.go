package service

import (
	"context"
	"log"

	repository "github.com/hamza4253/tiny-url/redirect/internal/repository"
)

type RedirectionService struct {
	repo repository.Repository
}

func NewRedirectionService(r repository.Repository) *RedirectionService {
	return &RedirectionService{
		repo: r,
	}
}

func (s *RedirectionService) GetURLByShortURL(ctx context.Context, shortUrl string) (string, error) {
	longUrl, err := s.repo.LookupURL(ctx, shortUrl)
	if err != nil {
		log.Printf("Error getting long url for %v. Error: %v", shortUrl, err)
		return "", err
	}
	return longUrl, nil
}
