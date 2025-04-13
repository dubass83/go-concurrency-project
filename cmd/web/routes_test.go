package main

import (
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
)

var routes = []string{
	"/login",
	"/logout",
	"/register",
	"/activate",
	"/test-email",
	"/members/plans",
	"/members/subscribe",
}

func TestRoutesExist(t *testing.T) {
	found := make(map[string]bool)

	// Walk all routes
	err := chi.Walk(testApp.Router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		found[route] = true
		return nil
	})

	if err != nil {
		t.Error(err)
	}

	// Check if each required route exists
	for _, route := range routes {
		if !found[route] {
			t.Errorf("Route %s does not exist", route)
		}
	}
}
