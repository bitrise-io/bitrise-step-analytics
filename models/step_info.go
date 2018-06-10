package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/slapec93/bitrise-step-analytics/database"
)

// StepInfoData ...
type StepInfoData struct {
	gorm.Model
	StepName   string     `db:"step_name"`
	Duration   float64    `db:"duration"`
	IsCI       bool       `db:"is_ci"`
	LaunchDate *time.Time `db:"launch_date"`
}

// Create ...
func (s *StepInfoData) Create() {
	db := database.GetDB()
	db.NewRecord(s)
	db.Create(s)
}

// FindByID ..
func (s *StepInfoData) FindByID(id int64) {
	db := database.GetDB()
	db.First(s, id)
}

// ListStepInfos ..
func ListStepInfos() []StepInfoData {
	db := database.GetDB()
	stepInfos := []StepInfoData{}
	db.Find(&stepInfos)
	return stepInfos
}
