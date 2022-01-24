package controller

import (
	"encoding/json"
	"fibonacciAPI/internal/pkg/parser"
	"fibonacciAPI/pkg/service"
	"fibonacciAPI/pkg/storage"
	"log"
	"net/http"
	"strings"
)

type fibonacciController struct {
	repo storage.Storage
}

func NewFibonacci() *fibonacciController {
	return &fibonacciController{repo: storage.New()}
}

func (f *fibonacciController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "I’m captain Jack Sparrow", http.StatusTeapot)
		return
	}

	// parse query string
	from, to, err := parser.Parse(r.URL.Query().Get("from"), r.URL.Query().Get("to"))
	if Error(w, err) {
		return
	}

	// get sequence Fibonacci
	srv := service.Fibonacci(f.repo)
	seqFib, err := srv.Execute(from, to)
	if Error(w, err) {
		return
	}

	w.Header().Add("content-type", "application/json")
	err = json.NewEncoder(w).Encode(seqFib)
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
	http.Error(w, "Status Internal Server Error", http.StatusInternalServerError)
	return true
}

func isClientError(err error) bool {
	log.Println(err)
	return strings.Contains(err.Error(), "Atoi: parsing") ||
		strings.Contains(err.Error(), "not a natural") ||
		strings.Contains(err.Error(), "is less than")
}
