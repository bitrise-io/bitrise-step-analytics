package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

// StepAnalytics ...
type StepAnalytics struct {
	gorm.Model
	BuildAnalyticsID uint64          `db:"build_analytics_id" json:"build_analytics_id"`
	StepID           string          `db:"step_id" json:"step_id"`
	Status           string          `db:"status" json:"status"`
	StartTime        time.Time       `db:"start_time" json:"start_time"`
	Runtime          time.Duration   `db:"run_time" json:"run_time"`
	RawJSONData      json.RawMessage `db:"raw_json_data" json:"raw_json_data" sql:"type:json"`
}

// Create ...
func (s *StepAnalytics) Create() {
	CreateInDB(s)
}
