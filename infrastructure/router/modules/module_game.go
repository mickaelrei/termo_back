package modules

import (
	"github.com/gorilla/mux"
	"net/http"
	"termo_back_end/infrastructure/router"
)

type moduleGame struct {
	name string
	path string
}

func NewGameModule() router.Module {
	return moduleGame{
		name: "Game",
		path: "/game",
	}
}

func (m moduleGame) Name() string {
	return m.name
}

func (m moduleGame) Path() string {
	return m.path
}

func (m moduleGame) Setup(r *mux.Router) ([]router.RouteDefinition, *mux.Router) {
	defs := []router.RouteDefinition{
		{
			Path:        "/start",
			Description: "Start a game",
			Handler:     m.start,
			HttpMethods: []string{http.MethodPost},
		},
	}

	for _, d := range defs {
		r.HandleFunc(d.Path, d.Handler).Methods(d.HttpMethods...)
	}

	return defs, r
}

func (m moduleGame) start(w http.ResponseWriter, _ *http.Request) {
	// TODO: Implement game start
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Start game"))
}
