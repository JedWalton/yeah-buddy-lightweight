## RUN TESTS 
#
# IF TESTS PASS
#
# RUN ANY MIGRATIONS
#
# DEPLOY TO CLOUD RUN
gcloud run deploy SERVICE-NAME \
  --image gcr.io/PROJECT-ID/IMAGE \
  --set-env-vars KEY1=VALUE1,KEY2=VALUE2 \
  --min-instances 0 \
  --max-instances 1 \
  --region europe-west1

