package controller

import (
	"encoding/json"
	"fibonacciAPI/pkg/model"
	"fibonacciAPI/pkg/service"
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
	parseQueryString(w http.ResponseWriter, r *http.Request) (*model.Request, bool)
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
	// parse query string
	d, ok := f.parseQueryString(w, r)
	if !ok {
		return
	}

	// get sequence Fibonacci
	srv := service.NewFibonacci(f.rdb)
	seqFib, err := srv.Execute(d.From, d.To)
	if f.controllerErrors(err, w) {
		return
	}

	// serialization for response
	res, err := json.Marshal(seqFib)
	if f.controllerErrors(err, w) {
		return
	}

	// send response
	_, err = w.Write(res)
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

	f.rdb.LPush(f.rdb.Context(), service.REDIS_FIB_LOG, "controller: "+err.Error()+time.Now().String())

	f.writeError(w, []byte("Status Internal Server Error"), http.StatusInternalServerError)
	return true
}

func (f *Fibonacci) writeError(w http.ResponseWriter, rsp []byte, code int) {
	w.WriteHeader(code)
	if _, err := w.Write(rsp); err != nil {
		f.rdb.LPush(f.rdb.Context(), service.REDIS_FIB_LOG, "controller: "+err.Error()+time.Now().String())
	}
}

func (f *Fibonacci) parseQueryString(w http.ResponseWriter, r *http.Request) (*model.Request, bool) {
	d := model.NewRequest()
	var err error
	d.From, err = strconv.Atoi(r.URL.Query().Get("from"))
	if f.controllerErrors(err, w) {
		return nil, false
	}
	d.To, err = strconv.Atoi(r.URL.Query().Get("to"))
	if f.controllerErrors(err, w) {
		return nil, false
	}
	return d, true
}

func isClientError(err error) bool {
	log.Println(err)
	return strings.Contains(err.Error(), "Atoi: parsing") ||
		strings.Contains(err.Error(), "not a natural") ||
		strings.Contains(err.Error(), "is less than")
}
