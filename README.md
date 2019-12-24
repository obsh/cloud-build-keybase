
## Build container
```bash
docker build -t keybase-cloud-build-bot .
```

## Run
```bash
docker run -it --rm \
    -v $PWD/service_account:/service_account \
    -e GOOGLE_APPLICATION_CREDENTIALS="/service_account/<service account file>.json" \
    -e KEYBASE_USERNAME="<user name>" \
    -e KEYBASE_PAPERKEY="<keybase paper key>" \
    -e KEYBASE_SERVICE="1" \
    -e PROJECT_ID="" \
    -e SUBSCTIPTION_ID="" \
    -e TEAM_NAME="" \
    keybase-cloud-build-bot cloud-build-bot
```