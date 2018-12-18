package models

import (
	"encoding/json"
	"time"
)

// BuildAnalytics ...
type BuildAnalytics struct {
	AppID       string          `json:"app_id"`
	StackID     string          `json:"stack_id"`
	Platform    string          `json:"platform"`
	CLIVersion  string          `json:"cli_version"`
	Status      string          `json:"status"`
	StartTime   time.Time       `json:"start_time"`
	Runtime     time.Duration   `json:"run_time"`
	RawJSONData json.RawMessage `json:"raw_json_data"`

	StepAnalytics []StepAnalytics `json:"step_analytics"`
}
