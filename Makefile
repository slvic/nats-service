# Database
.PHONY: "migrate"
migrate:
	goose -dir="./migrations" postgres "user=postgres dbname=nats-service password=postgres port=5432 sslmode=disable" up

# NATS
.PHONY: "create-example-nats-stream"
create-example-nats-stream:
	nats stream add ORDERS --subjects "ORDERS.*" --ack --max-msgs=-1 --max-bytes=-1 --max-age=1y --storage file --retention limits --max-msg-size=-1 --discard=old

# Docker image
.PHONY: "build"
build:
	CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd/nats-service/

.PHONY: "create-docker-image"
create-docker-image:
	docker build -t nats-service-scratch -f Dockerfile.scratch .

# Docker network
.PHONY: "create-docker-network"
create-docker-network:
	docker network create nats_service

.PHONY: "docker-connect-nats"
docker-connect-nats:
	docker network connect nats_service nats_service_nats_server

.PHONY: "docker-connect-postgres"
docker-connect-postgres:
	docker network connect nats_service nats_service_postgres

.PHONY: "docker-connect"
docker-connect: create-docker-network docker-connect-nats docker-connect-postgres

# Docker example
.PHONY: "run-first-example-container"
run-first-example-container:
	docker run --network nats_service -d -p 3000:3000 \
	--env NATS_SERVICE_DB_USER='postgres' \
	--env NATS_SERVICE_DB_NAME='nats-service' \
	--env NATS_SERVICE_DB_PASSWORD='postgres' \
	--env NATS_SERVICE_DB_HOST="nats_service_postgres" \
	--env NATS_SERVICE_DB_PORT='5432' \
	--env NATS_SERVICE_DB_SSL_MODE='disable' \
	--env NATS_SERVICE_NATS_HOST="nats_service_nats_server" \
	--env NATS_SERVICE_NATS_PORT="4222" \
	 slvic/nats-service-scratch:latest

.PHONY: "run-second-example-container"
run-second-example-container:
	docker run --network nats_service -d -p 3001:3000 \
    	--env NATS_SERVICE_DB_USER='postgres' \
    	--env NATS_SERVICE_DB_NAME='nats-service' \
    	--env NATS_SERVICE_DB_PASSWORD='postgres' \
   		--env NATS_SERVICE_DB_HOST="nats_service_postgres" \
    	--env NATS_SERVICE_DB_PORT='5432' \
    	--env NATS_SERVICE_DB_SSL_MODE='disable' \
		--env NATS_SERVICE_NATS_HOST="nats_service_nats_server" \
		--env NATS_SERVICE_NATS_PORT="4222" \
    	slvic/nats-service-scratch:latest

.PHONY: "run-example"
run-example: run-first-example-container run-second-example-container