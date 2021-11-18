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

### Run Trace Consumer
```shell
go run ./main -- consumer trace
```

### Run Consumer
```shell
go run ./main -- consumer
```

### Trace worflow
```shell
curl --location --request POST '0.0.0.0:3000/api/v1/workflows' \
--header 'Content-Type: application/json' \
--data-raw '{"name":"default workflow","start":"s1","end":"s2","payload":1, "operations":[{"name":"op1","from":"s1","to":"s2"},{"name":"op2","from":"s1","to":"s3"},{"name":"op3","from":"s3","to":"s2"}]}'
```

```shell
curl --location --request POST '0.0.0.0:3000/api/v1/workflows' \
--header 'Content-Type: application/json' \
--data-raw '{"name":"portfolio","start":"s1","end":"s3","payload":1, "operations":[{"name":"op1","from":"s1","to":"s2"},{"name":"op2","from":"s2","to":"s3"}]}'
```

### Handle create portfolio workflow
```shell
curl --location --request POST '0.0.0.0:3000/api/v1/create-portfolio' \
--header 'Content-Type: application/json' \
--data-raw '{"uuid":"test-uuid","company_uuid":"test-company-uuid","name":"test-name","description":"test-description","logo_url":"test-logo-url", "portfolio_building_ids":["uuid-1", "uuid-2"]}'
```

## Reset Kafka State
 ```shell
 docker-compose down -v
 ```

