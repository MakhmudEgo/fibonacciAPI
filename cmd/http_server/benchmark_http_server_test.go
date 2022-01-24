package main

import (
	"fibonacciAPI/pkg/controller"
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkStrong(b *testing.B) {
	svr := httptest.NewServer(controller.NewFibonacci())
	defer svr.Close()
	req, err := http.NewRequest("GET", svr.URL+"?from=1&to=93", nil)
	if err != nil {
		b.Error(err)
	}
	client := http.Client{}

	for i := 0; i < b.N; i++ {
		resp, err := client.Do(req)
		if err != nil {
			b.Error(err)
		}
		resp.Body.Close()
	}
}
