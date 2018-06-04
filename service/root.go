package service

import "net/http"

// RootHandler ...
func RootHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithSuccessNoErr(w, map[string]string{"message": "Bitrise Step Analytics"})
}
