package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

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

		build := CloudBuildMessage{}

		err = json.Unmarshal(msg.Data, &build)
		if err != nil {
			fail("Error decoding message; %s", err.Error())
		}

		message := fmt.Sprintf("Project *%s*, Repo: *%s*, Branch: *%s*, Status: *%s*\ncheck build details: %s",
			build.ProjectID,
			build.Substitutions.REPO_NAME,
			build.Substitutions.BRANCH_NAME,
			build.Status, build.LogURL)

		fmt.Printf("sending message '%s' to team: '%s'\n", message, teamName)
		if _, err = kbc.SendMessageByTeamName(teamName, channelPtr, message); err != nil {
			fail("Error sending message; %s", err.Error())
		}

		msg.Ack()
	})

	if err != nil {
		fail("Receive: %v", err)
	}
}

// copied from https://github.com/Jleagle/cloud-build-slack
type CloudBuildMessage struct {
	ID        string `json:"id"`
	ProjectID string `json:"projectId"`
	Status    string `json:"status"`
	Source    struct {
		RepoSource struct {
			ProjectID  string `json:"projectId"`
			RepoName   string `json:"repoName"`
			BranchName string `json:"branchName"`
		} `json:"repoSource"`
	} `json:"source"`
	Steps []struct {
		Name   string   `json:"name"`
		Args   []string `json:"args"`
		Timing struct {
			StartTime time.Time `json:"startTime"`
			EndTime   time.Time `json:"endTime"`
		} `json:"timing"`
		PullTiming struct {
			StartTime time.Time `json:"startTime"`
			EndTime   time.Time `json:"endTime"`
		} `json:"pullTiming"`
		Status string `json:"status"`
	} `json:"steps"`
	Results struct {
		Images []struct {
			Name       string `json:"name"`
			Digest     string `json:"digest"`
			PushTiming struct {
				StartTime time.Time `json:"startTime"`
				EndTime   time.Time `json:"endTime"`
			} `json:"pushTiming"`
		} `json:"images"`
		BuildStepImages []string `json:"buildStepImages"`
	} `json:"results"`
	CreateTime time.Time `json:"createTime"`
	StartTime  time.Time `json:"startTime"`
	FinishTime time.Time `json:"finishTime"`
	Timeout    string    `json:"timeout"`
	Images     []string  `json:"images"`
	Artifacts  struct {
		Images []string `json:"images"`
	} `json:"artifacts"`
	LogsBucket       string `json:"logsBucket"`
	SourceProvenance struct {
		ResolvedRepoSource struct {
			ProjectID string `json:"projectId"`
			RepoName  string `json:"repoName"`
			CommitSha string `json:"commitSha"`
		} `json:"resolvedRepoSource"`
	} `json:"sourceProvenance"`
	BuildTriggerID string `json:"buildTriggerId"`
	Options        struct {
		SubstitutionOption string `json:"substitutionOption"`
		Logging            string `json:"logging"`
	} `json:"options"`
	LogURL        string `json:"logUrl"`
	Substitutions struct {
		BRANCH_NAME string `json:"BRANCH_NAME"`
		COMMIT_SHA  string `json:"COMMIT_SHA"`
		REPO_NAME   string `json:"REPO_NAME"`
		REVISION_ID string `json:"REVISION_ID"`
		SHORT_SHA   string `json:"SHORT_SHA"`
	}
	Tags   []string `json:"tags"`
	Timing struct {
		BUILD struct {
			StartTime time.Time `json:"startTime"`
			EndTime   time.Time `json:"endTime"`
		} `json:"BUILD"`
		FETCHSOURCE struct {
			StartTime time.Time `json:"startTime"`
			EndTime   time.Time `json:"endTime"`
		} `json:"FETCHSOURCE"`
		PUSH struct {
			StartTime time.Time `json:"startTime"`
			EndTime   time.Time `json:"endTime"`
		} `json:"PUSH"`
	} `json:"timing"`
}
