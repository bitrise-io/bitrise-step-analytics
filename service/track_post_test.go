package service_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/bitrise-io/bitrise-step-analytics/event/mocks"
	"github.com/bitrise-io/bitrise-step-analytics/service"
	"github.com/stretchr/testify/require"
)

func Test_TrackPostHandler(t *testing.T) {
	handler := service.TrackPostHandler
	mockTracker := new(mocks.Tracker)
	mockTracker.On("Send", mock.Anything).Return(nil).Once()
	mockTracker.On("Send", mock.Anything).Return(errors.New("test error")).Once()

	for _, tc := range []struct {
		testName            string
		requestBody         string
		expectedInternalErr string
		expectedBody        string
		expectedStatusCode  int
	}{
		{
			testName: "ok, minimal",
			requestBody: `{
"id": "4d94adc2-33a7-4ad9-9f69-5c937a6da52a",
"event_name": "test_event"
}`,
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"message":"ok"}` + "\n",
		},
		{
			testName: "when no event provided",
			requestBody: `{
"id": "4d94adc2-33a7-4ad9-9f69-5c937a6da52a",
"test_property": "test_value"
}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, please provide event's name"}` + "\n",
		},
		{
			testName: "when invalid id provided",
			requestBody: `{
"id": "fake-id"
}`,
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, please provide event's name"}` + "\n",
		},
		{
			testName:           "when no request body provided",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, JSON decode failed: EOF"}` + "\n",
		},
		{
			testName:           "when invalid request body provided",
			requestBody:        "invalidJSON",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"message":"Invalid request body, JSON decode failed: invalid character 'i' looking for beginning of value"}` + "\n",
		},
		{
			testName: "when tracking has failed",
			requestBody: `{
"id": "4d94adc2-33a7-4ad9-9f69-5c937a6da52a",
"event_name": "test_event"
}`,
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody:       `{"message":"Couldn't send analytics event: test error"}` + "\n",
		},
	} {
		t.Run(tc.testName, func(t *testing.T) {
			r, err := http.NewRequest("POST", "/track", bytes.NewBuffer([]byte(tc.requestBody)))
			require.NoError(t, err)

			r = r.WithContext(service.ContextWithTracker(r.Context(), mockTracker))

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
