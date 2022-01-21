package controller

import (
	"fibonacciAPI/pkg/service"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Fibonacci struct {
	rdb *redis.Client
}

func NewFibonacci() *Fibonacci {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	f := &Fibonacci{rdb: rdb}
	pong, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println(pong)
	}
	return f
}

func (f *Fibonacci) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "I’m captain Jack Sparrow", http.StatusTeapot)
		return
	}

	// parse query string
	from, to, err := parseQueryString(r)
	if Error(w, err) {
		return
	}

	// get sequence Fibonacci
	srv := service.Fibonacci(f.rdb)
	seqFib, err := srv.Execute(from, to)
	if Error(w, err) {
		return
	}

	//w.Header().Add("content-type", "json/application")
	_, err = fmt.Fprint(w, seqFib)
	Error(w, err)
}

func Error(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	if isClientError(err) {
		if strings.Contains(err.Error(), "Atoi: parsing") {
			http.Error(w, "bad query args", http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return true
	}

	if strings.Contains(err.Error(), "overflow") {
		http.Error(w, "The service supports up to 93 Fibonacci numbers,\nbut we are already working on increasing the numbers)", http.StatusInternalServerError)
	} else {
		http.Error(w, "Status Internal Server Error", http.StatusInternalServerError)
	}
	return true
}

func parseQueryString(r *http.Request) (int, int, error) {
	var from, to int
	var err error
	from, err = strconv.Atoi(r.URL.Query().Get("from"))
	if err != nil {
		return from, to, err
	}
	to, err = strconv.Atoi(r.URL.Query().Get("to"))

	return from, to, err
}

func isClientError(err error) bool {
	log.Println(err)
	return strings.Contains(err.Error(), "Atoi: parsing") ||
		strings.Contains(err.Error(), "not a natural") ||
		strings.Contains(err.Error(), "is less than")
}
