package main

import (
	"log"
	"net/http"
	"os"
	"time"

	handler "github.com/hamza4253/tiny-url/gateway/internal/handlers"
	shorten "github.com/hamza4253/tiny-url/gateway/internal/services"
	pb "github.com/hamza4253/tiny-url/shared/api/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr                 = getEnv("API_GATEWAY", ":8080")
	urlShorteningServiceURL  = getEnv("URL_SHORTENING_URL", "http://localhost:8081")
	urlRedirectionServiceURL = getEnv("URL_REDIRECTION_GRPC_ADDR", "localhost:9000")
)

func main() {
	log.Println("Starting API Gateway: ", httpAddr)
	log.Println("Starting urlShorteningServiceURL: ", urlShorteningServiceURL)
	log.Println("Starting urlRedirectionServiceURL: ", urlRedirectionServiceURL)

	conn, err := grpc.NewClient(
		urlRedirectionServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to create a gRPC client: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to Redirection service at", urlRedirectionServiceURL)
	client := pb.NewRedirectionServiceClient(conn)

	mux := http.NewServeMux()
	urlShorteningService := shorten.NewURLShorteningClient(urlShorteningServiceURL)
	h := handler.NewHandler(urlShorteningService, client)
	h.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         httpAddr,
		Handler:      mux,
		IdleTimeout:  120 * time.Second, // max time to wait for next request. used if multiple requests from same client
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
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
