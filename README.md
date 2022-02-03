# NATS-SERVICE
- postgres
- jet-stream

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
```

- run `main.go`

## Scripts

- `create-stream.sh` creates nats stream named `ORDERS` with enabled jet-stream
- `publish.sh` sends an example message to the `ORDERS` stream

## Migrations

Run 
```
goose -dir=*migrations directory* *driver* *connection string* up`
```
Example
```
goose -dir="./migrations" postgres "user=postgres dbname=nats-service password=postgres port=5432 sslmode=disable" up
```