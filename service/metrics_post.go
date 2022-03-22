package service

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/bitrise-step-analytics/models"
	"github.com/pkg/errors"
)

func MetricsPostHandler(w http.ResponseWriter, r *http.Request) error {
	var buildAnalytics models.BuildAnalytics
	defer httprequest.BodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&buildAnalytics); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	if reflect.DeepEqual(buildAnalytics, models.BuildAnalytics{}) {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, please provide metrics data")
	}
	dogstatsd, err := GetClientFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	dogstatsd.Track(buildAnalytics)
	for _, aStepAnalytic := range buildAnalytics.StepAnalytics {
		aStepAnalytic.AppSlug = buildAnalytics.AppSlug
		aStepAnalytic.BuildSlug = buildAnalytics.BuildSlug
		dogstatsd.Track(aStepAnalytic)
	}

	return httpresponse.RespondWithSuccess(w, map[string]string{"message": "ok"})
}
