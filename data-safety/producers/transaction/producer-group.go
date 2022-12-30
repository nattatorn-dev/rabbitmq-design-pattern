package main

import (
	"log"
	"os"
	"strconv"
	"time"

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
		"transaction-message", // name
		"direct",              // type
		false,                 // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	round := string(os.Args[1])
	roundNumber, err := strconv.Atoi(round)

	if err != nil {
		failOnError(err, "round should be number")
	}

	log.Println(">>> round", roundNumber)

	err = ch.Tx()
	if err != nil {
		failOnError(err, "Failed to Start a transaction")
	}
	for i := 0; i < roundNumber; i++ {
		err = ch.Publish(
			"transaction-message", // exchange
			"",                    // routing key
			false,                 // mandatory
			false,                 // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(strconv.Itoa(i)),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf("Sent message round %d", i)
		time.Sleep(1 * time.Second)
	}
	err = ch.TxCommit()
	if err != nil {
		failOnError(err, "Failed to Commit the transaction")
	}

	log.Printf("Sent %s", round)
}
