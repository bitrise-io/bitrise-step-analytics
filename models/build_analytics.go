package models

import (
	"encoding/json"
	"fmt"
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

// GetProfileName ...
func (a BuildAnalytics) GetProfileName() string {
	return "build"
}

// GetTagArray ...
func (a BuildAnalytics) GetTagArray() []string {
	return []string{
		fmt.Sprintf("app_id:%s", a.AppID),
		fmt.Sprintf("stack_id:%s", a.StackID),
		fmt.Sprintf("cli_version:%s", a.CLIVersion),
		fmt.Sprintf("status:%s", a.Status),
	}
}

// GetRunTime ...
func (a BuildAnalytics) GetRunTime() time.Duration {
	return a.Runtime
}
