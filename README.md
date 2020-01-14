# Cloud build notifications to Keybase
Bot that listens to GCP cloud build notifications and publish it to Keybase team.

## Build container
```bash
docker build -t keybase-cloud-build-bot .
```

## Run in GCP
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

## Run with service account key
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
