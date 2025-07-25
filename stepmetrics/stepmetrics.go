package stepmetrics

import (
	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/bitrise-io/bitrise-step-analytics/models"
)

func CreateMetricsFromEvent(client *statsd.Client, event models.TrackEvent) {
	_ = client.Incr("events", []string{}, 1)

	switch event.EventName {
	}
}
