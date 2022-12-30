1 Publish message
```sh
$ go run delayed-queue/producers/producer-group.go <message> <delay_ms>
```

2 Run Consumer

```sh
$ go run delayed-queue/consumers/consumer-group.go <queue_name>
```

