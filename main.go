package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"termo_back_end/infrastructure"
	"time"
)

const ServerPort = 8080

func main() {
	corsOptions := handlers.AllowedOriginValidator(func(s string) bool {
		// TODO: Only allow front-end specific origins
		return true
	})

	r := mux.NewRouter()
	infrastructure.SetupModules(r)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", ServerPort),
		Handler: handlers.CORS(
			corsOptions,
			handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Accept"}),
			handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"}),
			handlers.AllowCredentials(),
		)(r),
		ReadTimeout:       time.Second * 120,
		WriteTimeout:      time.Second * 120,
		ReadHeaderTimeout: time.Second * 2,
		IdleTimeout:       time.Second * 60,
	}

	log.Printf("Starting server on address: %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("Error: %v", err)
	}
}
