package models

import (
	"fmt"
	"time"
)

// BuildAnalytics ...
type BuildAnalytics struct {
	AppSlug      string        `json:"app_slug"`
	BuildSlug    string        `json:"build_slug"`
	StackID      string        `json:"stack_id"`
	Platform     string        `json:"platform"`
	CLIVersion   string        `json:"cli_version"`
	Status       string        `json:"status"`
	StartTime    time.Time     `json:"start_time"`
	Runtime      time.Duration `json:"run_time"`
	RepositoryID string        `json:"repo_id"`
	WorkflowName string        `json:"workflow_name"`

	StepAnalytics []StepAnalytics `json:"step_analytics"`
}

// GetProfileName ...
func (a BuildAnalytics) GetProfileName() string {
	return "build"
}

// GetTagArray ...
func (a BuildAnalytics) GetTagArray() []string {
	return []string{
		fmt.Sprintf("app_slug:%s", a.AppSlug),
		fmt.Sprintf("build_slug:%s", a.BuildSlug),
		fmt.Sprintf("stack_id:%s", a.StackID),
		fmt.Sprintf("platform:%s", a.Platform),
		fmt.Sprintf("cli_version:%s", a.CLIVersion),
		fmt.Sprintf("status:%s", a.Status),
		fmt.Sprintf("repository_id:%s", a.RepositoryID),
		fmt.Sprintf("workflow_name:%s", a.Status),
	}
}

// GetRunTime ...
func (a BuildAnalytics) GetRunTime() time.Duration {
	return a.Runtime
}
