package rmq

import (
	"context"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func (e *Emitter) setup() error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return nil
}

func (e *Emitter) SendRPC(routingKey string, message string) error {
	ch, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"users", // name
		false,   // durable
		false,   // delete when unused
		true,    // exclusive
		false,   // noWait
		nil,     // arguments
	)

	corrId, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",            // exchange
		"user.create", // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId.String(),
			ReplyTo:       q.Name,
			Body:          []byte(message),
		})

	return nil
}
