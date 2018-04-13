package handlers

import (
	"net/http"
	"time"

	"github.com/datatogether/api/apiutil"
	"github.com/qri-io/registry"
	"github.com/sirupsen/logrus"
)

// logger
var log = logrus.New()

// NewRoutes allocates server handlers along standard routes
func NewRoutes(ps *registry.Profiles) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", apiutil.HealthCheckHandler)
	m.HandleFunc("/profile", logReq(NewProfileHandler(ps)))
	m.HandleFunc("/profiles", logReq(NewProfilesHandler(ps)))

	return m
}

func logReq(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	}
}
