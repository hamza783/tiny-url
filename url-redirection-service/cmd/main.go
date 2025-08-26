package main

import (
	"log"
	"net"
	"os"

	"github.com/hamza4253/tiny-url/redirect/internal/handler"
	"github.com/hamza4253/tiny-url/redirect/internal/repository"
	"github.com/hamza4253/tiny-url/redirect/internal/service"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

var (
	grpcAddr  = getEnv("REDIRECTION_SERVICE_ADDR", ":9000")
	redisAddr = getEnv("REDIS_URL", "localhost:6379")
)

func main() {
	// RPC server
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("Failed to start server on port %v. Error: %v", grpcAddr, err)
	}
	defer l.Close()

	// Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	// initialize grpc handler
	repo := repository.NewRedisRepository(redisClient)
	service := service.NewRedirectionService(repo)
	handler.NewGRPCRedirectionHandler(grpcServer, service)

	// start server
	log.Println("gRPC server listening on: ", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatalf("Failed to start gRPC server. Error: %v", err)
	}
}

// helper function to read env with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
