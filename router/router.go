package router

import (
	"github.com/bitrise-io/api-utils/httpresponse"
	"github.com/bitrise-team/bitrise-step-analytics/configs"
	"github.com/bitrise-team/bitrise-step-analytics/service"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/gorilla/mux"
)

// New ...
func New(config configs.ConfigModel) *mux.Router {
	r := mux.NewRouter(mux.WithServiceName("steps-mux")).StrictSlash(true)
	commonMiddleware := alice.New(
		cors.AllowAll().Handler,
	)

	r.Handle("/", commonMiddleware.ThenFunc(service.RootHandler))
	r.Handle("/log-analytics", commonMiddleware.Then(
		httpresponse.InternalErrHandlerFuncAdapter(service.AnalyticsLogHandler))).Methods("POST")

	return r
}
