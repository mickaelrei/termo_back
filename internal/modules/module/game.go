package module

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/modules/service"
	"termo_back_end/internal/rules"
	"termo_back_end/internal/status_codes"
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

type gameModule struct {
	service service.GameService
	path    string
}

func NewGameModule(service service.GameService) entities.Module {
	return gameModule{
		service: service,
		path:    "/game",
	}
}

func (m gameModule) Path() string {
	return m.path
}

func (m gameModule) Setup(r *mux.Router) ([]entities.RouteDefinition, *mux.Router) {
	defs := []entities.RouteDefinition{
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
		{
			Path:        "/getActive",
			Handler:     m.getActive,
			HttpMethods: []string{http.MethodGet},
		},
	}

	for _, d := range defs {
		r.HandleFunc(d.Path, d.Handler).Methods(d.HttpMethods...)
	}

	return defs, nil
}

func (m gameModule) start(w http.ResponseWriter, r *http.Request) {
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

	var maxTries *uint32
	if status == status_codes.GameStartSuccess {
		tries := rules.GetGameMaxAttempts(body.WordLength, body.GameCount)
		maxTries = &tries
	}
	response := struct {
		util.DefaultEndpointResponse[status_codes.GameStart]
		MaxTries *uint32 `json:"max_tries,omitempty"`
	}{
		DefaultEndpointResponse: util.BuildDefaultEndpointStatusResponse(status),
		MaxTries:                maxTries,
	}

	util.WriteResponseJSON(w, response)
}

func (m gameModule) attempt(w http.ResponseWriter, r *http.Request) {
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

	status, gameStatus, err := m.service.AttemptGame(r.Context(), user, body.Attempt)
	if err != nil {
		log.Printf("[AttemptGame] | %v", err)
		util.WriteInternalError(w)
		return
	}

	response := struct {
		util.DefaultEndpointResponse[status_codes.GameAttempt]
		GameStatus []entities.GameWordStatus `json:"game_status,omitempty"`
	}{
		DefaultEndpointResponse: util.BuildDefaultEndpointStatusResponse(status),
		GameStatus:              gameStatus,
	}

	util.WriteResponseJSON(w, response)
}

func (m gameModule) getActive(w http.ResponseWriter, r *http.Request) {
	user, err := util.GetUser(r)
	if err != nil {
		util.WriteInternalError(w)
		return
	}

	game, gameStatuses, err := m.service.GetUserActiveGame(r.Context(), user)
	if err != nil {
		log.Printf("[AttemptGame] | %v", err)
		util.WriteInternalError(w)
		return
	}

	util.WriteResponseJSON(
		w,
		game.ToResponse(gameStatuses, rules.GetGameMaxAttempts(game.GetWordLength(), game.GetCount())),
	)
}
