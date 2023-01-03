# wv-connector


## Deploy via cloud run button

[![Run on Google Cloud](https://deploy.cloud.run/button.svg)](https://deploy.cloud.run)

After successful deployment, navigate to *GCP IAM & Admin* and find the default compute service account created.

<img width="978" alt="Screen Shot 2022-12-21 at 9 39 25 AM" src="https://user-images.githubusercontent.com/35460203/208800600-bcd461e5-63ee-4678-82d1-be8a2006249c.png">

Add these roles to the default service account:
- Storage Object Viewer
- Cloud Task Enqueuer
- Cloud Run Developer


## Deploy manually

1. Clone this repository.

```
git clone https://github.com/sephdiza/gcs-connector.git
```
2. Install [docker](https://docs.docker.com/get-docker).
3. Create a file named `.env` at the root of this folder. Copy the contents of `placeholder.env` and set the values.
```
# GCP
GOOGLE_CLOUD_PROJECT='wv-gcs-connector-dev'

# Path where you put the service account file
GOOGLE_APPLICATION_CREDENTIALS_PATH="./secrets/service-account.json"

# WizyVision
# The application API key
WV_BEARER_API_KEY=""

# Public API url of your account
WV_PUBLIC_API_URL="https://mywizyvisiondomain.eu.wizyvision.app/api/v1/public/files"

# Comma-separated list of Tag IDs to set for the uploaded files
TAG_IDS=""

# Privacy to set for the uploaded files, Standard if this is left empty
PRIVACY_ID=""

# Header and token to be authenticated to call the Uploader service
UPLOADER_SERVICE_AUTH_HEADER="Wizdam-Dev-Csc-Token"
UPLOADER_SERVICE_AUTH_TOKEN=""

# Slack used for notification when errors are encountered. Can be left blank
SLACK_BEARER_TOKEN=""
SLACK_CHANNEL_ID=""
```
4. Build the image.
```
docker build -t IMAGE_NAME .
```
5. Push the image to the container registry.
```
docker tag IMAGE_NAME gcr.io/PROJECT_NAME_HERE/IMAGE_NAME
docker push gcr.io/PROJECT_NAME_HERE/IMAGE_NAME 
```
6. Deploy to cloud run.
```
gcloud run deploy IMAGE_NAME \
  --project PROJECT_NAME \
  --image IMAGE_URL \
  --region REGION
```
7. Take note of the URL of the service after successfully deploying. This will be needed in the Trigger service deployment.
