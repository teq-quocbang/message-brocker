package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/teq-quocbang/message-brocker/state"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

type OrderServiceMessage struct {
	OrderName string      `json:"order_name"`
	OrderType string      `json:"order_type"`
	Qty       int         `json:"qty"`
	State     state.State `json:"state"`
	OrderTime time.Time   `json:"order_time"`
}

func main() {
	conn, err := amqp.Dial("amqp://quocbang:quocbang@localhost:5672/")
	failOnError(err, "Failed to connect RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("order_service", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	orderMessage := OrderServiceMessage{
		OrderName: gofakeit.Name(),
		OrderType: gofakeit.CreditCard().Type,
		Qty:       2,
		State:     state.State_ORDER,
		OrderTime: time.Now(),
	}
	body, err := json.Marshal(orderMessage)
	failOnError(err, "Failed to marshal body")

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
