package service

import (
	"net/http"

	"github.com/bitrise-io/api-utils/httpresponse"
)

// RootHandler ...
func RootHandler(w http.ResponseWriter, r *http.Request) {
	httpresponse.RespondWithSuccessNoErr(w, map[string]string{"message": "Bitrise Step Analytics"})
}
