# NATS-SERVICE
- postgres
- jet-stream
- docker

### Example launch
```bash
$ docker-compose up # runs NATS server and postgres
$ make docker-connect # connects NATS server and postgres containers to docker network
$ make migrate
$ make create-example-nats-stream
$ make run-example
```
To publish a sample message to nats server run `/scripts/publish.sh` script.

To check if your message has been delivered open `/web/page/index.html` and hit "get message button" 


## Launch

### NATS server and postgres
```bash
$ docker-compose up
```

### Service
- set the following environmental variables:
```bash
POSTGRES_USER
POSTGRES_NAME
POSTGRES_PASSWORD
POSTGRES_PORT
POSTGRES_SSL_MODE
NATS_SERVICE_NATS_HOST
NATS_SERVICE_NATS_PORT
```

- run `main.go`

## Scripts

- `create-stream.sh` creates nats stream named `ORDERS` with enabled jet-stream
- `publish.sh` sends a sample message to the `ORDERS` stream

## Migrations

Run 
```
goose -dir=*migrations directory* *driver* *connection string* up`
```
Example
```
goose -dir="./migrations" postgres "user=postgres dbname=nats-service password=postgres port=5432 sslmode=disable" up
```