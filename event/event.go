package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/bitrise-io/bitrise-step-analytics/models"
	"github.com/bitrise-io/bitrise-step-analytics/stepmetrics"
	"google.golang.org/api/option"
)

type Tracker interface {
	Send(analytics models.TrackEvent) error
	HealthCheck() error
}

type tracker struct {
	topic         *pubsub.Topic
	context       *context.Context
	datadogClient *statsd.Client
}

func NewTracker(projectID string, topic string, credentialJSON string) Tracker {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID, option.WithCredentialsJSON([]byte(credentialJSON)))
	if err != nil {
		panic(fmt.Sprintf("Couldn't start PubSub Client: %s", err.Error()))
	}

	// Do not provide hostname, rely on the defaults set by DD_AGENT_HOST and DD_DOGSTATSD_PORT
	datadogClient, err := statsd.New("", statsd.WithNamespace("step_metrics"))
	if err != nil {
		panic(fmt.Sprintf("Couldn't start Datadog Client: %s", err.Error()))
	}

	return tracker{topic: client.Topic(topic), context: &ctx, datadogClient: datadogClient}
}

func (t tracker) Send(event models.TrackEvent) error {
	stepmetrics.CreateMetricsFromEvent(t.datadogClient, event)

	properties := map[string]interface{}{"id": event.ID, "ts": convertEpochInMicrosecondsToBigQueryTimestampFormat(event.Timestamp), "event_name": event.EventName}
	for k, v := range event.Properties {
		if k == "id" || k == "ts" || k == "event_name" {
			continue
		}
		properties[k] = v
	}
	payload, err := json.Marshal(properties)
	if err != nil {
		return err
	}
	_, err = t.topic.Publish(*t.context, &pubsub.Message{
		Data: payload,
	}).Get(*t.context)
	return err
}

func (t tracker) HealthCheck() error {
	ctx, cancel := context.WithTimeout(*t.context, 5*time.Second)
	defer cancel()
	
	// Test if the topic exists and we can access it
	exists, err := t.topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("failed to check topic existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("PubSub topic does not exist")
	}
	
	return nil
}

func convertEpochInMicrosecondsToBigQueryTimestampFormat(timestamp int64) string {
	t := time.Unix(0, timestamp*int64(time.Microsecond)).In(time.UTC)
	return t.Format("2006-01-02 15:04:05.000000")
}
