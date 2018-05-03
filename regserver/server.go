// Package regserver is a wrapper around the handlers package,
// turning it into a proper http server
package main

import (
	"net/http"
	"os"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/regserver/handlers"
	"github.com/sirupsen/logrus"
)

var (
	// in-memory profiles for now
	profiles = registry.NewProfiles()
	// in-memory datasets for now
	datasets = registry.NewDatasets()
	// logger
	log = logrus.New()
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	s := http.Server{
		Addr:    ":" + port,
		Handler: handlers.NewRoutes(profiles, datasets),
	}
	log.Infof("serving on: %s", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Info(err.Error())
	}
}
