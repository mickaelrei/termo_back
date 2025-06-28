package router

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"termo_back_end/internal/modules"
	"termo_back_end/internal/modules/auth"
	"termo_back_end/internal/modules/game"
	"termo_back_end/internal/modules/user"
	"time"
)

func SetupRouter(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Auth module
	authRepo := auth.NewRepo(db)
	authService := auth.NewService(authRepo)
	authModule := auth.NewModule(authService)

	// User module
	userRepo := user.NewRepo(db)
	userService := user.NewService(userRepo)
	userModule := user.NewModule(userService)

	// Game module
	gameRepo := game.NewRepo(db)
	gameService := game.NewService(gameRepo)
	gameModule := game.NewModule(gameService)

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
