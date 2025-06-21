package infrastructure

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"termo_back_end/infrastructure/router"
	"termo_back_end/infrastructure/router/modules"
	"time"
)

func SetupModules(r *mux.Router) {
	moduleAuthentication := modules.NewAuthenticationModule()
	moduleGame := modules.NewGameModule()

	apiModules := []router.Module{
		moduleGame,
	}

	// Set up the main auth module for API
	_, adminApiRouter := moduleAuthentication.Setup(r)
	for _, m := range apiModules {
		moduleSubRouter := adminApiRouter.PathPrefix(m.Path()).Subrouter()
		_, _ = m.Setup(moduleSubRouter)
	}

	// Home URL handler returns the current server time
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serverTime := time.Now()
		_, err := fmt.Fprintf(w, "%d", serverTime.UTC().Unix())
		if err != nil {
			log.Println(err)
		}
	})
}
