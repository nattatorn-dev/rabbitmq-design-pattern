FYI: At least one - The consumer should handle idempotency

Type: at least one without QOS  
Run
```
$ go run tools/queue-migration/abuse-mode/consumer-group.go <source_queue>
```
20000 message / 1 sec

Type: at least one without publisher confirm/transaction  
Run
```
$ go run tools/queue-migration/no-guarantee/consumer-group.go <source_queue>
```
20000 message / 28 sec

Type: at least one  

Run
```
$ go run tools/queue-migration/at-least-one/consumer-group.go <source_queue>
```
20000 message / 55 sec

Type: exactly-once  
Run
```
$ go run tools/queue-migration/exactly-once/consumer-group.go <source_queue>
```
20000 message / 48 sec

Perf  
1 Create a new queue <queue_name>
2 Run consumer
```sh
$ go run tools/queue-migration/consumer-group.go <queue_name>
```

3 Run Perf
```sh
$ rabbitmq-perf-test-2.18.0/bin/runjava com.rabbitmq.perf.PerfTest \
-h amqp://guest:guest@localhost:5672 \
--time 5 \
--queue <queue_name> \
-f persistent \
--auto-delete false \
--producers 2 \
--consumers 0 \
--size 100
```
