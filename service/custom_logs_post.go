package service

import (
	"encoding/json"
	"net/http"

	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/bitrise-step-analytics/models"
	"github.com/pkg/errors"
)

// CustomLogsPostHandler ...
func CustomLogsPostHandler(w http.ResponseWriter, r *http.Request) error {
	var log models.RemoteLog
	defer httpresponse.RequestBodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	if len(log.LogLevel) == 0 {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, please provide log_level")
	}
	if len(log.Message) == 0 {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, please provide message")
	}
	if _, ok := log.Data["step_id"].(string); !ok {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, please provide data.step_id")
	}
	if _, ok := log.Data["tag"].(string); !ok {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, please provide data.tag")
	}

	dogstatsd, err := GetClientFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	dogstatsd.Track(log)

	return httpresponse.RespondWithSuccess(w, map[string]string{"message": "ok"})
}
