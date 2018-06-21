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
	// logger
	log = logrus.New()

	adminKey string
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
	ps := registry.NewMemProfiles()
	reg := registry.Registry{
		Profiles: ps,
		Datasets: registry.NewMemDatasets(),
		Pinset:   &registry.MemPinset{Profiles: ps},
	}

	s := http.Server{
		Addr:    ":" + port,
		Handler: handlers.NewRoutes(pro, reg),
	}

	log.Infof("serving on: %s", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Info(err.Error())
	}
}
