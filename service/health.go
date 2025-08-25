package service

import (
	"net/http"

	"github.com/bitrise-io/api-utils/httpresponse"
)

func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	httpresponse.RespondWithSuccessNoErr(w, map[string]string{"status": "alive"})
}

func ReadinessHandler(w http.ResponseWriter, r *http.Request) error {
	tracker, err := GetTrackerFromContext(r.Context())
	if err != nil {
		return httpresponse.RespondWithError(w, "Tracker not available", http.StatusServiceUnavailable)
	}

	if tracker == nil {
		return httpresponse.RespondWithError(w, "PubSub tracker not initialized", http.StatusServiceUnavailable)
	}

	if err := tracker.HealthCheck(); err != nil {
		return httpresponse.RespondWithError(w, "PubSub connectivity failed: "+err.Error(), http.StatusServiceUnavailable)
	}

	// If we get here, all dependencies are ready
	return httpresponse.RespondWithSuccess(w, map[string]string{"status": "ready"})
}
