package router

import (
	"fmt"

	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-team/bitrise-step-analytics/configs"
	"github.com/bitrise-team/bitrise-step-analytics/metrics"
	"github.com/bitrise-team/bitrise-step-analytics/service"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

// New ...
func New(config configs.ConfigModel) *mux.Router {
	r := mux.NewRouter(mux.WithServiceName("steps-mux")).StrictSlash(true)
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("Failed to initialize zap logger")
	}

	middlewareProvider := service.MiddlewareProvider{
		LoggerProvider:   service.NewLoggerProvider(logger),
		DogStatsDMetrics: metrics.NewDogStatsDMetrics(""),
	}

	r.Handle("/", middlewareProvider.CommonMiddleware().ThenFunc(service.RootHandler))
	r.Handle("/metrics", middlewareProvider.CommonMiddleware().Then(
		httpresponse.InternalErrHandlerFuncAdapter(service.MetricsPostHandler))).Methods("POST")
	r.Handle("/logs", middlewareProvider.MiddlewareWithLoggerProvider().Then(
		httpresponse.InternalErrHandlerFuncAdapter(service.CustomLogsPostHandler))).Methods("POST")

	return r
}
