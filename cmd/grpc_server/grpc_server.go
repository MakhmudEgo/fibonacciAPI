package main

import (
	"encoding/json"
	"fibonacciAPI/internal/pkg/fibonaccigrpc"
	"fibonacciAPI/pkg/service"
	"fibonacciAPI/pkg/storage"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var repo storage.Storage

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
	repo = storage.New()
}

type FibGenerator struct {
	fibonaccigrpc.UnimplementedFibonacciServer
}

func (g *FibGenerator) Get(ctx context.Context, r *fibonaccigrpc.Request) (*fibonaccigrpc.Response, error) {
	res, err := service.Fibonacci(repo).Execute(r.X, r.Y)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}
	return &fibonaccigrpc.Response{Message: data}, nil
}

func main() {
	s := grpc.NewServer()
	fibonaccigrpc.RegisterFibonacciServer(s, &FibGenerator{})

	listen, err := net.Listen("tcp", os.Getenv("GRPC_SERVICE_PORT"))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("server start")
	}

	err = s.Serve(listen)
	if err != nil {
		log.Fatal(err)
	}
}
