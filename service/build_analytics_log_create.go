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

// BuildAnalyticsParams ...
type BuildAnalyticsParams struct {
	AppID       *string          `json:"app_id"`
	StackID     *string          `json:"stack_id"`
	Platform    *string          `json:"platform"`
	CLIVersion  *string          `json:"cli_version"`
	Status      *string          `json:"status"`
	StartTime   *time.Time       `json:"start_time"`
	Runtime     *int64           `json:"run_time"`
	RawJSONData *json.RawMessage `json:"raw_json_data"`
}

// BuildAnalyticsLogCreateHandler ...
func BuildAnalyticsLogCreateHandler(w http.ResponseWriter, r *http.Request) error {
	params := BuildAnalyticsParams{}
	defer utils.RequestBodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf(" [!] Exception: Internal Server Error: AnalyticsLogHandler: %+v", errors.Wrap(err, "Failed to JSON decode request body"))
		return RespondWithBadRequest(w, "Invalid request body, JSON decode failed")
	}

	buildAnalytics := models.BuildAnalytics{
		AppID:       *params.AppID,
		StackID:     *params.StackID,
		Platform:    *params.Platform,
		CLIVersion:  *params.CLIVersion,
		Status:      *params.Status,
		StartTime:   *params.StartTime,
		RawJSONData: *params.RawJSONData,
	}
	buildAnalytics.Create()

	return RespondWithSuccess(w, buildAnalytics)
}
