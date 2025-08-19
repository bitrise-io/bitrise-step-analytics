package stepmetrics

import (
	"fmt"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/bitrise-io/bitrise-step-analytics/models"
)

func CreateMetricsFromEvent(client *statsd.Client, event models.TrackEvent) {
	_ = client.Incr("events", []string{}, 1)

	switch event.EventName {
		case "step_appstoreconnect_request":
			appStoreConnectRequests(client, event)
		case "step_appstoreconnect_auth_error":
			appStoreConnectAuthErrors(client)
	}
}

// Schema: https://github.com/bitrise-io/data-events/blob/main/schemas/data-team-229312/events/step_appstoreconnect_request.json
func appStoreConnectRequests(client *statsd.Client, event models.TrackEvent) {
	tags := []string{
		"host:" + event.StringProp("host"),
		"http_method:" + event.StringProp("http_method"),
		"endpoint:" + event.StringProp("endpoint"),
	}
	statusCode := event.IntProp("status_code")
	if statusCode != 0 {
		tags = append(tags, fmt.Sprintf("status_code:%d", statusCode))
	}
	_ = client.Incr("appstoreconnect_requests", tags, 1)
}

// Schema: https://github.com/bitrise-io/data-events/blob/main/schemas/data-team-229312/events/step_appstoreconnect_auth_error.json
func appStoreConnectAuthErrors(client *statsd.Client) {
	_ = client.Incr("appstoreconnect_auth_errors", []string{}, 1)
}
