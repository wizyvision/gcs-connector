export $(grep -v '^#' .env | xargs)
gcloud app deploy app.yaml --project=$GOOGLE_CLOUD_PROJECT