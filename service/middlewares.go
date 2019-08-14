package service

import (
	"net/http"

	"github.com/bitrise-io/bitrise-step-analytics/metrics"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

// MiddlewareProvider ...
type MiddlewareProvider struct {
	LoggerProvider   LoggerInterface
	Client metrics.Interface
}

func createSetLoggerProviderMiddleware(loggerProvider LoggerInterface) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ContextWithLoggerProvider(r.Context(), loggerProvider)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func createSetClientMiddleware(Client metrics.Interface) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := ContextWithClient(r.Context(), Client)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CommonMiddleware ...
func (m MiddlewareProvider) CommonMiddleware() alice.Chain {
	return alice.New(
		cors.AllowAll().Handler,
	)
}

// MiddlewareWithLoggerProvider ...
func (m MiddlewareProvider) MiddlewareWithLoggerProvider() alice.Chain {
	return m.CommonMiddleware().Append(
		createSetLoggerProviderMiddleware(m.LoggerProvider),
	)
}

// MiddlewareWithClient ...
func (m MiddlewareProvider) MiddlewareWithClient() alice.Chain {
	return m.CommonMiddleware().Append(
		createSetClientMiddleware(m.Client),
	)
}
