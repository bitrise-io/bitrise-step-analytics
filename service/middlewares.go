package service

import (
	"net/http"

	"github.com/bitrise-io/bitrise-step-analytics/event"
	"github.com/bitrise-io/bitrise-step-analytics/metrics"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

type MiddlewareProvider struct {
	Client  metrics.Interface
	Tracker event.Tracker
}

func (m MiddlewareProvider) MiddlewareWithClient() alice.Chain {
	return m.CommonMiddleware().Append(
		createSetClientMiddleware(m.Client),
	)
}

func createSetClientMiddleware(Client metrics.Interface) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ContextWithClient(r.Context(), Client)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m MiddlewareProvider) MiddlewareWithTracker() alice.Chain {
	return m.CommonMiddleware().Append(
		createSetTrackerMiddleware(m.Tracker),
	)
}

func createSetTrackerMiddleware(tracker event.Tracker) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ContextWithTracker(r.Context(), tracker)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (m MiddlewareProvider) CommonMiddleware() alice.Chain {
	return alice.New(
		cors.AllowAll().Handler,
	)
}
