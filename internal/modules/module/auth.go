package module

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/modules/service"
	"termo_back_end/internal/status_codes"
	"termo_back_end/internal/util"
)

type authModule struct {
	service service.AuthService
	path    string
}

func NewAuthModule(service service.AuthService) entities.Module {
	return authModule{
		service: service,
		path:    "/",
	}
}

func (m authModule) Path() string {
	return m.path
}

func (m authModule) Setup(r *mux.Router) ([]entities.RouteDefinition, *mux.Router) {
	defs := []entities.RouteDefinition{
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

func (m *authModule) sessionMiddleware(next http.Handler) http.Handler {
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

func (m *authModule) register(w http.ResponseWriter, r *http.Request) {
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

func (m *authModule) login(w http.ResponseWriter, r *http.Request) {
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
