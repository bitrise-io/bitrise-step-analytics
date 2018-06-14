package router

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/slapec93/bitrise-step-analytics/configs"
	"github.com/slapec93/bitrise-step-analytics/service"
)

// New ...
func New(config configs.ConfigModel) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	commonMiddleware := alice.New(
		cors.AllowAll().Handler,
	)

	r.Handle("/", commonMiddleware.ThenFunc(service.RootHandler))
	r.Handle("/analytics", commonMiddleware.Then(
		service.InternalErrHandlerFuncAdapter(service.AnalyticsListHandler))).Methods("GET")
	r.Handle("/log-analytics", commonMiddleware.Then(
		service.InternalErrHandlerFuncAdapter(service.AnalyticsLogHandler))).Methods("POST")

	return r
}
