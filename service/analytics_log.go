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

// StepAnalyticsParams ...
type StepAnalyticsParams struct {
	StepID      *string          `json:"step_id"`
	Status      *string          `json:"status"`
	StartTime   *string          `json:"start_time"`
	Runtime     *int64           `json:"run_time"`
	RawJSONData *json.RawMessage `json:"raw_json_data"`
}

// BuildAnalyticsParams ...
type BuildAnalyticsParams struct {
	AppID         *string               `json:"app_id"`
	StackID       *string               `json:"stack_id"`
	Platform      *string               `json:"platform"`
	CLIVersion    *string               `json:"cli_version"`
	Status        *string               `db:"status" json:"status"`
	StartTime     *string               `json:"start_time"`
	Runtime       *int64                `json:"run_time"`
	RawJSONData   *json.RawMessage      `json:"raw_json_data"`
	StepAnalytics []StepAnalyticsParams `json:"step_analytics"`
}

// AnalyticsLogHandler ...
func AnalyticsLogHandler(w http.ResponseWriter, r *http.Request) error {
	params := BuildAnalyticsParams{}
	defer utils.RequestBodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf(" [!] Exception: Internal Server Error: AnalyticsLogHandler: %+v", errors.Wrap(err, "Failed to JSON decode request body"))
		return RespondWithBadRequest(w, "Invalid request body, JSON decode failed")
	}
	startTime, err := utils.StringTimeToTime(*params.StartTime)
	if err != nil {
		return RespondWithBadRequest(w, err.Error())
	}

	buildAnalytics := models.BuildAnalytics{
		AppID:       *params.AppID,
		StackID:     *params.StackID,
		Platform:    *params.Platform,
		CLIVersion:  *params.CLIVersion,
		Status:      *params.Status,
		StartTime:   startTime,
		Runtime:     time.Duration(*params.Runtime) * time.Second,
		RawJSONData: *params.RawJSONData,
	}
	buildAnalytics.Create()
	stepAnalyticsList := []models.StepAnalytics{}

	for _, aStepAnalyticsParam := range params.StepAnalytics {
		startTime, err := utils.StringTimeToTime(*aStepAnalyticsParam.StartTime)
		if err != nil {
			return RespondWithBadRequest(w, err.Error())
		}
		stepAnalytics := models.StepAnalytics{
			BuildAnalyticsID: uint64(buildAnalytics.ID),
			StepID:           *aStepAnalyticsParam.StepID,
			Status:           *aStepAnalyticsParam.Status,
			StartTime:        startTime,
			Runtime:          time.Duration(*aStepAnalyticsParam.Runtime) * time.Second,
			RawJSONData:      *aStepAnalyticsParam.RawJSONData,
		}
		stepAnalytics.Create()
		stepAnalyticsList = append(stepAnalyticsList, stepAnalytics)
	}

	buildAnalytics.StepAnalytics = stepAnalyticsList

	return RespondWithSuccess(w, buildAnalytics)
}
