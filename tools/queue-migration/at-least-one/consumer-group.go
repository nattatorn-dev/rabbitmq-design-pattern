package main

import (
	"log"
	"os"

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

	// consume
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	queueName := string(os.Args[1])
	log.Println(">>> Source Queue Name", queueName)

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	// publisher
	backupExchangeName := "backup_" + q.Name
	err = ch.ExchangeDeclare(
		backupExchangeName, // name
		"direct",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	targetQueue, err := ch.QueueDeclare(
		"backup_"+q.Name, // name
		true,             // durable
		false,            // delete when unused
		false,            // exclusive
		false,            // no-wait
		nil,              // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		targetQueue.Name,   // queue name
		"",                 // routing key
		backupExchangeName, // exchange
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to bind a queue")

	err = ch.Confirm(false)
	if err != nil {
		failOnError(err, "Failed to enable publisher confirms on the channel")
	}
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	forever := make(chan bool)

	// Consume and publish messages one by one
	go func() {
		for d := range msgs {
			// consumer
			log.Printf("Received Message %s", d.Body)

			// publish
			err = ch.Publish(
				backupExchangeName, // exchange
				"",                 // routing key
				false,              // mandatory
				false,              // immediate
				amqp.Publishing{
					Body:    d.Body,
					Headers: d.Headers,
				},
			)
			if err != nil {
				failOnError(err, "Failed to publish a message")
				return
			}

			// Wait for the publish confirmation
			isConfirm := confirmOne(confirms)

			// consumer ack
			if !isConfirm {
				log.Printf("Publish not confirmed")
				return
			}

			log.Printf("Sent %s", d.Body)

			err = d.Ack(false)
			if err != nil {
				log.Fatalln("Failed to acknowledge message:", err)
				return
			}
			log.Printf("Consumer ack message successfully")

		}
	}()

	log.Printf("Waiting for logs. To exit press CTRL+C")
	<-forever
}

func confirmOne(confirms <-chan amqp.Confirmation) bool {
	confirmed := <-confirms
	if confirmed.Ack {
		log.Println("Publish confirmed")
		return true
	} else {
		log.Println("Publish not confirmed")
		return false
	}
}
