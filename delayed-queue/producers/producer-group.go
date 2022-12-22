package main

import (
	"log"
	"os"
	"strconv"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"delayed-message",   // name
		"x-delayed-message", // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		amqp.Table{
			"x-delayed-type": "direct",
		}, // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	message := os.Args[1]
	log.Println(">>> message", message)

	delay := string(os.Args[2])
	delayMilliSeconds, err := strconv.Atoi(delay)

	if err != nil {
		failOnError(err, "delay should be number")
	}

	log.Println(">>> delay", delayMilliSeconds, "ms")

	err = ch.Publish(
		"delayed-message", // exchange
		"",                // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			Headers: amqp.Table{
				"x-delay": delayMilliSeconds,
			},
			ContentType: "text/plain",
			Body:        []byte(message),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", message)
}
