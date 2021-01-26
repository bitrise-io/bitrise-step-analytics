package service_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bitrise-io/bitrise-step-analytics/service"
	"github.com/stretchr/testify/require"
)

func Test_ToolsPostHandler(t *testing.T) {
	handler := service.ToolsPostHandler

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
				`"build_slug":"build-slug"` +
				`,"tool_usage":[` +
				`{"name":"test-tool","version":"1.2.3","fresh_install":false}` +
				`,{"name":"test-tool-2","version":"1.2.3","fresh_install":true}` +
				`]` +
				`}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"ok"}` + "\n",
		},
		{
			testName: "when no tool usage provided ",
			requestBody: `{` +
				`"build_slug":"build-slug"` +
				`}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, please provide at least one tool_usage"}` + "\n",
		},
		{
			testName: "when no name provided for a tool",
			requestBody: `{` +
				`"build_slug":"build-slug"` +
				`,"tool_usage":[` +
				`{"version":"1.2.3","fresh_install":false}` +
				`,{"name":"test-tool-2","version":"1.2.3","fresh_install":true}` +
				`]` +
				`}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, please fill name for tool_usage"}` + "\n",
		},
		{
			testName:           "when no request body provided",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, JSON decode failed"}` + "\n",
		},
		{
			testName:           "when no tools data provided",
			requestBody:        `{}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, please provide build_slug"}` + "\n",
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			r, err := http.NewRequest("POST", "/tools", bytes.NewBuffer([]byte(tc.requestBody)))
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
