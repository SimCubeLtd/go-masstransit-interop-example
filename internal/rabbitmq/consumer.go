package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"go-masstransit-interop/internal/masstransit"
	"go-masstransit-interop/internal/worker"
)

type (
	Consumer struct {
		conn       *amqp091.Connection
		channel    *amqp091.Channel
		queueName  string
		workerPool *worker.Pool
	}
)

func NewConsumer(amqpURL, queueName string, workerPool *worker.Pool) (*Consumer, error) {
	var conn, err = amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		err := conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}

	return &Consumer{
		conn:       conn,
		channel:    ch,
		queueName:  queueName,
		workerPool: workerPool,
	}, nil
}

func (c *Consumer) StartConsuming(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return err
	}

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				// Channel closed
				return nil
			}
			var envelope masstransit.Envelope
			if err := json.Unmarshal(d.Body, &envelope); err != nil {
				log.Printf("Error decoding JSON: %s", err)
				continue
			}

			if err := masstransit.ProcessMessage(envelope, c.workerPool); err != nil {
				log.Printf("Error processing message: %s", err)
			}
		case <-ctx.Done():
			// Context cancelled
			return nil
		}
	}
}

func (c *Consumer) Close() {
	if err := c.channel.Close(); err != nil {
		log.Printf("Error closing channel: %s", err)
	}
	if err := c.conn.Close(); err != nil {
		log.Printf("Error closing connection: %s", err)
	}
}
