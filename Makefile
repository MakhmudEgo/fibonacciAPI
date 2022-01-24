run: build grpc_client
	docker-compose up -d
build:
	docker-compose build



grpc_client:
	go build cmd/grpc_client/grpc_client.go

stop:
	docker-compose stop

delete:
	docker-compose down

.PHONY: all build run grpc_client stop delete
