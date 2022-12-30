1 Publish message
```sh
$ go run limit-consuming-rate/producers/producer-group.go <message>
```

2 Run Consumer

```sh
$ go run limit-consuming-rate/consumers/consumer-group.go <queue_name_1>
$ go run limit-consuming-rate/consumers/consumer-group.go <queue_name_2>
$ go run limit-consuming-rate/consumers/consumer-group.go <queue_name_3>
```
