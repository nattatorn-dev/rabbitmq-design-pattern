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

	err = ch.Qos(
		10,   // prefetch count
		0,    // prefetch size
		true, // global
	)
	failOnError(err, "Failed to set QoS")

	queueName := string(os.Args[1])
	log.Println(">>> Queue Name", queueName)

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{ // arguments
			"x-max-priority": 9, // set the maximum priority to 9
		},
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,             // queue name
		"",                 // routing key
		"priority-message", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

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

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("priority %d", d.Priority)
			log.Printf("%s", d.Body)
			if err := d.Ack(false); err != nil {
				log.Fatalln("Failed to acknowledge message:", err)
			}
		}
	}()

	log.Printf("Waiting for logs. To exit press CTRL+C")
	<-forever
}
