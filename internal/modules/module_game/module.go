package module_game

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"termo_back_end/internal/modules"
	"termo_back_end/internal/util"
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
		path:    "/game",
	}
}

func (m module) Path() string {
	return m.path
}

func (m module) Setup(r *mux.Router) ([]modules.RouteDefinition, *mux.Router) {
	defs := []modules.RouteDefinition{
		{
			Path:        "/start",
			Handler:     m.start,
			HttpMethods: []string{http.MethodPost},
		},
		{
			Path:        "/attempt",
			Handler:     m.attempt,
			HttpMethods: []string{http.MethodPost},
		},
	}

	for _, d := range defs {
		r.HandleFunc(d.Path, d.Handler).Methods(d.HttpMethods...)
	}

	return defs, nil
}

func (m module) start(w http.ResponseWriter, r *http.Request) {
	user, err := util.GetUser(r)
	if err != nil {
		util.WriteInternalError(w)
		return
	}

	var body struct {
		WordLength uint32 `json:"word_length"`
		GameCount  uint32 `json:"game_count"`
	}
	if !util.ReadBody(w, r, &body) {
		return
	}

	status, err := m.service.StartGame(r.Context(), user, body.WordLength, body.GameCount)
	if err != nil {
		log.Printf("[StartGame] | %v", err)
		util.WriteInternalError(w)
		return
	}

	util.WriteResponseJSON(w, util.BuildDefaultEndpointStatusResponse(status))
}

func (m module) attempt(w http.ResponseWriter, r *http.Request) {
	user, err := util.GetUser(r)
	if err != nil {
		util.WriteInternalError(w)
		return
	}

	var body struct {
		Attempt string `json:"attempt"`
	}
	if !util.ReadBody(w, r, &body) {
		return
	}

	status, err := m.service.AttemptGame(r.Context(), user, body.Attempt)
	if err != nil {
		log.Printf("[AttemptGame] | %v", err)
		util.WriteInternalError(w)
		return
	}

	util.WriteResponseJSON(w, util.BuildDefaultEndpointStatusResponse(status))
}
