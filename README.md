# Google Cloud Storage Connector

The GCS Connector enables automated, real-time transfer of all files added to the designated source GCS bucket to WizyVision.

![gcs connector](https://user-images.githubusercontent.com/4800851/211000369-70e9be5f-36a6-4e60-8232-f6b73d892d8b.png)

## Setup and Requirements
1. Google Cloud project
2. WizyVision account
3. Cloud Storage bucket
4. Enable the APIs required in GCP
```
gcloud services enable \
    cloudfunctions.googleapis.com \
    pubsub.googleapis.com \
    cloudbuild.googleapis.com \
    artifactregistry.googleapis.com \
    containerregistry.googleapis.com \
    run.googleapis.com \
    --quiet
```
5. Grant the pubsub.publisher role to the Cloud Storage service account. This will allow the service account to publish events when images are uploaded into the bucket.
```
SERVICE_ACCOUNT="$(gsutil kms serviceaccount -p GOOGLE_CLOUD_PROJECT_ID_HERE)"

gcloud projects add-iam-policy-binding GOOGLE_CLOUD_PROJECT_ID_HERE \
    --member="serviceAccount:${SERVICE_ACCOUNT}" \
    --role='roles/pubsub.publisher'
```
1. For manual installation you will need [Docker](https://docs.docker.com/get-docker) and Go 1.17.


## Installation
The GCS connector is composed of 2 components, refer to [__Uploader__](cmd/uploader) and [__Trigger Function__](cmd/trigger_function) for instructions on how to deploy.

## Testing the Connector
> Both the Uploader service and the Trigger function should already be deployed
1. Upload files on the GCS bucket
2. In WizyVision webapp, you should see the file uploaded there with the correct tags and privacy set during the deployment
