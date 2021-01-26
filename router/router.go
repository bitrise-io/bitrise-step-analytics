package router

import (
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-io/bitrise-step-analytics/configs"
	"github.com/bitrise-io/bitrise-step-analytics/metrics"
	"github.com/bitrise-io/bitrise-step-analytics/service"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

// New ...
func New(config configs.ConfigModel) *mux.Router {
	r := mux.NewRouter(mux.WithServiceName("step-analytics-mux")).StrictSlash(true)

	middlewareProvider := service.MiddlewareProvider{
		Client: metrics.NewClient(config.SegmentWriteKey),
	}

	r.Handle("/", middlewareProvider.CommonMiddleware().ThenFunc(service.RootHandler))
	r.Handle("/metrics", middlewareProvider.MiddlewareWithClient().Then(
		httpresponse.InternalErrHandlerFuncAdapter(service.MetricsPostHandler))).Methods("POST")
	r.Handle("/logs", middlewareProvider.MiddlewareWithClient().Then(
		httpresponse.InternalErrHandlerFuncAdapter(service.CustomLogsPostHandler))).Methods("POST")
	r.Handle("/tools", middlewareProvider.MiddlewareWithClient().Then(
		httpresponse.InternalErrHandlerFuncAdapter(service.ToolsPostHandler))).Methods("POST")

	return r
}
