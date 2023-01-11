. yaml-source.sh .env.yaml

gcloud functions deploy gcsConnectorTrigger \
--runtime go116 \
--env-vars-file .env.yaml \
--trigger-resource $BUCKET_NAME \
--trigger-event google.storage.object.finalize \
--project $GOOGLE_CLOUD_PROJECT \
--region $GCLOUD_REGION