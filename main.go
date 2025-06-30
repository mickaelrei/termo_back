package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"termo_back_end/internal/entities"
	"termo_back_end/internal/router"
	"termo_back_end/internal/util"
	"time"
)

const ServerPort = 8080

func loadWordListFile() ([]string, error) {
	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatalf("failed to load word list: %v", err)
	}
	defer util.DeferFileClose(file)

	var words []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed reading word list: %v", err)
	}

	log.Printf("loaded %d words", len(words))

	return words, nil
}

func openDB(config entities.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
	)

	// Open database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("[sql.Open] | %v", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("[db.Ping] | %v", err)
	}

	return db, nil
}

func readConfig() (*entities.Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, fmt.Errorf("[os.Open] | %v", err)
	}
	defer util.DeferFileClose(file)

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("[io.ReadAll] | %v", err)
	}

	var config entities.Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("[json.Unmarshal] | %v", err)
	}

	return &config, nil
}

func createServer(router *mux.Router) *http.Server {
	corsOptions := handlers.AllowedOriginValidator(func(s string) bool {
		// TODO: Only allow front-end specific origins
		log.Printf("origin: %s", s)
		return true
	})

	server := http.Server{
		Addr: fmt.Sprintf(":%d", ServerPort),
		Handler: handlers.CORS(
			corsOptions,
			handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Accept"}),
			handlers.AllowedMethods([]string{"GET", "PUT", "POST", "DELETE"}),
			handlers.AllowCredentials(),
		)(router),
		ReadTimeout:       time.Second * 120,
		WriteTimeout:      time.Second * 120,
		ReadHeaderTimeout: time.Second * 2,
		IdleTimeout:       time.Second * 60,
	}

	return &server
}

func main() {
	// Load the words-list file
	words, err := loadWordListFile()
	if err != nil {
		log.Fatalf("[loadWordListFile] | %v", err)
		return
	}

	// Ensure the list is not empty
	if len(words) == 0 {
		log.Fatalf("word list is empty")
		return
	}

	// Load config file
	config, err := readConfig()
	if err != nil {
		log.Fatalf("[readConfig] | %v", err)
		return
	}

	// Open database
	db, err := openDB(*config)
	if err != nil {
		log.Fatalf("[openDB] | %v", err)
		return
	}

	// Set up all route handlers
	r := router.Setup(*config, words, db)

	// Create server
	server := createServer(r)

	// Start listening to requests
	log.Printf("Starting server on address: %s", server.Addr)
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("[server.ListenAndServe] | %v", err)
	}
}
