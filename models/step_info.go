package models

import (
	"github.com/jinzhu/gorm"
	"github.com/markbates/pop/nulls"
)

// StepInfo ...
type StepInfo struct {
	gorm.Model
	StepName   string     `db:"step_name"`
	Duration   float64    `db:"duration"`
	DurationCI float64    `db:"duration_ci"`
	LaunchDate nulls.Time `db:"launch_date"`
	CreatedAt  nulls.Time `db:"created_at"`
}
