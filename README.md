# NATS-SERVICE
- postgres
- jet-stream

## Launch
<hr>

### NATS server and postgres
```bash
$ docker-compose up
```

### Service
- set environmental variables for db, needed variables described in `scripts/set-env-vars.sh`
- Run `main.go`

## Scripts
<hr>

- `create-stream.sh` creates nats stream named `ORDERS` with enabled jet-stream
- `publish.sh` sends an example message to the `ORDERS` stream

## Migrations
<hr>

Run 
```
goose -dir=*migrations directory* *driver* *connection string* up`
```
Example
```
goose -dir="./migrations" postgres "user=postgres dbname=nats-service password=postgres port=5432 sslmode=disable" up
```