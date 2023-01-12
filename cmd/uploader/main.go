package main

import (
	"log"
	"net/http"
	"os"

	"github.com/wizyvision/gcs-connector/cmd/uploader/handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	handlers.SetupEndpoints(router)

	handleAppengineInternalEndpoint(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
		log.Printf("Defaulting to port %s", port)
	}
	log.Printf("Listening on port %s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

func handleAppengineInternalEndpoint(router *mux.Router) {
	// Called when the instance is started by App Engine
	router.HandleFunc("/_ah/start", func(w http.ResponseWriter, r *http.Request) {
	})

	// Called when the instance is stopped by App Engine
	router.HandleFunc("/_ah/stop", func(w http.ResponseWriter, r *http.Request) {
	})
}
