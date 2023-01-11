package logger

import (
	"context"
	"log"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
)

type ErrorLogPayload struct {
	Message string
	Error   string
}

type InfoLogPayload struct {
	Message string
}

func createLog(entry logging.Entry) {
	ctx := context.Background()
	projectID, err := metadata.ProjectID()
	if err != nil {
		log.Fatalf("Failed to fetch project ID: %v", err)
	}

	// Creates a client.
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	logName := "google-cloud-logger"
	logger := client.Logger(logName)

	// Adds an entry to the log buffer.
	logger.Log(entry)

	// Closes the client and flushes the buffer to the Cloud Logging service.
	if err := client.Close(); err != nil {
		log.Fatalf("Failed to close client: %v", err)
	}
}

func LogInfo(message string) {
	entry := InfoLogPayload{Message: message}
	if !metadata.OnGCE() {
		log.Println(entry)
		return
	} else {
		createLog(
			logging.Entry{
				Payload:  entry,
				Severity: logging.Info,
			},
		)
	}
}

func LogError(message string, errorData string) {
	entry := ErrorLogPayload{
		Message: message,
		Error:   errorData,
	}

	if !metadata.OnGCE() {
		log.Println(entry)
		return
	} else {
		createLog(
			logging.Entry{
				Payload:  entry,
				Severity: logging.Error,
			},
		)
	}
}
