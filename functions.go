package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
)

// User-set environment variables.
var projectID = os.Getenv("NOTIFICATIONS_PROJECT")
var topic = os.Getenv("NOTIFICATIONS_TOPIC")

var client *pubsub.Client

// PubSubMessage is the payload of a Pub/Sub event.
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func init() {
	// err is pre-declared to avoid shadowing client.
	var err error

	// client is initialized with context.Background() because it should
	// persist between function invocations.
	client, err = pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		log.Fatalf("pubsub.NewClient: %v", err)
	}
}

func ReportCloudBuild(ctx context.Context, m PubSubMessage) error {
	build := CloudBuildMessage{}

	err := json.Unmarshal(m.Data, &build)
	if err != nil {
		log.Printf("Error decoding message; %s", err.Error())
		return err
	}

	message := fmt.Sprintf("Project *%s*, Repo: *%s*, Branch: *%s*, Status: *%s*\ncheck build details: %s",
		build.ProjectID,
		build.Substitutions.REPO_NAME,
		build.Substitutions.BRANCH_NAME,
		build.Status, build.LogURL)
	// TODO: pass as env variable
	channel := "builds"
	log.Printf(message)
	err = publish(ctx, message, channel)
	if err != nil {
		return err
	}

	return nil
}

func ReportAlerts(ctx context.Context, m PubSubMessage) error {
	monitoringMessage := MonitoringMessage{}

	err := json.Unmarshal(m.Data, &monitoringMessage)
	if err != nil {
		log.Printf("Error decoding message; %s", err.Error())
		return err
	}

	incident := monitoringMessage.Incident
	endedAt := "-"
	if incident.EndedAt != 0 {
		endedAt = fmt.Sprintf("%s", time.Unix(incident.EndedAt, 0))
	}
	message := fmt.Sprintf(
		":rotating_light: Incident with resource *%s*.\n"+
			"Condition: *%s*\n"+
			"State: *%s*\n"+
			"Started: *%s*, Ended: *%s*\n"+
			"Documentation: %s\n" +
			"Check details: %s",
		incident.ResourceDisplayName,
		incident.ConditionName,
		incident.State,
		time.Unix(incident.StartedAt, 0), endedAt,
		incident.Documentation.Content,
		incident.URL,
	)
	// TODO: pass as env variable
	channel := "devops"
	err = publish(ctx, message, channel)
	if err != nil {
		return err
	}

	return nil
}

func publish(ctx context.Context, message string, channel string) error {
	botMessage := &pubsub.Message{
		Data:       []byte(message),
		Attributes: map[string]string{"channel": channel},
	}

	_, err := client.Topic(topic).Publish(ctx, botMessage).Get(ctx)
	if err != nil {
		log.Printf("topic(%s).Publish.Get: %v", topic, err)
		return fmt.Errorf("error publishing message: %v", err)
	}

	return nil
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

type MonitoringMessage struct {
	Incident struct {
		IncidentID   string `json:"incident_id"`
		ResourceID   string `json:"resource_id"`
		ResourceName string `json:"resource_name"`
		Resource     struct {
			Type string `json:"type"`
			//Labels struct {
			//	SubscriptionID string `json:"subscription_id"`
			//} `json:"labels"`
		} `json:"resource"`
		ResourceDisplayName     string `json:"resource_display_name"`
		ResourceTypeDisplayName string `json:"resource_type_display_name"`
		Metric                  struct {
			Type        string `json:"type"`
			DisplayName string `json:"displayName"`
		} `json:"metric"`
		StartedAt     int64 `json:"started_at"`
		EndedAt       int64 `json:"ended_at"`
		PolicyName    string    `json:"policy_name"`
		ConditionName string    `json:"condition_name"`
		Condition     struct {
			Name               string `json:"name"`
			DisplayName        string `json:"displayName"`
			ConditionThreshold struct {
				Filter string `json:"filter"`
				//Aggregations []struct {
				//	AlignmentPeriod  string `json:"alignmentPeriod"`
				//	PerSeriesAligner string `json:"perSeriesAligner"`
				//} `json:"aggregations"`
				Comparison     string `json:"comparison"`
				ThresholdValue int    `json:"thresholdValue"`
				Duration       string `json:"duration"`
				Trigger        struct {
					Count int `json:"count"`
				} `json:"trigger"`
			} `json:"conditionThreshold"`
		} `json:"condition"`
		URL           string `json:"url"`
		Documentation struct {
			Content  string `json:"content"`
			MimeType string `json:"mime_type"`
		} `json:"documentation"`
		State   string `json:"state"`
		Summary string `json:"summary"`
	} `json:"incident"`
	Version string `json:"version"`
}
