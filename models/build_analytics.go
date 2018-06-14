package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

// BuildAnalytics ...
type BuildAnalytics struct {
	gorm.Model
	AppID       string          `db:"app_id" json:"app_id"`
	StackID     string          `db:"stack_id" json:"stack_id"`
	Platform    string          `db:"platform" json:"platform"`
	CLIVersion  string          `db:"cli_version" json:"cli_version"`
	Status      string          `db:"status" json:"status"`
	StartTime   time.Time       `db:"start_time" json:"start_time"`
	Runtime     time.Duration   `db:"run_time" json:"run_time"`
	RawJSONData json.RawMessage `db:"raw_json_data" json:"raw_json_data" sql:"type:json"`

	StepAnalytics []StepAnalytics `gorm:"foreignkey:BuildAnalyticsID" json:"step_analytics"`
}

// Create ...
func (b *BuildAnalytics) Create() {
	CreateInDB(b)
}
