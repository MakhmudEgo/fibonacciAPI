package controller

import (
	"encoding/json"
	"fibonacciAPI/pkg/service"
	"fibonacciAPI/pkg/storage"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type iFibonacci interface {
	controllerErrors(error, http.ResponseWriter) bool
	writeError(w http.ResponseWriter, rsp []byte, code int)
	parseQueryString(w http.ResponseWriter, r *http.Request) (int, int, bool)
}

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
		f.writeError(w, []byte("I’m captain Jack Sparrow "), http.StatusTeapot)
		return
	}
	// parse query string
	from, to, ok := f.parseQueryString(w, r)
	if !ok {
		return
	}

	// get sequence Fibonacci
	srv := service.NewFibonacci(f.rdb)
	seqFib, err := srv.Execute(from, to)
	if f.controllerErrors(err, w) {
		return
	}

	str := fmt.Sprint(seqFib)
	err = json.NewEncoder(w).Encode(strings.Split(str[1:len(str)-1], " "))
	if f.controllerErrors(err, w) {
		return
	}

}

func (f *Fibonacci) controllerErrors(err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}
	if isClientError(err) {
		if strings.Contains(err.Error(), "Atoi: parsing") {
			f.writeError(w, []byte("bad query args"), http.StatusBadRequest)
		} else {
			f.writeError(w, []byte(err.Error()), http.StatusBadRequest)
		}
		return true
	}

	f.rdb.LPush(f.rdb.Context(), storage.REDIS_FIB_LOG, "controller: "+err.Error()+time.Now().String())

	if strings.Contains(err.Error(), "overflow") {
		f.writeError(w, []byte("The service supports up to 93 Fibonacci numbers,\nbut we are already working on increasing the numbers)"), http.StatusInternalServerError)
	} else {
		f.writeError(w, []byte("Status Internal Server Error"), http.StatusInternalServerError)
	}
	return true
}

func (f *Fibonacci) writeError(w http.ResponseWriter, rsp []byte, code int) {
	w.WriteHeader(code)
	if _, err := w.Write(rsp); err != nil {
		f.rdb.LPush(f.rdb.Context(), storage.REDIS_FIB_LOG, "controller: "+err.Error()+time.Now().String())
	}
}

func (f *Fibonacci) parseQueryString(w http.ResponseWriter, r *http.Request) (int, int, bool) {
	var from, to int
	var err error
	from, err = strconv.Atoi(r.URL.Query().Get("from"))
	if f.controllerErrors(err, w) {
		return from, to, false
	}
	to, err = strconv.Atoi(r.URL.Query().Get("to"))
	if f.controllerErrors(err, w) {
		return from, to, false
	}
	return from, to, true
}

func isClientError(err error) bool {
	log.Println(err)
	return strings.Contains(err.Error(), "Atoi: parsing") ||
		strings.Contains(err.Error(), "not a natural") ||
		strings.Contains(err.Error(), "is less than")
}
