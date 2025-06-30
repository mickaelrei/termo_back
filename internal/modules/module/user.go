package module

import (
	"github.com/gorilla/mux"
	"net/http"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/modules/service"
	"termo_back_end/internal/rules"
	"termo_back_end/internal/util"
)

type module struct {
	service     service.UserService
	gameService service.GameService
	path        string
}

func NewUserModule(service service.UserService, gameService service.GameService) entities.Module {
	return module{
		service:     service,
		gameService: gameService,
		path:        "/user",
	}
}

func (m module) Path() string {
	return m.path
}

func (m module) Setup(r *mux.Router) ([]entities.RouteDefinition, *mux.Router) {
	defs := []entities.RouteDefinition{
		{
			Path:        "/getData",
			Handler:     m.getData,
			HttpMethods: []string{http.MethodGet},
		},
		{
			Path:        "/updateName",
			Handler:     m.updateName,
			HttpMethods: []string{http.MethodPost},
		},
		{
			Path:        "/updatePassword",
			Handler:     m.updatePassword,
			HttpMethods: []string{http.MethodPost},
		},
	}

	for _, d := range defs {
		r.HandleFunc(d.Path, d.Handler).Methods(d.HttpMethods...)
	}

	return defs, nil
}

func (m module) getData(w http.ResponseWriter, r *http.Request) {
	user, err := util.GetUser(r)
	if err != nil {
		util.WriteInternalError(w)
		return
	}

	// Get active game
	game, gameStatuses, err := m.gameService.GetUserActiveGame(r.Context(), user)
	if err != nil {
		util.WriteInternalError(w)
		return
	}

	util.WriteResponseJSON(w, user.ToResponse(
		game,
		gameStatuses,
		rules.GetGameMaxAttempts(game.GetWordLength(), game.GetCount()),
	))
}

func (m module) updateName(w http.ResponseWriter, r *http.Request) {
	user, err := util.GetUser(r)
	if err != nil {
		util.WriteInternalError(w)
		return
	}

	var body struct {
		NewName string `json:"new_name"`
	}
	if !util.ReadBody(w, r, &body) {
		return
	}

	status, err := m.service.UpdateName(r.Context(), user, body.NewName)
	if err != nil {
		util.WriteInternalError(w)
		return
	}

	util.WriteResponseJSON(w, util.BuildDefaultEndpointStatusResponse(status))
}

func (m module) updatePassword(w http.ResponseWriter, r *http.Request) {
	user, err := util.GetUser(r)
	if err != nil {
		util.WriteInternalError(w)
		return
	}

	var body struct {
		NewPassword string `json:"new_password"`
	}
	if !util.ReadBody(w, r, &body) {
		return
	}

	status, err := m.service.UpdatePassword(r.Context(), user, body.NewPassword)
	if err != nil {
		util.WriteInternalError(w)
		return
	}

	util.WriteResponseJSON(w, util.BuildDefaultEndpointStatusResponse(status))
}
