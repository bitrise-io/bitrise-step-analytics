package main

import (
	"log"
	"net/http"

	"github.com/bitrise-io/bitrise-step-analytics/configs"
	"github.com/bitrise-io/bitrise-step-analytics/router"
	"github.com/pkg/errors"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func initilize() error {
	tracer.Start(tracer.WithServiceName("step-analytics"))
	defer tracer.Stop()

	conf, err := configs.CreateAndValidate()
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "Failed to read Configs")
	}

	// Routing
	http.Handle("/", router.New(conf))

	log.Println("Starting - using port:", conf.Port)
	if err := http.ListenAndServe(":"+conf.Port, nil); err != nil {
		return errors.Wrap(errors.WithStack(err), "Failed to ListenAndServe")
	}
	return nil
}

func main() {
	err := initilize()
	if err != nil {
		log.Fatalf(" [!] Exception: Failed to initialize Bitrise Step Analytics: %+v", err)
	}
}
