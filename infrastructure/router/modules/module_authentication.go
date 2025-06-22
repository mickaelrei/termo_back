package modules

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"termo_back_end/infrastructure/router"
)

type moduleAuthentication struct {
	name string
	path string
}

func NewAuthenticationModule() router.Module {
	return moduleAuthentication{
		name: "Authentication",
		path: "/",
	}
}

func (m moduleAuthentication) Name() string {
	return m.name
}

func (m moduleAuthentication) Path() string {
	return m.path
}

func (m moduleAuthentication) Setup(r *mux.Router) ([]router.RouteDefinition, *mux.Router) {
	defs := []router.RouteDefinition{
		{
			Path:        "/login",
			Description: "User login",
			Handler:     m.login,
			HttpMethods: []string{http.MethodPost},
		},
	}

	for _, d := range defs {
		r.HandleFunc(d.Path, d.Handler).Methods(d.HttpMethods...)
	}

	// Add /api prefix to all other modules
	api := r.PathPrefix("/api").Subrouter()
	api.Use(m.sessionMiddleware)

	return defs, api
}

func (m *moduleAuthentication) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Check valid auth token
		// TODO: Check valid access permission

		c := r.Context()
		next.ServeHTTP(w, r.WithContext(c))
	})
}

func (m *moduleAuthentication) login(w http.ResponseWriter, _ *http.Request) {
	// TODO: Implement login
	log.Printf("login endpoint requested")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Login"))
}
