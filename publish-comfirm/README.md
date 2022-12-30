1 Publish message
```sh
$ go run delayed-queue/producers/producer-group.go <message> <delay_ms>
```

2 Run Consumer

```sh
$ go run delayed-queue/consumers/consumer-group.go <queue_name>
```


Download Tools
https://github.com/rabbitmq/rabbitmq-perf-test/releases

Run Command
```sh
$ rabbitmq-perf-test-2.18.0/bin/runjava com.rabbitmq.perf.PerfTest --help
```

Example
```sh
$ rabbitmq-perf-test-2.18.0/bin/runjava com.rabbitmq.perf.PerfTest \
-h amqp://guest:guest@localhost:5672 \
--time 20 \
--queue "p1" \
--queue-args x-max-priority=9 \
--auto-delete false \
--producers 2 \
--consumers 0 \
--size 100 \
--message-properties priority=1
```


```sh
$ rabbitmq-perf-test-2.18.0/bin/runjava com.rabbitmq.perf.PerfTest \
-h amqp://guest:guest@localhost:5672 \
--time 20 \
--queue "p1" \
--queue-args x-max-priority=9 \
--auto-delete false \
--producers 2 \
--consumers 0 \
--size 100 \
--message-properties priority=5
```