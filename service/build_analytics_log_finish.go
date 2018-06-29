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

// BuildAnalyticsLogFinishHandler ...
func BuildAnalyticsLogFinishHandler(w http.ResponseWriter, r *http.Request) error {
	params := BuildAnalyticsParams{}
	defer utils.RequestBodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf(" [!] Exception: Internal Server Error: AnalyticsLogHandler: %+v", errors.Wrap(err, "Failed to JSON decode request body"))
		return RespondWithBadRequest(w, "Invalid request body, JSON decode failed")
	}
	urlVars := mux.Vars(r)
	buildAnalyticsID, err := strconv.Atoi(urlVars["id"])
	if err != nil {
		return RespondWithBadRequest(w, "The provided ID is invalid")
	}

	buildAnalytics := models.FindBuildAnalyticsByID(int64(buildAnalyticsID))
	updateAttributes := models.BuildAnalytics{
		Status:  *params.Status,
		Runtime: time.Duration(*params.Runtime),
	}
	if params.RawJSONData != nil {
		updateAttributes.RawJSONData = *params.RawJSONData
	}
	buildAnalytics.Update(&updateAttributes)

	return RespondWithSuccess(w, buildAnalytics)
}
