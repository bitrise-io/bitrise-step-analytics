package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

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
