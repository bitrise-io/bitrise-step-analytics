package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-team/bitrise-step-analytics/models"
	"github.com/pkg/errors"
)

// StepAnalyticsParams ...
type StepAnalyticsParams struct {
	StepID      string          `json:"step_id"`
	Status      string          `json:"status"`
	StartTime   time.Time       `json:"start_time"`
	Runtime     time.Duration   `json:"run_time"`
	RawJSONData json.RawMessage `json:"raw_json_data"`
}

// BuildAnalyticsParams ...
type BuildAnalyticsParams struct {
	AppID         string                `json:"app_id"`
	StackID       string                `json:"stack_id"`
	Platform      string                `json:"platform"`
	CLIVersion    string                `json:"cli_version"`
	Status        string                `json:"status"`
	StartTime     time.Time             `json:"start_time"`
	Runtime       time.Duration         `json:"run_time"`
	RawJSONData   json.RawMessage       `json:"raw_json_data"`
	StepAnalytics []StepAnalyticsParams `json:"step_analytics"`
}

// AnalyticsLogHandler ...
func AnalyticsLogHandler(w http.ResponseWriter, r *http.Request) error {
	params := BuildAnalyticsParams{}
	defer httpresponse.RequestBodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf(" [!] Exception: Internal Server Error: AnalyticsLogHandler: %+v", errors.Wrap(err, "Failed to JSON decode request body"))
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	buildAnalytics := models.BuildAnalytics{
		AppID:       params.AppID,
		StackID:     params.StackID,
		Platform:    params.Platform,
		CLIVersion:  params.CLIVersion,
		Status:      params.Status,
		StartTime:   params.StartTime,
		Runtime:     params.Runtime,
		RawJSONData: params.RawJSONData,
	}

	stepAnalyticsList := []models.StepAnalytics{}

	for _, aStepAnalyticsParam := range params.StepAnalytics {
		stepAnalytics := models.StepAnalytics{
			StepID:      aStepAnalyticsParam.StepID,
			Status:      aStepAnalyticsParam.Status,
			StartTime:   aStepAnalyticsParam.StartTime,
			Runtime:     aStepAnalyticsParam.Runtime,
			RawJSONData: aStepAnalyticsParam.RawJSONData,
		}
		stepAnalyticsList = append(stepAnalyticsList, stepAnalytics)
	}

	buildAnalytics.StepAnalytics = stepAnalyticsList
	fmt.Println(buildAnalytics)

	return httpresponse.RespondWithSuccess(w, buildAnalytics)
}
