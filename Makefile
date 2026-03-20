deploy:
	gcloud functions deploy blablappy --entry-point OnPubSubMessage --runtime go116 --trigger-topic blablappy --env-vars-file env.yml --memory 128

update-secrets:
	gcloud functions deploy blablappy --set-secrets SECURE_SLACK_TOKEN=SLACK_TOKEN:latest,SECURE_WORKATODAY_TOKEN=WORKATODAY_TOKEN:latest
