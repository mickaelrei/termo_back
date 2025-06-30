package module_user

import (
	"github.com/gorilla/mux"
	"net/http"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/modules"
	"termo_back_end/internal/util"
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
		path:    "/user",
	}
}

func (m module) Path() string {
	return m.path
}

func (m module) Setup(r *mux.Router) ([]modules.RouteDefinition, *mux.Router) {
	defs := []modules.RouteDefinition{
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

	response := entities.UserResponse{
		ID:   user.ID,
		Name: user.Name,
	}

	util.WriteResponseJSON(w, response)
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
