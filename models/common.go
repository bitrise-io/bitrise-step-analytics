package models

import "github.com/slapec93/bitrise-step-analytics/database"

// GetModelList ...
func GetModelList() []interface{} {
	return []interface{}{
		StepAnalytics{},
		BuildAnalytics{},
	}
}

// CreateInDB ...
func CreateInDB(object interface{}) {
	db := database.GetDB()
	if db.NewRecord(object) {
		db.Create(object)
	}
}
