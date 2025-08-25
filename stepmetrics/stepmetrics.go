package stepmetrics

import (
	"fmt"
	"log"
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
	case "cli_tool_setup":
		cliToolSetup(client, event)
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

// Schema: https://github.com/bitrise-io/data-events/blob/main/schemas/data-team-229312/events/cli_tool_setup.json
func cliToolSetup(client *statsd.Client, event models.TrackEvent) {
	tags := []string{
		"provider:" + event.StringProp("provider"),
		"tool_name:" + event.StringProp("tool_name"),
		"cli_version:" + event.StringProp("cli_version"),
	}

	isSuccessful, err := event.BoolProp("is_successful")
	if err == nil {
		tags = append(tags, fmt.Sprintf("is_successful:%t", isSuccessful))
	} else {
		// TODO: proper WARN log level once logging is fixed
		log.Println("'is_successful' property of 'cli_tool_setup' is not a bool")
	}

	isCI, err := event.BoolProp("is_ci")
	if err == nil {
		tags = append(tags, fmt.Sprintf("is_ci:%t", isCI))
	} else {
		// TODO: proper WARN log level once logging is fixed
		log.Println("'is_ci' property of 'cli_tool_setup' is not a bool")
	}

	_ = client.Incr("cli_tool_setup", tags, 1)
}
