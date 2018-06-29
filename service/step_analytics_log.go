package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/slapec93/bitrise-step-analytics/models"
	"github.com/slapec93/bitrise-step-analytics/utils"
)

// StepAnalyticsParams ...
type StepAnalyticsParams struct {
	StepID      *string          `json:"step_id"`
	Status      *string          `json:"status"`
	StartTime   *time.Time       `json:"start_time"`
	Runtime     *int64           `json:"run_time"`
	RawJSONData *json.RawMessage `json:"raw_json_data"`
}

// StepAnalyticsLogHandler ...
func StepAnalyticsLogHandler(w http.ResponseWriter, r *http.Request) error {
	params := StepAnalyticsParams{}
	defer utils.RequestBodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf(" [!] Exception: Internal Server Error: StepAnalyticsLogHandler: %+v", errors.Wrap(err, "Failed to JSON decode request body"))
		return RespondWithBadRequest(w, "Invalid request body, JSON decode failed")
	}
	urlVars := mux.Vars(r)
	buildAnalyticsID, err := strconv.Atoi(urlVars["build-analytics-id"])
	if err != nil {
		return RespondWithBadRequest(w, "The provided ID is invalid")
	}

	stepAnalytics := models.StepAnalytics{
		BuildAnalyticsID: uint64(buildAnalyticsID),
		StepID:           *params.StepID,
		Status:           *params.Status,
		StartTime:        *params.StartTime,
		Runtime:          time.Duration(*params.Runtime),
		RawJSONData:      *params.RawJSONData,
	}
	stepAnalytics.Create()

	return RespondWithSuccess(w, stepAnalytics)
}
