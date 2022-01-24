// Package routes – lists http endpoints
package routes

import (
	"fibonacciAPI/pkg/controller"
	"net/http"
)

func init() {
	http.Handle("/fibonacci", controller.NewFibonacci())
}
