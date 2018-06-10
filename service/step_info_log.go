package service

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/slapec93/bitrise-step-analytics/models"
	"github.com/slapec93/bitrise-step-analytics/utils"
)

// StepInfo ...
type StepInfo struct {
	StepName   string     `json:"step_name"`
	Duration   float64    `json:"duration"`
	IsCI       bool       `json:"is_ci"`
	LaunchDate *time.Time `json:"launch_date"`
}

// NewStepInfoFromStepInfoData ...
func NewStepInfoFromStepInfoData(stepInfoData models.StepInfoData) StepInfo {
	return StepInfo{
		StepName:   stepInfoData.StepName,
		Duration:   stepInfoData.Duration,
		IsCI:       stepInfoData.IsCI,
		LaunchDate: stepInfoData.LaunchDate,
	}
}

// LogStepInfoParams ...
type LogStepInfoParams struct {
	StepName   *string  `json:"step_name"`
	Duration   *float64 `json:"duration"`
	IsCI       *bool    `json:"is_ci"`
	LaunchDate *string  `json:"launch_date"`
}

// StepInfoLogHandler ...
func StepInfoLogHandler(w http.ResponseWriter, r *http.Request) error {
	params := LogStepInfoParams{}
	defer utils.RequestBodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf(" [!] Exception: Internal Server Error: AppWebhookCreateHandler: %+v", errors.Wrap(err, "Failed to JSON decode request body"))
		return RespondWithBadRequest(w, "Invalid request body, JSON decode failed")
	}
	stepInfo := models.StepInfoData{}
	stepInfo.StepName = *params.StepName
	stepInfo.Duration = *params.Duration
	stepInfo.IsCI = *params.IsCI
	stepInfo.Create()
	return RespondWithSuccess(w, NewStepInfoFromStepInfoData(stepInfo))
}
