package consumers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hamza4253/tiny-url/shortener/internal/repository"
	shorten "github.com/hamza4253/tiny-url/shortener/internal/service"
	"github.com/streadway/amqp"
)

var CONSUMER_NAME = "shorten_url_consumer"

type Consumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	service   *shorten.ShortenService
	repo      repository.Repository
	queueName string
}

func NewConsumer(conn *amqp.Connection, svc *shorten.ShortenService, r repository.Repository, queueName string) (*Consumer, error) {
	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// create queue if doesn't exist
	// params must be same if declared elsewhere(probably in producer)
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		conn:      conn,
		channel:   ch,
		service:   svc,
		repo:      r,
		queueName: q.Name,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Printf("=======Starting consumer %s for queue %s=============\n", CONSUMER_NAME, c.queueName)
	msgs, err := c.channel.Consume(
		c.queueName,   // queue
		CONSUMER_NAME, // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
			// unmarshal message
			var payload struct {
				URL     string `json:"url"`
				BatchId string `json:"batch_id"`
			}
			if err := json.Unmarshal(msg.Body, &payload); err != nil {
				log.Printf("failed to unmarshal message: %v", err)
				continue
			}
			longUrl := payload.URL
			// shorten url
			shortUrl, err := c.service.Shorten(ctx, longUrl)
			if err != nil {
				log.Printf("failed to shorten URL %s: %v", longUrl, err)
				continue
			}
			// save shortened url list to redis
			c.repo.SaveURLBatch(ctx, payload.BatchId, longUrl, shortUrl)
			log.Println("done saving url batch")
		}
	}()

	log.Printf("Waiting for messages.")

	return nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
}
