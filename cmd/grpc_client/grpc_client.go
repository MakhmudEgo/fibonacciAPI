package main

import (
	"encoding/json"
	"fibonacciAPI/internal/pkg/fibonaccigrpc"
	"fibonacciAPI/internal/pkg/parser"
	"fibonacciAPI/pkg/numbers"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/big"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Enter 2 arguments.")
	}
	x, y, err := parser.Parse(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal("Bad arguments. Enter 2 numbers.")
	}
	conn, err := grpc.Dial(os.Getenv("GRPC_SERVICE_PORT"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	cl := fibonaccigrpc.NewFibonacciClient(conn)

	resp, err := cl.Get(context.Background(), &fibonaccigrpc.Request{X: x, Y: y})

	if err != nil {
		log.Fatalln(err)
	}
	res := make([]*big.Int, 0, y-x+1)

	err = json.Unmarshal(resp.Message, &res)
	if err != nil {
		log.Fatalln(err)
	}
	numbers.PrintFibonacci(res)
}

// max 1 6330
