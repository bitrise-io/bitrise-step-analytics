package stepmetrics

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/bitrise-io/bitrise-step-analytics/models"
)

// jsonLogger writes structured lines with no prefix. The package-level log.Println prepends
// "2026/07/31 09:15:27 ", and Datadog only auto-parses a log body that is *entirely* valid
// JSON — a timestamp in front would leave every field as unparsed text and no @build_slug
// facet would exist. Destination is stderr to match the existing log calls in this file,
// whose collection into Datadog is already proven.
var jsonLogger = log.New(os.Stderr, "", 0)

// logJSON writes fields as a single-line JSON log. Datadog parses JSON log bodies into
// attributes, so every key becomes a queryable facet (@build_slug, @step_ref, ...). That is
// why these identifiers belong in logs rather than metric tags: build_slug is unbounded
// cardinality and would blow up the metric, but costs nothing here.
func logJSON(fields map[string]any) {
	// Drop empty strings so a facet's presence is meaningful: a CLI older than v2.42.2 sends
	// no step_execution_id, and emitting it as "" would make @step_execution_id:* match every
	// log line instead of only the ones that can be joined to their failure reason. Only
	// empty strings are dropped — a numeric zero can be a real measurement.
	present := make(map[string]any, len(fields))
	for k, v := range fields {
		if s, ok := v.(string); ok && s == "" {
			continue
		}
		present[k] = v
	}

	line, err := json.Marshal(present)
	if err != nil {
		log.Printf("failed to marshal log fields: %s", err)
		return
	}
	jsonLogger.Println(string(line))
}

func CreateMetricsFromEvent(client *statsd.Client, event models.TrackEvent) {
	_ = client.Incr("events", []string{}, 1)

	switch event.EventName {
	case "step_appstoreconnect_request":
		appStoreConnectRequests(client, event)
	case "step_appstoreconnect_auth_error":
		appStoreConnectAuthErrors(client)
	case "cli_tool_setup":
		cliToolSetup(client, event)
	case "cli_step_activation":
		cliStepActivation(client, event)
	case "step_finished":
		stepFinished(client, event)
	case "step_preparation_failed":
		stepPreparationFailed(event)
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

	isRetry, err := event.BoolProp("is_retry")
	if err == nil {
		tags = append(tags, fmt.Sprintf("is_retry:%t", isRetry))
	} else {
		// TODO: proper WARN log level once logging is fixed
		log.Println("'is_retry' property of 'appstoreconnect_requests' is not a bool")
	}

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

// Schema: https://github.com/bitrise-io/data-events/blob/main/schemas/data-team-229312/events/cli_step_activation.json
func cliStepActivation(client *statsd.Client, event models.TrackEvent) {
	tags := []string{
		"activation_type:" + event.StringProp("activation_type"),
		"cli_version:" + event.StringProp("cli_version"),
	}

	// Which StepLib inventory served the step: "steplib_api" or "steplib". Only
	// steplib refs have one, so it is absent for git and path refs — tagged only
	// when present, to avoid an empty tag value on those.
	inventorySource := event.StringProp("inventory_source")
	if inventorySource != "" {
		tags = append(tags, "inventory_source:"+inventorySource)
	}

	isSuccessful, err := event.BoolProp("is_successful")
	if err == nil {
		tags = append(tags, fmt.Sprintf("is_successful:%t", isSuccessful))
	} else {
		log.Println("'is_successful' property of 'cli_step_activation' is not a bool")
	}

	// A failed activation is logged as well as counted: the metric cannot carry
	// build_slug (unbounded cardinality), so without this there is no way to get from
	// the Datadog failure count to the build that failed. Only on a parsed false —
	// BoolProp yields (false, err) for a missing or non-bool property, and those are
	// unknown outcomes rather than failures.
	if err == nil && !isSuccessful {
		logJSON(failedActivationLogFields(event))
	}

	isCI, err := event.BoolProp("is_ci")
	if err == nil {
		tags = append(tags, fmt.Sprintf("is_ci:%t", isCI))
	} else {
		log.Println("'is_ci' property of 'cli_step_activation' is not a bool")
	}

	_ = client.Incr("cli_step_activation", tags, 1)

	duration := event.IntProp("duration_ms")
	if duration != 0 {
		_ = client.Histogram("cli_step_activation_duration", float64(duration), tags, 1)
	}
}

// failedActivationLogFields is the log payload for a failed activation. Split out from
// cliStepActivation because these key names are the contract the Datadog facets are built
// on: renaming one silently breaks every saved query and dashboard link using it.
func failedActivationLogFields(event models.TrackEvent) map[string]any {
	return map[string]any{
		"event":             "cli_step_activation_failed",
		"build_slug":        event.StringProp("build_slug"),
		"step_execution_id": event.StringProp("step_execution_id"),
		"step_ref":          event.StringProp("step_ref"),
		"activation_type":   event.StringProp("activation_type"),
		"inventory_source":  event.StringProp("inventory_source"),
		"cli_version":       event.StringProp("cli_version"),
		"duration_ms":       event.IntProp("duration_ms"),
	}
}

// Schema: https://github.com/bitrise-io/data-events/blob/main/schemas/data-team-229312/events/step_preparation_failed.json
//
// Logged, not counted. This event carries the error_message that cli_step_activation lacks,
// so joining the two on step_execution_id turns a failure count into a cause. step_finished
// already counts failures, so a second metric here would only duplicate it.
func stepPreparationFailed(event models.TrackEvent) {
	logJSON(preparationFailedLogFields(event))
}

// preparationFailedLogFields is the other half of the drill-down. It has no build_slug —
// the event does not carry one — which is why this is a separate log line joined on
// step_execution_id rather than extra fields on the activation log.
func preparationFailedLogFields(event models.TrackEvent) map[string]any {
	return map[string]any{
		"event":                 "step_preparation_failed",
		"step_execution_id":     event.StringProp("step_execution_id"),
		"workflow_execution_id": event.StringProp("workflow_execution_id"),
		"step_id":               event.StringProp("step_id"),
		"step_version":          event.StringProp("step_version"),
		"error_message":         event.StringProp("error_message"),
	}
}

var trackedSteps = map[string]bool{
	"git-clone":                         true,
	"deploy-to-bitrise-io":              true,
	"activate-ssh-key":				     true,
	"script":                            true,
}

// Schema: https://github.com/bitrise-io/data-events/blob/main/schemas/data-team-229312/events/step_finished.json
func stepFinished(client *statsd.Client, event models.TrackEvent) {
	tags := []string{
		"status:" + event.StringProp("status"),
	}

	// Only tag with step_id for the known, bounded set of tracked steps to keep cardinality in check.
	if stepID := event.StringProp("step_id"); trackedSteps[stepID] {
		tags = append(tags, "step_id:"+stepID)
	}

	_ = client.Incr("step_finished_total", tags, 1)

	runtime := event.IntProp("runtime")
	if runtime != 0 {
		_ = client.Histogram("step_runtime_seconds", float64(runtime), tags, 1)
	}
}
