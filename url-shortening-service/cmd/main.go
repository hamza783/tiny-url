package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hamza4253/tiny-url/shortener/internal/consumers"
	"github.com/hamza4253/tiny-url/shortener/internal/handler"
	"github.com/hamza4253/tiny-url/shortener/internal/repository"
	shorten "github.com/hamza4253/tiny-url/shortener/internal/service"
	"github.com/streadway/amqp"
)

// TODO: Move everything to env
var (
	httpAddr = getEnv("SHORTENING_SERVICE_ADDR", ":8081")
	// redisAddr   = getEnv("REDIS_URL", "localhost:6379")
	QUEUE_NAME  = "shorten_url_batch"
	rabbitURL   = getEnv("RABBITMQ_URL", "amqp://guest:guest@rabbitmq:5672/")
	postgresURL = getEnv("DATABASE_URL", "")
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Starting URL Shortening Service: ", httpAddr)
	mux := http.NewServeMux()
	// Create redis client
	// redisClient := redis.NewClient(&redis.Options{
	// 	Addr: redisAddr,
	// })

	// Create DB
	// memRepo := repository.NewMemRepository()
	// redisRepo := repository.NewRedisRepository(redisClient)
	postgresRepo, err := repository.NewDBClient(ctx, postgresURL)
	failOnError("Failed to connect to rabbit mq", err)
	defer postgresRepo.Close()

	service := shorten.NewShortenService(postgresRepo)

	// Create handler and register routes
	h := handler.NewHandler(service)
	h.RegisterRoutes(mux)

	// Connect to RabbitMQ
	conn, err := connectToMQ()
	failOnError("Failed to connect to rabbit mq", err)
	defer conn.Close()

	// Create a consumer
	consumer, err := consumers.NewConsumer(conn, service, postgresRepo, QUEUE_NAME)
	failOnError("Error setting up a RabbitMQ consumer", err)

	// Start consumer
	err = consumer.Start(ctx)
	failOnError("Error starting the RabbitMQ consumer", err)
	defer consumer.Close()

	// start server
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
		log.Fatalf("%s: %v", msg, err)
	}
}
