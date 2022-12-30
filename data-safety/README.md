1 Run Publisher

```sh
$ go run data-safety/transaction/producers/producer-group.go <round>
```

2 Run Consumer

```sh
$ go run data-safety/transaction/consumers/consumer-group.go <queue_name>
```

### Transaction
1 Run Publisher

```sh
$ go run data-safety/transaction/producers/transaction/producer-group <round>
```


### Publisher Confirm
better than transactions
1 Run Publisher

```sh
$ go run data-safety/publisher-confirm/producers/confirm/producer-group <round>
```
