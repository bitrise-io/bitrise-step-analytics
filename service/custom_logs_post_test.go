package service_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitrise-io/bitrise-step-analytics/service"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func Test_CustomLogsPostHandler(t *testing.T) {
	handler := service.CustomLogsPostHandler

	for _, tc := range []struct {
		testName            string
		requestBody         string
		expectedInternalErr string
		expectedBody        string
		expectedLogContent  map[string]interface{}
		expectedStatusCode  int
	}{
		{
			testName:           "ok, minimal",
			requestBody:        `{}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"ok"}` + "\n",
		},
		{
			testName:           "ok, more complex",
			expectedLogContent: map[string]interface{}{"key1": "value1"},
			requestBody:        `{"log_level":"info","message":"test message","data":{"key1":"value1"}}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"ok"}` + "\n",
		},
		{
			testName:           "when request body isn't a valid JSON, it retrieves bad reqest error",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, JSON decode failed"}` + "\n",
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			core, recorded := observer.New(zapcore.InfoLevel)
			zl := zap.New(core)

			r, err := http.NewRequest("POST", "/logs", bytes.NewBuffer([]byte(tc.requestBody)))
			require.NoError(t, err)

			r = r.WithContext(service.ContextWithLoggerProvider(r.Context(), service.NewLoggerProvider(zl)))

			rr := httptest.NewRecorder()
			internalServerError := handler(rr, r)

			for _, logs := range recorded.All() {
				for _, field := range logs.Context {
					require.Equal(t, tc.expectedLogContent, field.Interface)
				}
			}

			if tc.expectedBody != "" {
				require.Equal(t, tc.expectedBody, rr.Body.String())
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
		})
	}
}
