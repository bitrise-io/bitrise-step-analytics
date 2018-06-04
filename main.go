package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/slapec93/bitrise-step-analytics/configs"
	"github.com/slapec93/bitrise-step-analytics/database"
	"github.com/slapec93/bitrise-step-analytics/router"
)

func initilize() error {
	conf, err := configs.CreateAndValidate()
	if err != nil {
		return errors.Wrap(errors.WithStack(err), "Failed to read Configs")
	}

	// Database
	if err := database.InitAndOpenDatabase(); err != nil {
		return errors.Wrap(errors.WithStack(err), "Failed to init and open the database")
	}
	defer database.Close()
	log.Println(" [OK] Database connection established")

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
