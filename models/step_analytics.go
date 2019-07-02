package models

import (
	"fmt"
	"time"
)

// StepAnalytics ...
type StepAnalytics struct {
	StepID      string        `json:"step_id"`
	StepTitle   *string       `json:"step_title"`
	StepVersion string        `json:"step_verion"`
	StepSource  *string       `json:"step_source"`
	Status      string        `json:"status"`
	StartTime   time.Time     `json:"start_time"`
	Runtime     time.Duration `json:"run_time"`
}

// GetProfileName ...
func (a StepAnalytics) GetProfileName() string {
	return "step"
}

// GetTagArray ...
func (a StepAnalytics) GetTagArray() []string {
	return []string{
		fmt.Sprintf("step_id:%s", a.StepID),
		fmt.Sprintf("status:%s", a.Status),
	}
}

// GetRunTime ...
func (a StepAnalytics) GetRunTime() time.Duration {
	return a.Runtime
}
