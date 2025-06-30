package router

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/modules/module"
	"termo_back_end/internal/modules/repo"
	"termo_back_end/internal/modules/service"
	"time"
)

func Setup(config entities.Config, words []string, db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Repositories
	userRepo := repo.NewUserRepo(db)
	gameRepo := repo.NewGameRepo(db)

	// Services
	userService := service.NewUserService(userRepo)
	gameService := service.NewGameService(words, gameRepo, userRepo)
	authService := service.NewAuthService(config, userRepo)

	// Modules
	userModule := module.NewUserModule(userService, gameService)
	gameModule := module.NewGameModule(gameService)
	authModule := module.NewAuthModule(authService)

	apiModules := []entities.Module{
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
