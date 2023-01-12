package createTask

import (
	"context"
	"fmt"
	"os"

	"github.com/wizyvision/gcs-connector/cmd/uploader/notifications"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
)

// createTask creates a new task in your App Engine queue.
func CreateTask(projectID, locationID, queueID, message string, Url string, isHttpRequest bool, service string) (*cloudtaskspb.Task, error) {
	// Create a new Cloud Tasks client instance.
	// See https://godoc.org/cloud.google.com/go/cloudtasks/apiv2
	ctx := context.Background()
	// serviceAccountKey, err := secret.LoadSecret(ctx, "projects/wv-gcs-connector-dev/secrets/gcs-conn-service-acct/versions/2")

	client, err := cloudtasks.NewClient(
		ctx,
		// option.WithCredentialsJSON(serviceAccountKey),
	)
	if err != nil {
		errMsg := &notifications.Message{
			Pretext: "[Trigger] Failed to create cloud task client.",
			Text:    fmt.Sprintf("Error: %s", err),
		}
		notifications.SendSlackErrorMessage(*errMsg)
		return nil, fmt.Errorf("NewClient: %v", err)
	}
	defer client.Close()

	// Build the Task queue path.
	queuePath := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, locationID, queueID)

	var task *cloudtaskspb.Task

	if isHttpRequest {
		task = &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_HttpRequest{
				HttpRequest: &cloudtaskspb.HttpRequest{
					Headers: map[string]string{
						(os.Getenv("AUTH_HEADER")): os.Getenv("AUTH_TOKEN"),
					},
					HttpMethod: cloudtaskspb.HttpMethod_GET,
					Url:        Url,
				},
			},
		}
	} else {
		task = &cloudtaskspb.Task{
			MessageType: &cloudtaskspb.Task_AppEngineHttpRequest{
				AppEngineHttpRequest: &cloudtaskspb.AppEngineHttpRequest{
					HttpMethod:  cloudtaskspb.HttpMethod_GET,
					RelativeUri: Url,
					AppEngineRouting: &cloudtaskspb.AppEngineRouting{
						Service: "default",
					},
				},
			},
		}
	}

	// Build the Task payload.
	// https://godoc.org/google.golang.org/genproto/googleapis/cloud/tasks/v2#CreateTaskRequest
	req := &cloudtaskspb.CreateTaskRequest{
		Parent: queuePath,
		Task:   task,
	}

	// Add a payload message if one is present.
	if isHttpRequest {
		req.Task.GetHttpRequest().Body = []byte(message)
	} else {
		req.Task.GetAppEngineHttpRequest().Body = []byte(message)
	}

	createdTask, err := client.CreateTask(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks.CreateTask: %v", err)
	}

	return createdTask, nil
}
