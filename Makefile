deploy:
	gcloud functions deploy blablappy --entry-point OnPubSubMessage --runtime go113 --trigger-topic blablappy --env-vars-file env.yml --memory 128
