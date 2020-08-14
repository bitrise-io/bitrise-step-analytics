package service

import (
	"net/http"

	"github.com/bitrise-io/bitrise-step-analytics/metrics"
	"github.com/justinas/alice"
	"github.com/rs/cors"
)

// MiddlewareProvider ...
type MiddlewareProvider struct {
	Client metrics.Interface
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

// MiddlewareWithClient ...
func (m MiddlewareProvider) MiddlewareWithClient() alice.Chain {
	return m.CommonMiddleware().Append(
		createSetClientMiddleware(m.Client),
	)
}
