package service

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitrise-io/bitrise-step-analytics/event/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLivenessHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/liveness", nil)
	w := httptest.NewRecorder()

	LivenessHandler(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "alive")
}

func TestReadinessHandler_Success(t *testing.T) {
	mockTracker := &mocks.Tracker{}
	mockTracker.On("HealthCheck").Return(nil)
	ctx := ContextWithTracker(context.Background(), mockTracker)
	req := httptest.NewRequest("GET", "/readiness", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	err := ReadinessHandler(w, req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ready")
	mockTracker.AssertExpectations(t)
}

func TestReadinessHandler_PubSubConnectivityFailed(t *testing.T) {
	mockTracker := &mocks.Tracker{}
	mockTracker.On("HealthCheck").Return(assert.AnError)
	ctx := ContextWithTracker(context.Background(), mockTracker)
	req := httptest.NewRequest("GET", "/readiness", nil).WithContext(ctx)
	w := httptest.NewRecorder()

	err := ReadinessHandler(w, req)

	assert.NoError(t, err) // httpresponse.RespondWithError returns nil
	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
	assert.Contains(t, w.Body.String(), "PubSub connectivity failed")
	mockTracker.AssertExpectations(t)
}
