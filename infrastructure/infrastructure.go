package infrastructure

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func SetupModules(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serverTime := time.Now()
		_, err := fmt.Fprintf(w, "%v - server time", serverTime.UTC().Unix())
		if err != nil {
			log.Println(err)
		}
	})
}
