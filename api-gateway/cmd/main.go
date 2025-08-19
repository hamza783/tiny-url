package main

import (
	"log"
	"net/http"

	handler "github.com/hamza4253/tiny-url/gateway/internal/handlers"
	shorten "github.com/hamza4253/tiny-url/gateway/internal/services"
	pb "github.com/hamza4253/tiny-url/shared/api/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr                 = ":8080"
	urlShorteningServiceURL  = "http://localhost:8081"
	urlRedirectionServiceURL = "localhost:9000"
)

func main() {
	log.Println("Starting API Gateway: ", httpAddr)

	conn, err := grpc.NewClient(urlRedirectionServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		Addr:    httpAddr,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
