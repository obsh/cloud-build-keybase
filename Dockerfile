FROM golang:1.13.3-stretch AS builder
WORKDIR /bot
RUN go get -u github.com/keybase/go-keybase-chat-bot/kbchat
RUN go get -u cloud.google.com/go/pubsub
COPY main.go /bot/
RUN go build -o build/cloud-build-bot

FROM keybaseio/client
COPY --from=builder /bot/build/cloud-build-bot /usr/bin/cloud-build-bot
