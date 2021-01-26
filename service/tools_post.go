package service

import (
	"encoding/json"
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/bitrise-step-analytics/models"
	"github.com/pkg/errors"
	"net/http"
)

// ToolsPostHandler ...
func ToolsPostHandler(w http.ResponseWriter, r *http.Request) error {
	var toolAnalytics models.ToolAnalytics
	defer httpresponse.RequestBodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&toolAnalytics); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	if toolAnalytics.BuildSlug == "" {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, please provide build_slug")
	}
	if len(toolAnalytics.ToolUsage) < 1 {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, please provide at least one tool_usage")
	}
	for _, toolUsage := range toolAnalytics.ToolUsage {
		if toolUsage.Name == "" {
			return httpresponse.RespondWithBadRequestError(w, "Invalid request body, please fill name for tool_usage")
		}
	}

	segmentClient, err := GetClientFromContext(r.Context())
	if err != nil {
		return errors.WithStack(err)
	}

	for _, toolUsage := range toolAnalytics.ToolUsage {
		toolUsage.BuildSlug = toolAnalytics.BuildSlug
		segmentClient.Track(toolUsage)
	}

	return httpresponse.RespondWithSuccess(w, map[string]string{"message": "ok"})
}
