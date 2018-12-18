package models

import (
	"encoding/json"
	"time"
)

// StepAnalytics ...
type StepAnalytics struct {
	StepID      string          `json:"step_id"`
	Status      string          `json:"status"`
	StartTime   time.Time       `json:"start_time"`
	Runtime     time.Duration   `json:"run_time"`
	RawJSONData json.RawMessage `json:"raw_json_data"`
}
