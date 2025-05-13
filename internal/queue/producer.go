package queue

import (
	"encoding/json"

	"github.com/streadway/amqp"
)

type Producer interface {
	PublishTransaction(message interface{}) error
	Close() error
}

type RabbitMQProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

// creates a new RabbitMQ producer
func NewRabbitMQProducer(url, queueName string) (*RabbitMQProducer, error) {
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
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQProducer{
		conn:    conn,
		channel: ch,
		queue:   queueName,
	}, nil
}

// publishes a transaction message to the queue
func (p *RabbitMQProducer) PublishTransaction(message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",      // exchange
		p.queue, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Ensure message persistence
		})
}

// Close closes the connection and channel
func (p *RabbitMQProducer) Close() error {
	if err := p.channel.Close(); err != nil {
		return err
	}
	return p.conn.Close()
}
