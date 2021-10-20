# Saga Coordinator POC
Saga coordinator demonstrates distributed transaction processing

## Getting Started
### Clone the Coordinator
Clone the coordinator locally.
```shell
git clone git@github.com:awe76/saga-coordinator.git
```

### Instal Golang
Kafka library demands go version 1.16 or greater
```shell
curl -o golang.pkg https://dl.google.com/go/go1.16.4.darwin-amd64.pkg
```
Then launch installer

### Launch Docker
```shell
docker-compose up -d
```

## Using the Coordinator
### Run Gateway
```shell
go run ./main
```

### Run Consumer
```shell
go run ./main -- consumer
```

### Send a test message
```shell
curl --location --request POST '0.0.0.0:3000/api/v1/comments' \
--header 'Content-Type: application/json' \
--data-raw '{ "text":"nice boy" }'
```

