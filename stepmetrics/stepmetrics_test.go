package stepmetrics

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
	"testing"

	"github.com/bitrise-io/bitrise-step-analytics/models"
	"github.com/stretchr/testify/require"
)

// The log-field tests below assert exact maps rather than substrings: the keys become Datadog
// log facets (@build_slug, @step_execution_id, ...), so renaming one silently breaks every
// saved query and dashboard link built on it. require.Equal on the whole map also catches an
// accidentally added or dropped key, which a per-key assertion would not.

func Test_failedActivationLogFields(t *testing.T) {
	event := models.TrackEvent{
		EventName: "cli_step_activation",
		Properties: map[string]any{
			"build_slug":        "45d53f84-351d-488b-a22b-59794f759531",
			"step_execution_id": "7f3a1c92-0000-4aaa-bbbb-ccccdddd0001",
			"step_ref":          "deploy-to-bitrise-io",
			"activation_type":   "steplib_source",
			"inventory_source":  "steplib_api",
			"cli_version":       "2.42.2",
			"duration_ms":       float64(217), // JSON numbers decode as float64
			"is_successful":     false,
		},
	}

	require.Equal(t, map[string]any{
		"event":             "cli_step_activation_failed",
		"build_slug":        "45d53f84-351d-488b-a22b-59794f759531",
		"step_execution_id": "7f3a1c92-0000-4aaa-bbbb-ccccdddd0001",
		"step_ref":          "deploy-to-bitrise-io",
		"activation_type":   "steplib_source",
		"inventory_source":  "steplib_api",
		"cli_version":       "2.42.2",
		"duration_ms":       217,
	}, failedActivationLogFields(event))
}

// A CLI older than v2.42.2 sends neither step_execution_id nor inventory_source. Those
// activations are still worth logging: build_slug is present, so the failure is still
// traceable to a build even though it cannot be joined to its reason.
func Test_failedActivationLogFields_missingPropertiesAreEmpty(t *testing.T) {
	event := models.TrackEvent{
		EventName: "cli_step_activation",
		Properties: map[string]any{
			"build_slug": "af4a4b03-40af-4b4a-a637-5b9ac28432fe",
			"step_ref":   "deploy-to-bitrise-io",
		},
	}

	fields := failedActivationLogFields(event)

	require.Equal(t, "af4a4b03-40af-4b4a-a637-5b9ac28432fe", fields["build_slug"])
	require.Equal(t, "", fields["step_execution_id"])
	require.Equal(t, "", fields["inventory_source"])
	require.Equal(t, 0, fields["duration_ms"])
}

func Test_preparationFailedLogFields(t *testing.T) {
	event := models.TrackEvent{
		EventName: "step_preparation_failed",
		Properties: map[string]any{
			"step_execution_id":     "7f3a1c92-0000-4aaa-bbbb-ccccdddd0001",
			"workflow_execution_id": "9b2e4d51-0000-4eee-ffff-000011112222",
			"step_id":               "deploy-to-bitrise-io",
			"step_version":          "1",
			"error_message":         "Preparing Step (deploy-to-bitrise-io@1) failed: resolve step version",
		},
	}

	require.Equal(t, map[string]any{
		"event":                 "step_preparation_failed",
		"step_execution_id":     "7f3a1c92-0000-4aaa-bbbb-ccccdddd0001",
		"workflow_execution_id": "9b2e4d51-0000-4eee-ffff-000011112222",
		"step_id":               "deploy-to-bitrise-io",
		"step_version":          "1",
		"error_message":         "Preparing Step (deploy-to-bitrise-io@1) failed: resolve step version",
	}, preparationFailedLogFields(event))
}

func Test_normalizeEndpoint(t *testing.T) {
	tests := []struct {
		endpoint string
		want     string
	}{
		{
			endpoint: "/v1/certificates",
			want:     "/v1/certificates",
		},
		{
			endpoint: "/v1/devices",
			want:     "/v1/devices",
		},
		{
			endpoint: "/v1/profiles",
			want:     "/v1/profiles",
		},
		{
			endpoint: "/v1/bundleids",
			want:     "/v1/bundleids",
		},
		{
			endpoint: "/v1/profiles/qx2rdp4p9r",
			want:     "/v1/profiles/{id}",
		},
		{
			endpoint: "/v1/profiles/22x3m5y2l3/bundleid",
			want:     "/v1/profiles/{id}/bundleid",
		},
		{
			endpoint: "/v1/profiles/22x3m5y2l3/certificates",
			want:     "/v1/profiles/{id}/certificates",
		},
		{
			endpoint: "/v1/bundleids/22x3m5y2l3/bundleidcapabilities",
			want:     "/v1/bundleids/{id}/bundleidcapabilities",
		},
		{
			endpoint: "/v1/certificates/abc123def456",
			want:     "/v1/certificates/{id}",
		},
		{
			endpoint: "/v1/devices/xyz789abc123",
			want:     "/v1/devices/{id}",
		},
		{
			endpoint: "/v1/certificates/1234567890",
			want:     "/v1/certificates/{id}",
		},
		{
			endpoint: "/v2/profiles/someid",
			want:     "/v2/profiles/someid",
		},
		{
			endpoint: "/v1/bundleids/verylongid123456/bundleidcapabilities",
			want:     "/v1/bundleids/{id}/bundleidcapabilities",
		},
		{
			endpoint: "",
			want:     "",
		},
		{
			endpoint: "/v1/profiles/abc-123_def456",
			want:     "/v1/profiles/{id}",
		},
	}

	for _, tt := range tests {
		t.Run(strings.ReplaceAll(tt.want, "/", "_"), func(t *testing.T) {
			got := normalizeEndpoint(tt.endpoint)
			require.Equal(t, tt.want, got)
		})
	}
}

// The whole drill-down depends on Datadog auto-parsing the body as JSON, which requires the
// emitted line to be valid JSON with nothing prepended — log.Println would prefix a
// timestamp and silently defeat it. Empty strings are dropped so that @step_execution_id:*
// selects only the failures that can be joined to a reason.
func Test_logJSON_emitsBareJSONAndOmitsEmpty(t *testing.T) {
	var buf bytes.Buffer
	orig := jsonLogger
	jsonLogger = log.New(&buf, "", 0)
	defer func() { jsonLogger = orig }()

	logJSON(failedActivationLogFields(models.TrackEvent{
		Properties: map[string]any{"build_slug": "abc", "duration_ms": float64(0)},
	}))

	line := strings.TrimRight(buf.String(), "\n")
	require.True(t, strings.HasPrefix(line, "{"), "must be bare JSON, got: %q", line)

	var parsed map[string]any
	require.NoError(t, json.Unmarshal([]byte(line), &parsed), "must be valid JSON: %q", line)

	require.Equal(t, "cli_step_activation_failed", parsed["event"])
	require.Equal(t, "abc", parsed["build_slug"])
	// Absent, not empty — this is what makes the facet filter usable.
	require.NotContains(t, parsed, "step_execution_id")
	require.NotContains(t, parsed, "inventory_source")
	// A zero measurement is kept: only empty strings are dropped.
	require.Equal(t, float64(0), parsed["duration_ms"])
}
