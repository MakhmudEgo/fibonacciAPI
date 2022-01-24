FROM golang:latest

WORKDIR /go/src/
COPY ./cmd/grpc_server ./cmd/grpc_server
COPY ./pkg ./pkg
COPY ./internal/pkg ./internal/pkg
COPY ./go.mod ./
COPY ./.env ./
RUN go mod tidy
RUN go build ./cmd/grpc_server
CMD ["./grpc_server"]
