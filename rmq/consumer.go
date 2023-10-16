package rmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zura-t/go_delivery_system-accounts/cmd/api"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	server  *api.Server
}

func NewConsumer(conn *amqp.Connection, channel *amqp.Channel, server *api.Server) (Consumer, error) {
	consumer := Consumer{
		conn:    conn,
		channel: channel,
		server:  server,
	}

	return consumer, nil
}

func (consumer *Consumer) Listen(topics []string) error {
	err := consumer.ListenCreateUserRPC()
	if err != nil {
		return err
	}

	return nil
}

func (consumer *Consumer) ListenCreateUserRPC() error {
	q, err := consumer.declareQueue("queue.user.create")
	if err != nil {
		return err
	}

	messages, err := consumer.ConsumeRPCMessages(q.Name)

	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			resp, err := consumer.HandleRPC(context.Background(), d.Body)
			if err != nil {
				res, err := json.Marshal(&Response{Data: "error"})
				err = consumer.publishRPCResponse(d.ReplyTo, d.CorrelationId, res)

				if err != nil {
					fmt.Printf("Can't send event :", err)
					d.Ack(false)
					return
				}
				d.Ack(false)
				return
			}
			response, err := json.Marshal(&resp)
			log.Println(resp)

			err = consumer.publishRPCResponse(d.ReplyTo, d.CorrelationId, response)

			if err != nil {
				fmt.Printf("Can't send event :", err)
				d.Ack(false)
				return
			}

			d.Ack(false)
		}
	}()

	fmt.Printf("Wait for message on [Exchange, Queue ] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func (consumer *Consumer) publishRPCResponse(replyTo string, correlationId string, response []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return consumer.channel.PublishWithContext(ctx,
		"",      // exchange
		replyTo, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: correlationId,
			Body:          []byte(response),
		},
	)
}

func (consumer *Consumer) declareQueue(queue string) (*amqp.Queue, error) {
	q, err := consumer.channel.QueueDeclare(
		queue, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (consumer *Consumer) ConsumeRPCMessages(queue string) (<-chan amqp.Delivery, error) {
	messages, err := consumer.channel.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
