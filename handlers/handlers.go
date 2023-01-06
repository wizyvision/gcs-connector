package handlers

import (
	"cloud-storage-connector/run"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func SetupEndpoints(router *mux.Router) {
	fmt.Println("SetupEndpoints")
	router.HandleFunc("/", verifyToken(homeHandler))
	router.HandleFunc("/run", verifyToken(runHandler)).Methods("GET")
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&map[string]string{
		"service": "Upload to wizy visions's upload file public api",
		"usage":   "/run?gcsObject=gcsObjectValue&gcsBucket=gcsBucketValue",
	})
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	gcsObject := r.URL.Query().Get("gcsObject")
	gcsBucket := r.URL.Query().Get("gcsBucket")
	if gcsObject == "" || gcsBucket == "" {
		http.Error(w, "gcsObject and gcsBucet are required", http.StatusBadRequest)
		return
	}

	status, err := run.Execute(gcsObject, gcsBucket)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error on upload : %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&map[string]string{
		"status":     status,
		"gcs object": gcsObject,
		"gcs bucket": gcsBucket,
	})
}

func verifyToken(endpointHandler func(writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tokens := request.Header[os.Getenv("UPLOADER_SERVICE_AUTH_HEADER")]
		if tokens != nil {
			if tokens[0] != os.Getenv("UPLOADER_SERVICE_AUTH_TOKEN") {
				writer.WriteHeader(http.StatusUnauthorized)
				_, err := writer.Write([]byte("You're Unauthorized"))
				if err != nil {
					return
				}
			} else {
				endpointHandler(writer, request)
			}
		} else {
			writer.WriteHeader(http.StatusUnauthorized)
			_, err := writer.Write([]byte("You're Unauthorized due to No token in the header"))
			if err != nil {
				return
			}
		}
	})
}
