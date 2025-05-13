package queue

import (
	"context"
	"log"

	"github.com/streadway/amqp"
)

type MessageHandler func([]byte) error

type Consumer interface {
	Consume(ctx context.Context, handler MessageHandler) error
	Close() error
}

// Implements the Consumer interface with RabbitMQ
type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

// creates a new RabbitMQ consumer
func NewRabbitMQConsumer(url, queueName string) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare the queue
	_, err = ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	// Set QoS for fair dispatching
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}

// starts consuming messages from the queue
func (c *RabbitMQConsumer) Consume(ctx context.Context, handler MessageHandler) error {
	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Consumer context cancelled, stopping...")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("Channel closed, stopping consumer...")
					return
				}

				if err := handler(msg.Body); err != nil {
					log.Printf("Error processing message: %v", err)
					// Nack the message and requeue it
					if err := msg.Nack(false, true); err != nil {
						log.Printf("Error nacking message: %v", err)
					}
				} else {
					// Acknowledge the message
					if err := msg.Ack(false); err != nil {
						log.Printf("Error acknowledging message: %v", err)
					}
				}
			}
		}
	}()

	log.Printf("Started consuming from queue: %s", c.queue)
	return nil
}

// Close closes the connection and channel
func (c *RabbitMQConsumer) Close() error {
	if err := c.channel.Close(); err != nil {
		return err
	}
	return c.conn.Close()
}
