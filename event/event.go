package event

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/bitrise-io/bitrise-step-analytics/models"
)

// Tracker ...
type Tracker interface {
	Send(analytics models.TrackEvent) error
}

type tracker struct {
	topic   *pubsub.Topic
	context *context.Context
}

func NewTracker(projectID string, topic string) Tracker {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		panic(fmt.Sprintf("Couldn't start PubSub Client: %s", err.Error()))
	}
	return tracker{topic: client.Topic(topic), context: &ctx}
}

func (t tracker) Send(analytics models.TrackEvent) error {
	properties := map[string]interface{}{"id": analytics.ID, "ts": convertEpochInMicrosecondsToBigQueryTimestampFormat(analytics.Timestamp), "event_name": analytics.EventName}
	for k, v := range analytics.Properties {
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

func convertEpochInMicrosecondsToBigQueryTimestampFormat(timestamp int64) string {
	t := time.Unix(0, timestamp*int64(time.Microsecond)).In(time.UTC)
	return t.Format("2006-01-02 15:04:05.000000")
}
