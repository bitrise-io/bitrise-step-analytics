package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// StandardErrorResponse ...
type StandardErrorResponse struct {
	Message string `json:"message"`
}

// RespondWithJSON ...
func RespondWithJSON(w http.ResponseWriter, httpCode int, respModel interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(&respModel); err != nil {
		return errors.Wrapf(err, "Failed to respond (encode) with JSON for response model: %#v", respModel)
	}
	return nil
}

// RespondWithJSONNoErr ...
func RespondWithJSONNoErr(w http.ResponseWriter, httpCode int, respModel interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(httpCode)
	if err := json.NewEncoder(w).Encode(&respModel); err != nil {
		log.Printf(" [!] Exception: failed to respond with JSON, error: %+v", errors.WithStack(err))
	}
}

// RespondWithSuccess ...
func RespondWithSuccess(w http.ResponseWriter, respModel interface{}) error {
	return RespondWithJSON(w, http.StatusOK, respModel)
}

// RespondWithSuccessNoErr ...
func RespondWithSuccessNoErr(w http.ResponseWriter, respModel interface{}) {
	RespondWithJSONNoErr(w, http.StatusOK, respModel)
}

// RespondWithError ...
func RespondWithError(w http.ResponseWriter, errMsg string, httpErrCode int) error {
	return RespondWithJSON(w, httpErrCode, StandardErrorResponse{
		Message: errMsg,
	})
}

// RespondWithBadRequest ...
func RespondWithBadRequest(w http.ResponseWriter, errMsg string) error {
	return RespondWithError(w, errMsg, http.StatusBadRequest)
}

// RespondWithInternalServerError ...
func RespondWithInternalServerError(w http.ResponseWriter, errorToLog error) {
	log.Printf(" [!] Exception: Internal Server Error: %+v", errors.WithStack(errorToLog))

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_, err := fmt.Fprintln(w, `{"message":"Internal Server Error"}`)
	if err != nil {
		log.Printf(" [!] Exception: failed to write Internal Server Error response, error: %+v", errors.WithStack(err))
	}
}

// HanderFuncWithInternalError ...
type HanderFuncWithInternalError func(http.ResponseWriter, *http.Request) error

// InternalErrHandlerFuncAdapter ...
func InternalErrHandlerFuncAdapter(h HanderFuncWithInternalError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		intServErr := h(w, r)
		if intServErr != nil {
			RespondWithInternalServerError(w, errors.WithStack(intServErr))
		}
	})
}
