FROM golang:latest

WORKDIR /go/src/
COPY ./cmd/http_server ./cmd/http_server
COPY ./pkg ./pkg
COPY ./internal/pkg/parser ./internal/pkg/parser
COPY ./go.mod ./
COPY ./.env ./
RUN go mod tidy
RUN go build ./cmd/http_server
CMD ["./http_server"]
