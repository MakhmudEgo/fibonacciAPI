package main

import (
	"fibonacciAPI/pkg/controller"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
}

func main() {
	http.Handle("/fibonacci", controller.NewFibonacci())
	log.Fatalln(http.ListenAndServe(os.Getenv("SERVER_PORT"), nil))
}
