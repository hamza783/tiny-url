package publisher

import (
	"context"
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type ShorteningPublisher struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

// url_shortening_queue
func NewShorteningPublisher(conn *amqp.Connection, queueName string) (*ShorteningPublisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// create queue if doesn't exist
	// params must be same if declared elsewhere(probably in consumer)
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

	return &ShorteningPublisher{
		conn:      conn,
		channel:   ch,
		queueName: q.Name,
	}, nil
}

func (p *ShorteningPublisher) Publish(ctx context.Context, batchId, longUrl string) error {
	payload := struct {
		URL     string `json:"url"`
		BatchId string `json:"batch_id"`
	}{
		URL:     longUrl,
		BatchId: batchId,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("failed to marshal payload: %v", err)
		return err
	}

	// publish the message
	err = p.channel.Publish(
		"",          // exchange
		p.queueName, // routing key = queue name
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}

func (p *ShorteningPublisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
}
