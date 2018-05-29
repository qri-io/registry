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
	profiles = registry.NewMemProfiles()
	// in-memory datasets for now
	datasets = registry.NewMemDatasets()
	// logger
	log = logrus.New()

	adminKey string

	nilSearch registry.NilSearch
)

func init() {
	adminKey = handlers.NewAdminKey()
	log.Infof("admin key: %s", adminKey)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	pro := handlers.NewBAProtector("username", adminKey)

	s := http.Server{
		Addr:    ":" + port,
		Handler: handlers.NewRoutes(pro, profiles, datasets, nilSearch),
	}
	log.Infof("serving on: %s", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Info(err.Error())
	}
}
