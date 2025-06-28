package user

import (
	"github.com/gorilla/mux"
	"termo_back_end/internal/modules"
)

/*

Endpoints:
 - /getData
 - /changeName
 - /changePassword

*/

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
	// TODO: Implement
	return nil, nil
}
