package main

import (
	"context"
	"log"
	"strconv"
	"time"

	libredis "github.com/go-redis/redis/v8"
	limiter "github.com/ulule/limiter/v3"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func rate() {
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
	limiterCtx, err := rateLimiter.Get(ctx, "api.example.com")
	if err != nil {
		failOnError(err, "err")
	}

	log.Println("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
	log.Println("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
	log.Println("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

	if limiterCtx.Reached {
		log.Printf("Too Many Requests from %s", "api.example.com")
	}
}
