package service

import (
	"net/http"
	"time"

	"github.com/slapec93/bitrise-step-analytics/models"
)

// StepInfoAnalytics ...
type StepInfoAnalytics struct {
	LaunchDate time.Time `json:"launch_date"`
	Launches   int64     `json:"launches"`
	ErrorRate  float64   `json:"error_rate"`
}

// StepInfoAnalyticsHandler ...
func StepInfoAnalyticsHandler(w http.ResponseWriter, r *http.Request) error {
	stepInfos := models.ListStepInfos()

	return RespondWithSuccess(w, NewStepInfosFromStepInfoDatas(stepInfos))
}
