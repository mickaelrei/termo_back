package util

import (
	"fmt"
	"log"
	"net/http"
	"termo_back_end/internal/entities"
)

// DefaultEndpointResponse defines a default response format for action endpoints with a template status type
type DefaultEndpointResponse[T any] struct {
	Status        T      `json:"status"`
	StatusMessage string `json:"status_message"`
}

// BuildDefaultEndpointStatusResponse builds a DefaultEndpointResponse with the provided status
func BuildDefaultEndpointStatusResponse[T fmt.Stringer](status T) DefaultEndpointResponse[T] {
	return DefaultEndpointResponse[T]{
		Status:        status,
		StatusMessage: status.String(),
	}
}

// GetUser attempts to retrieve the user in the request's context
func GetUser(r *http.Request) (*entities.User, error) {
	contextUser := r.Context().Value("user")
	if contextUser == nil {
		log.Printf("user not found in request")
		return nil, fmt.Errorf("user not found in request")
	}

	return contextUser.(*entities.User), nil
}
