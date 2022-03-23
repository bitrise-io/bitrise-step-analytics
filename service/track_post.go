package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitrise-io/api-utils/httprequest"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/bitrise-step-analytics/models"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

func TrackPostHandler(w http.ResponseWriter, r *http.Request) error {
	var trackAnalytics models.TrackEvent
	defer httprequest.BodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&trackAnalytics); err != nil {
		return httpresponse.RespondWithBadRequestError(w, fmt.Sprintf("Invalid request body, JSON decode failed: %s", err.Error()))
	}

	if trackAnalytics.EventName == "" {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, please provide event's name")
	}

	if _, err := uuid.FromString(trackAnalytics.ID); err != nil {
		return httpresponse.RespondWithBadRequestError(w, fmt.Sprintf("Invalid request body, provided id (%s) is not valid uuid: %s", trackAnalytics.ID, err.Error()))
	}

	tracker, err := GetTrackerFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}
	if err := tracker.Send(trackAnalytics); err != nil {
		return httpresponse.RespondWithError(w, fmt.Sprintf("Couldn't send analytics event: %s", err.Error()), http.StatusInternalServerError)
	}

	return httpresponse.RespondWithSuccess(w, map[string]string{"message": "ok"})
}
