# Pub/Sub notifications to Keybase
Bot that listens to Pub/Sub notifications and publish it to Keybase team#channel.

## Bot

### Build container
```bash
docker build -t keybase-cloud-build-bot .
```

### Run in GCP
```bash
docker run --rm \
    -e KEYBASE_USERNAME="<user name>" \
    -e KEYBASE_PAPERKEY="<keybase paper key>" \
    -e KEYBASE_SERVICE="1" \
    -e PROJECT_ID="<GCP project ID>" \
    -e SUBSCTIPTION_ID="<subscription ID>" \
    -e TEAM_NAME="<keybase team name>" \
    keybase-cloud-build-bot cloud-build-bot
```

optionally you can also pass `CHANNEL` environment variable

### Run with service account key
```bash
docker run --rm \
    -v $PWD/service_account:/service_account \
    -e GOOGLE_APPLICATION_CREDENTIALS="/service_account/<service account file>.json" \
    -e KEYBASE_USERNAME="<user name>" \
    -e KEYBASE_PAPERKEY="<keybase paper key>" \
    -e KEYBASE_SERVICE="1" \
    -e PROJECT_ID="<GCP project ID>" \
    -e SUBSCTIPTION_ID="<subscription ID>" \
    -e TEAM_NAME="<keybase team name>" \
    keybase-cloud-build-bot cloud-build-bot
```

optionally you can also pass `CHANNEL` environment variable

## Helper function

There are ready to use functions to parse cloud-builds notifications and monitoring notifications and re-publish them to the bot topic.

### Cloud Build messages
```
gcloud functions deploy ReportCloudBuild --runtime go113 --trigger-topic cloud-builds --set-env-vars 'NOTIFICATIONS_PROJECT=<>,NOTIFICATIONS_TOPIC=<>'
```

### Monitoring alerts
```
gcloud functions deploy ReportAlerts --runtime go113 --trigger-topic <ALERTS TOPIC> --set-env-vars 'NOTIFICATIONS_PROJECT=<>,NOTIFICATIONS_TOPIC=<>'
```
