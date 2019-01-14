package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitrise-io/api-utils/httpresponse"
	"go.uber.org/zap"
)

const (
	logLvlInfo  string = "info"
	logLvlWarn  string = "warn"
	logLvlError string = "error"
)

// LogPostRequestBody ...
type LogPostRequestBody struct {
	LogLevel string                 `json:"log_level"`
	Message  string                 `json:"message"`
	Data     map[string]interface{} `json:"data"`
}

// CustomLogsPostHandler ...
func CustomLogsPostHandler(w http.ResponseWriter, r *http.Request) error {
	logData := LogPostRequestBody{}
	defer httpresponse.RequestBodyCloseWithErrorLog(r)
	if err := json.NewDecoder(r.Body).Decode(&logData); err != nil {
		return httpresponse.RespondWithBadRequestError(w, "Invalid request body, JSON decode failed")
	}

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %s", err)
	}
	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Printf("Failed to sync logger: %s", err)
		}
	}()

	switch logData.LogLevel {
	case logLvlInfo:
		logger.Info(logData.Message, zap.Any("raw_json_data", logData.Data))
	case logLvlWarn:
		logger.Warn(logData.Message, zap.Any("raw_json_data", logData.Data))
	case logLvlError:
		logger.Error(logData.Message, zap.Any("raw_json_data", logData.Data))
	}

	return httpresponse.RespondWithSuccess(w, map[string]string{"message": "ok"})
}
