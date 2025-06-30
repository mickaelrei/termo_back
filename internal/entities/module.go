package entities

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Module is an interface that represents a set of routes within the application with related logic
//
// It is designed to have a base route path and set all subroutes's handlers
type Module interface {
	// Path returns the module base path
	Path() string

	// Setup sets all the route handlers to the subroutes and returns all the route definitions
	//
	// Optionally also returns a new router base. This is useful if a module has important logic and their subroutes
	// should have some rules, like the AuthModule which defines a session middleware for all /api subroutes
	Setup(r *mux.Router) ([]RouteDefinition, *mux.Router)
}

type RouteDefinition struct {
	// Path is the path for the route
	Path string

	// Handler is the function handler for the route
	Handler http.HandlerFunc

	// HttpMethods is a list of HTTP methods accepted by the route
	HttpMethods []string
}
