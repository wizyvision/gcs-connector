package run

import (
	"bytes"
	"cloud-storage-connector/notifications"
	secret "cloud-storage-connector/secret_manager"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"time"

	createTask "cloud-storage-connector/create_task"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

func runGetRequest(url string) error {
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	locationID := os.Getenv("GCLOUD_REGION")
	queueID := os.Getenv("QUEUE_ID")

	if !appengine.IsAppEngine() {
		if _, err := http.Get(url); err != nil {
			log.Fatalln(err)
			return err
		}
	} else {
		task, err := createTask.CreateTask(projectID, locationID, queueID, "", url, true, "")
		if err != nil {
			fmt.Printf("createTask: %v", err)
			return err
		}
		fmt.Printf("Create Task: %s\n", task.GetName())
	}

	return nil
}

func Execute(gcsObject string, gcsBucket string) (string, error) {
	// download file
	// file name and bucket should be from the query parameter
	data, err := downloadFileIntoMemory(gcsBucket, gcsObject)
	if err != nil {
		fmt.Printf("Error downloading file into memory: %v", err)
	}

	// Determine the content type of the file
	dataMimeType := http.DetectContentType(data)
	fmt.Println("dataMimeType", dataMimeType)

	// Create new multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	var fw io.Writer
	if fw, err = createFormFile(writer, gcsObject, dataMimeType); err != nil {
		fmt.Printf("Error creating writer: %v", err)
	}

	// copy the file
	_, err = io.Copy(fw, bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error COPY:", err)
	}

	// Add the addTagIds field
	err = writer.WriteField("addTagIds", os.Getenv("TAG_IDS"))
	if err != nil {
		fmt.Println("Error addTagids value:", err)
	}
	// Add privacyId
	err = writer.WriteField("privacyId", os.Getenv("PRIVACY_ID"))
	if err != nil {
		fmt.Println("Error privacyId value:", err)
	}

	writer.Close()

	url := os.Getenv("WV_PUBLIC_API_URL")

	// Create a new request using http
	req, err := http.NewRequest("POST", url, body)

	setAuth(req)
	setHeaders(req, writer)

	// send request with headers
	client := &http.Client{}
	response, responseErr := client.Do(req)
	if responseErr != nil {
		errMsg := &notifications.Message{
			Pretext: "[Uploader] Failed to upload image.",
			Text:    fmt.Sprintf("Error: %s", responseErr),
		}
		notifications.SendSlackErrorMessage(*errMsg)
	}
	fmt.Println(response)

	status := response.Status

	return status, nil
}

// downloadFileIntoMemory downloads an object.
func downloadFileIntoMemory(bucket, object string) ([]byte, error) {
	ctx := context.Background()
	// var w io.Writer
	serviceAccountKey, err := secret.LoadSecret(ctx, "projects/wv-gcs-connector-dev/secrets/gcs-conn-service-acct/versions/2")
	if err != nil {
		fmt.Printf("error fetching service acct: %s", err)
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsJSON(serviceAccountKey))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %v", object, err)
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		errMsg := &notifications.Message{
			Pretext: "[Uploader] Failed to downlaod file from Google Cloud Storage.",
			Text:    fmt.Sprintf("Error: %s", err),
		}
		notifications.SendSlackErrorMessage(*errMsg)
		return nil, fmt.Errorf("ioutil.ReadAll: %v", err)
	}

	fmt.Printf("Blob %v downloaded.\n", object)
	return data, nil
}

func setAuth(req *http.Request) {
	bearer := "Bearer " + os.Getenv("WV_BEARER_API_KEY")
	req.Header.Add("Authorization", bearer)
}

func setHeaders(req *http.Request, w *multipart.Writer) {
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Add("Accept", "application/json")
}

func createFormFile(w *multipart.Writer, filename, mimeType string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	if !strings.Contains(mimeType, "image") {
		mimeType = "application/octet-stream"
	}
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filename))
	h.Set("Content-Type", mimeType)
	fmt.Println("mimeType:", mimeType)
	return w.CreatePart(h)
}
