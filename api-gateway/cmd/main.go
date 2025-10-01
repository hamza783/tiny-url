package main

import (
	"log"
	"net/http"
	"os"
	"time"

	handler "github.com/hamza4253/tiny-url/gateway/internal/handlers"
	"github.com/hamza4253/tiny-url/gateway/internal/publisher"
	shorten "github.com/hamza4253/tiny-url/gateway/internal/services"
	pb "github.com/hamza4253/tiny-url/shared/api/gen"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr                 = getEnv("API_GATEWAY", ":8080")
	urlShorteningServiceURL  = getEnv("URL_SHORTENING_URL", "http://localhost:8081")
	urlRedirectionServiceURL = getEnv("URL_REDIRECTION_GRPC_ADDR", "localhost:9000")
	QUEUE_NAME               = "shorten_url_batch"
	rabbitURL                = getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
)

func main() {
	log.Println("Starting API Gateway: ", httpAddr)

	// gRPC client for redirection service
	conn, err := grpc.NewClient(
		urlRedirectionServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to create a gRPC client: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewRedirectionServiceClient(conn)
	log.Println("Connected to Redirection service at", urlRedirectionServiceURL)

	// REST client for shortening service
	urlShorteningService := shorten.NewURLShorteningClient(urlShorteningServiceURL)
	log.Println("Connected to Shortening service at", urlShorteningServiceURL)

	// publisher
	pubConn, err := connectToMQ()
	failOnError("Failed to connect to RabbitMQ", err)
	defer pubConn.Close()

	p, err := publisher.NewShorteningPublisher(pubConn, QUEUE_NAME)
	failOnError("Failed to connect to RabbitMQ", err)
	defer p.Close()
	log.Println("Connected to rabbit mq queue", QUEUE_NAME)

	// Handler
	h := handler.NewHandler(urlShorteningService, grpcClient, p)

	mux := http.NewServeMux()
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

func connectToMQ() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	for i := 1; i <= 30; i++ { // ~60s total with 2s sleep
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("RabbitMQ not ready (%v). Retry %d/30 in 2s...", err, i)
		time.Sleep(2 * time.Second)
	}
	return conn, err
}

func failOnError(msg string, err error) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
