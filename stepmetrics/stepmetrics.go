package stepmetrics

import (
	"fmt"
	"regexp"

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
	}
	statusCode := event.IntProp("status_code")
	if statusCode != 0 {
		tags = append(tags, fmt.Sprintf("status_code:%d", statusCode))
	}

	endpoint := normalizeEndpoint(event.StringProp("endpoint"))
	tags = append(tags, "endpoint:"+endpoint)

	_ = client.Incr("appstoreconnect_requests", tags, 1)

	duration := event.IntProp("duration_ms")
	if duration != 0 {
		_ = client.Histogram("appstoreconnect_request_duration", float64(duration), tags, 1)
	}
}

// Schema: https://github.com/bitrise-io/data-events/blob/main/schemas/data-team-229312/events/step_appstoreconnect_auth_error.json
func appStoreConnectAuthErrors(client *statsd.Client) {
	_ = client.Incr("appstoreconnect_auth_errors", []string{}, 1)
}

func normalizeEndpoint(endpoint string) string {
	// Reduce metric cardinality by normalizing unique IDs in endpoint tags
	// Examples:
	// /v1/certificates
	// /v1/devices
	// /v1/profiles
	// /v1/profiles/qx2rdp4p9r
	// /v1/profiles/22x3m5y2l3/bundleid
	// /v1/profiles/22x3m5y2l3/certificates
	// /v1/bundleids/22x3m5y2l3/bundleidcapabilities
	normalizePattern := regexp.MustCompile(`(\/v1\/(?:certificates|devices|profiles|bundleids))\/([a-zA-Z0-9-_]+)`)
	return normalizePattern.ReplaceAllString(endpoint, "$1/{id}")
}
