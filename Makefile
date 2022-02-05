.PHONY: "build"
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/nats-service/

.PHONY: "create-docker-image"
create-docker-image:
	docker build -t nats-service-scratch -f Dockerfile.scratch .

.PHONY: "create-docker-network"
create-docker-network:
	docker network create nats_service

.PHONY: "run-first-example-container"
run-first-example-container:
	docker run --network nats_service -d -p 3000:3000 \
	--env NATS_SERVICE_DB_USER='postgres' \
	--env NATS_SERVICE_DB_NAME='nats-service' \
	--env NATS_SERVICE_DB_PASSWORD='postgres' \
	--enc NATS_SERVICE_DB_HOST="nats_service_postgres" \
	--env NATS_SERVICE_DB_PORT='5432' \
	--env NATS_SERVICE_DB_SSL_MODE='disable' \
	nats-service-scratch

.PHONY: "run-second-example-container"
run-second-example-container:
	docker run --network nats_service -d -p 3001:3000 \
    	--env NATS_SERVICE_DB_USER='postgres' \
    	--env NATS_SERVICE_DB_NAME='nats-service' \
    	--env NATS_SERVICE_DB_PASSWORD='postgres' \
   		--enc NATS_SERVICE_DB_HOST="nats_service_postgres" \
    	--env NATS_SERVICE_DB_PORT='5432' \
    	--env NATS_SERVICE_DB_SSL_MODE='disable' \
    	nats-service-scratch