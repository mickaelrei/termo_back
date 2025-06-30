package module_auth

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/modules"
	"termo_back_end/internal/status_codes"
	"termo_back_end/internal/util"
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
		c := r.Context()
		authHeader := r.Header.Get("Authorization")

		var token string
		if authHeader != "" {
			token = strings.ReplaceAll(authHeader, "Bearer ", "")
		}

		if token == "" {
			log.Printf("No token found in the request\n")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := m.service.GetUserFromToken(c, token)
		if err != nil {
			log.Printf("[GetUserFromToken] | %v", err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if user == nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(c, "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *module) register(w http.ResponseWriter, r *http.Request) {
	var credentials entities.UserCredentials
	if !util.ReadBody(w, r, &credentials) {
		return
	}

	status, token, err := m.service.RegisterUser(r.Context(), credentials)
	if err != nil {
		util.WriteInternalError(w)
		return
	}

	response := struct {
		util.DefaultEndpointResponse[status_codes.UserRegister]
		Token string `json:"token,omitempty"`
	}{
		DefaultEndpointResponse: util.BuildDefaultEndpointStatusResponse(status),
		Token:                   token,
	}

	util.WriteResponseJSON(w, response)
}

func (m *module) login(w http.ResponseWriter, r *http.Request) {
	var credentials entities.UserCredentials
	if !util.ReadBody(w, r, &credentials) {
		return
	}

	status, token, err := m.service.LoginUser(r.Context(), credentials)
	if err != nil {
		util.WriteInternalError(w)
		return
	}

	response := struct {
		util.DefaultEndpointResponse[status_codes.UserLogin]
		Token string `json:"token,omitempty"`
	}{
		DefaultEndpointResponse: util.BuildDefaultEndpointStatusResponse(status),
		Token:                   token,
	}

	util.WriteResponseJSON(w, response)
}
