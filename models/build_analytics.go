package models

import (
	"time"
)

// BuildAnalytics ...
type BuildAnalytics struct {
	Status        string          `json:"status" track:"status"`
	AppSlug       string          `json:"app_slug" track:"app_slug"`
	Runtime       time.Duration   `json:"run_time" track:"run_time"`
	StackID       string          `json:"stack_id" track:"stack_id"`
	Platform      string          `json:"platform" track:"platform"`
	BuildSlug     string          `json:"build_slug" track:"build_slug"`
	StartTime     time.Time       `json:"start_time" track:"start_time"`
	CLIVersion    string          `json:"cli_version" track:"cli_version"`
	RepositoryID  string          `json:"repo_id" track:"repository_id"`
	WorkflowName  string          `json:"workflow_name" track:"workflow_name"`
	StepAnalytics []StepAnalytics `json:"step_analytics"`
}

// Event ...
func (a BuildAnalytics) Event() string {
	return "build_finished"
}

// Model ...
func (a BuildAnalytics) Model() interface{} {
	return a
}

// UserID ...
func (a BuildAnalytics) UserID() string {
	return a.AppSlug
}
