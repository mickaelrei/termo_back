package router

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/modules"
	"termo_back_end/internal/modules/module_auth"
	"termo_back_end/internal/modules/module_game"
	"termo_back_end/internal/modules/module_user"
	"time"
)

func Setup(config entities.Config, words []string, db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// User module
	userRepo := module_user.NewRepo(db)
	userService := module_user.NewService(userRepo)
	userModule := module_user.NewModule(userService)

	// Game module
	gameRepo := module_game.NewRepo(db)
	gameService := module_game.NewService(words, gameRepo)
	gameModule := module_game.NewModule(gameService)

	// Auth module
	authService := module_auth.NewService(config, userRepo)
	authModule := module_auth.NewModule(authService)

	apiModules := []modules.Module{
		gameModule,
		userModule,
	}

	// Set up the main auth module for API
	_, adminApiRouter := authModule.Setup(r)

	// Set up all other modules using the auth module's router
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

	return r
}
