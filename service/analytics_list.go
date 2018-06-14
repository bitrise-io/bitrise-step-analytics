package service

import (
	"net/http"

	"github.com/slapec93/bitrise-step-analytics/models"
)

// AnalyticsListHandler ...
func AnalyticsListHandler(w http.ResponseWriter, r *http.Request) error {
	buildAnalytics := models.ListBuildAnalytics()

	return RespondWithSuccess(w, buildAnalytics)
}
