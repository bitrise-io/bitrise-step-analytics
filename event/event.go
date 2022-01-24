package event

import (
	"github.com/bitrise-io/bitrise-step-analytics/models"
	segment "gopkg.in/segmentio/analytics-go.v3"
)

// Tracker ...
type Tracker interface {
	Send(analytics models.TrackEvent) error
}

type tracker struct {
	segmentClient segment.Client
}

func NewTracker(writeKey string) Tracker {
	return tracker{segmentClient: segment.New(writeKey)}
}

func (t tracker) Send(analytics models.TrackEvent) error {
	properties := map[string]interface{}{"id": analytics.ID, "timestamp": analytics.Timestamp}
	for k, v := range analytics.Properties {
		properties[k] = v
	}
	return t.segmentClient.Enqueue(segment.Track{
		Event:      analytics.EventName,
		UserId:     "mobile-devops",
		Properties: properties,
	})
}
