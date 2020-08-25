package service_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitrise-io/bitrise-step-analytics/service"
	"github.com/stretchr/testify/require"
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
			requestBody:        `{"log_level":"warn","message":"test","data":{"step_id":"test-step-id","tag":"test-tag"}}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"ok"}` + "\n",
		},
		{
			testName:           "error, when request body isn't a valid JSON, it retrieves bad reqest error",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, JSON decode failed"}` + "\n",
		},
		{
			testName:           "error, when all required param missing",
			requestBody:        `{}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, please provide log_level"}` + "\n",
		},
		{
			testName:           "error, missing message",
			requestBody:        `{"log_level":"warn"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, please provide message"}` + "\n",
		},
		{
			testName:           "error, when message also set",
			requestBody:        `{"log_level":"warn","message":"test"}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, please provide data.step_id"}` + "\n",
		},
		{
			testName:           "error, when step id also set",
			requestBody:        `{"log_level":"warn","message":"test","data":{"step_id":"test-step-id"}}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, please provide data.tag"}` + "\n",
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			r, err := http.NewRequest("POST", "/logs", bytes.NewBuffer([]byte(tc.requestBody)))
			require.NoError(t, err)

			r = r.WithContext(service.ContextWithClient(r.Context(), &testClient{}))

			rr := httptest.NewRecorder()
			internalServerError := handler(rr, r)

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

func Test_CustomLogsPostHandlerInvalidContext(t *testing.T) {
	handler := service.CustomLogsPostHandler

	r, err := http.NewRequest("POST", "/logs", bytes.NewBuffer([]byte(`{"log_level":"warn","message":"test","data":{"step_id":"test-step-id","tag":"test-tag"}}`)))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	internalServerError := handler(rr, r)
	expectedInternalErr := "DogStatsD not found in Context"

	require.EqualError(t, internalServerError, expectedInternalErr,
		"Expected internal err: %s", expectedInternalErr)
}
