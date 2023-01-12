# GCS Connector - Trigger Function

> This directory contains the **Trigger** app that is deployed in Cloud Function. To run the whole Storage connector, you will also need the [**Uploader**](/cmd/uploader) to be deployed in Cloud Run.


The GCS Connector enables automated, real-time transfer of all files added to the designated source GCS bucket to WizyVision,


![gcs connector](https://user-images.githubusercontent.com/4800851/211000369-70e9be5f-36a6-4e60-8232-f6b73d892d8b.png)


## Installation
1. Install Go 1.17
2. Clone this repo
3. Run `go mod tidy`
4. Create a file named `.env.yaml` at the root of this folder. Copy the contents of `placeholder.env.yaml` and set the values

```
# Google Cloud project ID
GOOGLE_CLOUD_PROJECT: "INPUT_YOUR_GOOGLE_CLOUD_PROJECT_HERE"

# Region where this service will be deployed
GCLOUD_REGION: "europe-west1"

# Cloud storage bucket name where the files will be uploaded
BUCKET_NAME: "INPUT_YOUR_GOOGLE_CLOUD_STORAGE_HERE"

# Uploader service endpoint
UPLOADER_SERVICE: "<INPUT_SERVICE_ENDPOINT_HERE>/run"

# Header and token to be authenticated to call the Uploader service
UPLOADER_SERVICE_AUTH_HEADER: "Wizdam-Dev-Csc-Token"
UPLOADER_SERVICE_AUTH_TOKEN: ""

# Slack used for notification when errors are encountered. Can be left blank
SLACK_BEARER_TOKEN: ""
SLACK_CHANNEL_ID: ""
```

2. To deploy, run
```
sh deploy.sh
```


## Testing the Connector
> Both the Uploader service and the Trigger function should be deployed
1. Upload files on the GCS bucket
2. In WizyVision webapp, you should see the file uploaded there with the correct tags and privacy set during the deployment
