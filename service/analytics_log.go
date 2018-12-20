package service

import (
	"encoding/json"
	"net/http"

	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-team/bitrise-step-analytics/metrics"
	"github.com/bitrise-team/bitrise-step-analytics/models"
)

// AnalyticsLogHandler ...
func AnalyticsLogHandler(w http.ResponseWriter, r *http.Request) error {
	buildAnalytics := models.BuildAnalytics{}
	defer httpresponse.RequestBodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&buildAnalytics); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	// logger, loggerSync := utils.GetLogger()
	// defer loggerSync()

	dogstatsd := metrics.NewDogStatsDMetrics("")
	// fmt.Printf(`{"msg":"Build finished"}`)
	// logger.Warn("Build finished",
	// 	zap.String("app_id", buildAnalytics.AppID),
	// 	zap.String("stack_id", buildAnalytics.StackID),
	// 	zap.String("platform", buildAnalytics.Platform),
	// 	zap.String("cli_version", buildAnalytics.CLIVersion),
	// 	zap.String("status", buildAnalytics.Status),
	// 	zap.Time("start_time", buildAnalytics.StartTime),
	// 	zap.Duration("run_time", buildAnalytics.Runtime),
	// 	zap.Any("raw_json_data", buildAnalytics.RawJSONData),
	// )

	dogstatsd.Track(buildAnalytics, metrics.DogStatsDBuildCounterMetricName)
	for _, aStepAnalytic := range buildAnalytics.StepAnalytics {
		dogstatsd.Track(aStepAnalytic, metrics.DogStatsDStepCounterMetricName)
		// logger.Warn(fmt.Sprintf("Step %s finished", aStepAnalytic.StepID),
		// 	zap.String("step_id", aStepAnalytic.StepID),
		// 	zap.String("status", aStepAnalytic.Status),
		// 	zap.Time("start_time", aStepAnalytic.StartTime),
		// 	zap.Duration("run_time", aStepAnalytic.Runtime),
		// 	zap.Any("raw_json_data", buildAnalytics.RawJSONData),
		// )
	}

	return httpresponse.RespondWithSuccess(w, map[string]string{"message": "ok"})
}
