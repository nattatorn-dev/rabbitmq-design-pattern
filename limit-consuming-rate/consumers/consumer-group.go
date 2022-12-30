package main

import (
	"log"
	"os"

	"context"
	"strconv"
	"time"

	"github.com/streadway/amqp"

	libredis "github.com/go-redis/redis/v8"
	limiter "github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	ctx := context.Background()
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  10,
	}

	// Create a redis client.
	option, err := libredis.ParseURL("redis://localhost:6379")
	if err != nil {
		failOnError(err, "err")
	}
	client := libredis.NewClient(option)

	// Create a store with the redis client.
	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   "limiter_example",
		MaxRetry: 3,
	})

	if err != nil {
		failOnError(err, "err")
	}

	rateLimiter := limiter.New(store, rate)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"limit-consuming-rate-message", // name
		"direct",                       // type
		true,                           // durable
		false,                          // auto-deleted
		false,                          // internal
		false,                          // no-wait
		nil,                            // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
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
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,                         // queue name
		"",                             // routing key
		"limit-consuming-rate-message", // exchange
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
			endpoint := "api.example.com"
			limiterCtx, err := rateLimiter.Get(ctx, endpoint)
			if err != nil {
				failOnError(err, "err")
			}

			log.Println("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
			log.Println("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
			log.Println("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

			// max limit
			if limiterCtx.Reached {
				log.Printf("Too Many Requests from %s", endpoint)
				// handle max limit case
				return
			}

			// Call External API
			time.Sleep(3 * time.Second)
			// handle max limit case from external API
			// handle error case

			// ACK
			log.Printf("message: %s", d.Body)
			if err := d.Ack(false); err != nil {
				log.Fatalln("Failed to acknowledge message:", err)
				return
			}
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
