package auth

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"termo_back_end/internal/modules"
)

type module struct {
	service Service
	path    string
}

func NewModule(service Service) modules.Module {
	return module{
		service: service,
		path:    "/",
	}
}

func (m module) Path() string {
	return m.path
}

func (m module) Setup(r *mux.Router) ([]modules.RouteDefinition, *mux.Router) {
	defs := []modules.RouteDefinition{
		{
			Path:        "/register",
			Handler:     m.register,
			HttpMethods: []string{http.MethodPost},
		},
		{
			Path:        "/login",
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

func (m *module) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Check valid auth token
		// TODO: Check valid access permission

		c := r.Context()
		next.ServeHTTP(w, r.WithContext(c))
	})
}

func (m *module) register(w http.ResponseWriter, _ *http.Request) {
	log.Printf("register endpoint requested")

	// TODO: call m.service.RegisterUser()

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Register"))
}

func (m *module) login(w http.ResponseWriter, _ *http.Request) {
	log.Printf("login endpoint requested")

	// TODO: call m.service.LoginUser()

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Login"))
}
