package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type RabbitMQAdapter struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewRabbitMQAdapter(amqpURL, queueName string) (*RabbitMQAdapter, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQAdapter{
		conn:    conn,
		channel: ch,
		queue:   q,
	}, nil
}

func (r *RabbitMQAdapter) PublishMessage(body string) error {
	err := r.channel.Publish(
		"",
		r.queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	if err != nil {
		return err
	}
	log.Printf("Message published: %s", body)
	return nil
}

func (r *RabbitMQAdapter) ConsumeMessages(handler func(string)) {
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for m := range msgs {
			handler(string(m.Body))
		}
	}()
}
