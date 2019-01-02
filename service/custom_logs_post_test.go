package service_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitrise-team/bitrise-step-analytics/service"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func Test_CustomLogsPostHandler(t *testing.T) {
	handler := service.CustomLogsPostHandler
	core, recorded := observer.New(zapcore.InfoLevel)
	zl := zap.New(core)

	for _, tc := range []struct {
		requestBody         string
		expectedInternalErr string
		expectedBody        string
		expectedStatusCode  int
	}{
		{},
	} {
		r, err := http.NewRequest("POST", "/logs", bytes.NewBuffer([]byte(tc.requestBody)))
		require.NoError(t, err)

		r.WithContext(service.ContextWithLoggerProvider(r.Context(), service.NewLoggerProvider(zl)))

		rr := httptest.NewRecorder()
		internalServerError := handler(rr, r)

		for _, logs := range recorded.All() {
			fmt.Println(logs.Message)
		}

		if tc.expectedInternalErr != "" {
			require.EqualError(t, internalServerError, tc.expectedInternalErr,
				"Expected internal err: %s | Request Body: %s | Response Code: %d, Expected Response Body: %s | Got Body: %s", tc.expectedInternalErr, tc.requestBody, rr.Code, tc.expectedBody, rr.Body.String())
		} else {
			require.NoError(t, internalServerError)
			if tc.expectedStatusCode != 0 {
				require.Equal(t, tc.expectedStatusCode, rr.Code,
					"Expected body: %s | Got body: %s", tc.expectedBody, rr.Body.String())
			}
		}
	}
}
