## Start RabbitMQ Server

```
$ docker-compose up -d
```

Admin page at `http://localhost:15672/`

- username=guest
- password=guest

## Start Consumer

```
$go run consumer-group.go <queue_name>
```

## Start Producer

```
$go run producer-group.go <message>
```
