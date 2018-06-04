package router

import (
	"github.com/gorilla/mux"
	"github.com/slapec93/bitrise-step-analytics/configs"
	"github.com/slapec93/bitrise-step-analytics/service"
)

// New ...
func New(config configs.ConfigModel) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", service.RootHandler)
	return r
}
