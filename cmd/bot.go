package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(3)
}

func main() {
	var kbLoc string
	var kbc *kbchat.API
	var projectID string
	var subID string
	var teamName string
	var channelPtr *string
	var err error

	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.StringVar(&projectID, "project_id", os.Getenv("PROJECT_ID"), "GCP project ID")
	flag.StringVar(&subID, "subscription_id", os.Getenv("SUBSCTIPTION_ID"), "GCP PubSub subscription name")
	flag.StringVar(&teamName, "team_name", os.Getenv("TEAM_NAME"), "Keybase team name to send message to")
	channelPtr = flag.String("channel", os.Getenv("CHANNEL"), "Keybase channel name to send message to")
	flag.Parse()

	if projectID == "" || subID == "" || teamName == "" {
		fail("All three project_id, subscription_id and team_name must be specified,\n" +
			"using either flags -project_id, -subscription_id, -team_name,\n" +
			"or environmet variables PROJECT_ID, SUBSCTIPTION_ID, TEAM_NAME")
	}

	// ignore empty channel value
	if *channelPtr == "" {
		channelPtr = nil
	}

	if kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc}); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fail("pubsub.NewClient: %v", err)
	}

	sub := client.Subscription(subID)
	cctx, _ := context.WithCancel(ctx)
	err = sub.Receive(cctx, func(ctx context.Context, msg *pubsub.Message) {
		message := string(msg.Data)
		targetTeamName := teamName
		if msg.Attributes["team"] != "" {
			targetTeamName = msg.Attributes["team"]
		}
		targetChannelPtr := channelPtr
		if msg.Attributes["channel"] != "" {
			channel := msg.Attributes["channel"]
			targetChannelPtr = &channel
		}

		fmt.Printf("sending message '%s' to team: '%s'\n", message, targetTeamName)
		if _, err = kbc.SendMessageByTeamName(targetTeamName, targetChannelPtr, message); err != nil {
			fail("Error sending message; %s", err.Error())
		}

		msg.Ack()
	})

	if err != nil {
		fail("Receive: %v", err)
	}
}
