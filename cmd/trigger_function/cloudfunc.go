package cloudfunc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"test.com/gocloudfunc/logger"
	"test.com/gocloudfunc/notifications"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
)

func init() {
	// Register a CloudEvent function with the Functions Framework
	functions.CloudEvent("gcsConnectorTrigger", gcsConnectorTrigger)
}

type Data struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Bucket         string `json:"bucket"`
	Link           string `json:"mediaLink"`
	Size           string `json:"size"`
	ContentType    string `json:"contentType"`
	ComponentCount string `json:"componentCount"`
}

// Function myCloudEventFunction accepts and handles a CloudEvent object
func gcsConnectorTrigger(ctx context.Context, event event.Event) error {
	data := &Data{}
	if err := event.DataAs(data); err != nil {
		logger.LogError("Got Data Error", err.Error())
	}
	fmt.Printf("ID: %+v \nName: %+v \n MediaLink: %+v \n Size: %+v \n Bucket: %+v \n ContentType: %+v \n ComponentCount: %+v \n",
		data.ID,
		data.Name,
		data.Link,
		data.Size,
		data.Bucket,
		data.ContentType,
		data.ComponentCount,
	)

	uploaderServiceUrl := os.Getenv("UPLOADER_SERVICE")

	url := fmt.Sprintf("%s?gcsObject=%s&gcsBucket=%s", uploaderServiceUrl, url.QueryEscape(data.Name), data.Bucket)
	fmt.Println("URL: ", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.LogError("Error creating request: ", err.Error())
	}
	setAuth(req)

	response, err := client.Do(req)
	// Success is indicated with 2xx status codes:
	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		errMsg := fmt.Sprintf("%v", response)
		logger.LogError(fmt.Sprintf("Non-OK HTTP status: %v", response.StatusCode), errMsg)
		return errors.New(errMsg)
	}

	if err != nil {
		logger.LogError("Error calling request: ", err.Error())
		errMsg := &notifications.Message{
			Pretext: "[Trigger] Failed on creating task.",
			Text: fmt.Sprintf("Error: %s \nID: %+v \nName: %+v \n MediaLink: %+v \n Size: %+v \n Bucket: %+v \n ContentType: %+v \n ComponentCount: %+v \n",
				err,
				data.ID,
				data.Name,
				data.Link,
				data.Size,
				data.Bucket,
				data.ContentType,
				data.ComponentCount,
			),
		}
		notifications.SendSlackErrorMessage(*errMsg)
		fmt.Printf("createTask: %v", err)
		return err
	}
	fmt.Println(response)

	return nil
}

func setAuth(req *http.Request) {
	req.Header.Add(os.Getenv("UPLOADER_SERVICE_AUTH_HEADER"), os.Getenv("UPLOADER_SERVICE_AUTH_TOKEN"))
}
