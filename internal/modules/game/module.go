package game

import (
	"github.com/gorilla/mux"
	"termo_back_end/internal/modules"
)

/*

Endpoints:
 - /startGame: returns a game id/hash
 - /attempt: receives a game id/hash and a word attempt; returns what's wrong/right

Notes:
 - When the user logs in, we return their unfinished game (if any)
 - Maybe add a lifetime to games? so that games of users inactive for too long are deleted

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
