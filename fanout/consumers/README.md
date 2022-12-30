1 Run Consumer

```sh
$ go run fanout/consumers/consumer-group.go <queue_name_1>
$ go run fanout/consumers/consumer-group.go <queue_name_2>
$ go run fanout/consumers/consumer-group.go <queue_name_3>
```

2 Publish message via RabbitMQ Management UI