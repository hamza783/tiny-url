package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hamza4253/tiny-url/shortener/internal/handler"
	"github.com/hamza4253/tiny-url/shortener/internal/repository"
	shorten "github.com/hamza4253/tiny-url/shortener/internal/service"
	"github.com/redis/go-redis/v9"
)

var (
	httpAddr  = getEnv("SHORTENING_SERVICE_ADDR", ":8081")
	redisAddr = getEnv("REDIS_URL", "localhost:6379")
)

func main() {
	log.Println("Starting URL Shortening Service: ", httpAddr)
	mux := http.NewServeMux()
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	// memRepo := repository.NewMemRepository()
	redisRepo := repository.NewRedisRepository(redisClient)
	service := shorten.NewShortenService(redisRepo)
	h := handler.NewHandler(service)
	h.RegisterRoutes(mux)

	server := http.Server{
		Addr:    httpAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}

// helper function to read env with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
