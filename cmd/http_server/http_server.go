package main

import (
	_ "fibonacciAPI/pkg/routes"
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
	log.Println("server start")
	err := http.ListenAndServe(os.Getenv("SERVER_PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
