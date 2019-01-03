package service_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitrise-team/bitrise-step-analytics/service"
	"github.com/stretchr/testify/require"
)

func Test_MetricsPostHandler(t *testing.T) {
	handler := service.MetricsPostHandler

	for _, tc := range []struct {
		testName            string
		requestBody         string
		expectedInternalErr string
		expectedBody        string
		expectedLogContent  map[string]interface{}
		expectedStatusCode  int
	}{
		{
			testName: "ok, minimal",
			requestBody: `{` +
				`"app_id":"app-slug","stack_id":"standard1","platform":"ios","cli_version":"1.21","status":"success","start_time":"2019-01-03T18:11:53.171409Z","run_time":121` +
				`,"step_analytics":[` +
				`{"step_id":"deploy_to_bitrise_io","status":"0","start_time":"2019-01-03T18:11:53.171409Z","run_time":120}` +
				`,{"step_id":"script","status":"0","start_time":"2019-01-03T18:11:53.171409Z","run_time":210}` +
				`]` +
				`}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"ok"}` + "\n",
		},
		{
			testName:           "when no request body provided",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, JSON decode failed"}` + "\n",
		},
		{
			testName:           "when no metrics data provided",
			requestBody:        `{}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, please provide metrics data"}` + "\n",
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			r, err := http.NewRequest("POST", "/metrics", bytes.NewBuffer([]byte(tc.requestBody)))
			require.NoError(t, err)

			r = r.WithContext(service.ContextWithDogStatsDMetrics(r.Context(), &testDogStatsdMetrics{}))

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
