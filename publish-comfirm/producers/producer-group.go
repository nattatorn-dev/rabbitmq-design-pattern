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
		"priority-message", // name
		"direct",           // type
		false,              // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	message := os.Args[1]
	log.Println(">>> message", message)

	priority := string(os.Args[2])
	priorityNumber, err := strconv.Atoi(priority)

	if err != nil {
		failOnError(err, "priority should be number")
		return
	}

	log.Println(">>> priority", priorityNumber)

	err = ch.Publish(
		"priority-message", // exchange
		"",                 // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
			Priority:    uint8(priorityNumber), // 0-255
		})
	failOnError(err, "Failed to publish a message")

	log.Printf("Sent %s", message)
}
