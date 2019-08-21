// Package handlers creates HTTP handler functions for registry interface implementations
package handlers

import (
	"net/http"
	"time"

	"github.com/qri-io/dag/dsync"
	"github.com/qri-io/registry"
	"github.com/qri-io/registry/pinset"
	"github.com/sirupsen/logrus"
)

var (
	// logger
	log = logrus.New()
)

// SetLogLevel controls how detailed handler logging is
func SetLogLevel(level string) error {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}
	log.SetLevel(lvl)
	return nil
}

// RouteOptions defines configuration details for NewRoutes
type RouteOptions struct {
	Protector MethodProtector
	Pinset    pinset.Pinset
	Dsync     *dsync.Dsync
}

// AddPinset creates a configuration func for passing to NewRoutes
func AddPinset(ps pinset.Pinset) func(o *RouteOptions) {
	return func(o *RouteOptions) {
		o.Pinset = ps
	}
}

// AddDsync creates a configuration func for passing to NewRoutes
func AddDsync(rec *dsync.Dsync) func(o *RouteOptions) {
	return func(o *RouteOptions) {
		o.Dsync = rec
	}
}

// AddProtector creates a configuration func for passing to NewRoutes
func AddProtector(p MethodProtector) func(o *RouteOptions) {
	return func(o *RouteOptions) {
		o.Protector = p
	}
}

// NewRoutes allocates server handlers along standard routes
func NewRoutes(reg registry.Registry, opts ...func(o *RouteOptions)) *http.ServeMux {
	o := &RouteOptions{
		Protector: NoopProtector(0),
	}
	for _, opt := range opts {
		opt(o)
	}

	pro := o.Protector
	m := http.NewServeMux()
	m.HandleFunc("/", HealthCheckHandler)

	if ps := reg.Profiles; ps != nil {
		m.HandleFunc("/profile", logReq(NewProfileHandler(ps)))
		m.HandleFunc("/profiles", pro.ProtectMethods("POST")(logReq(NewProfilesHandler(ps))))
	}

	if ds := reg.Datasets; ds != nil {
		m.HandleFunc("/dataset", logReq(NewDatasetHandler(ds, reg.Indexer)))
		m.HandleFunc("/dataset/", logReq(NewDatasetHandler(ds, reg.Indexer)))
		m.HandleFunc("/datasets", pro.ProtectMethods("POST")(logReq(NewDatasetsHandler(ds, reg.Indexer))))
	}

	if s := reg.Search; s != nil {
		m.HandleFunc("/search", logReq(NewSearchHandler(s)))
	}
	if rs := reg.Reputations; rs != nil {
		m.HandleFunc("/reputation", (logReq(NewReputationHandler(rs))))
	}

	if o.Pinset != nil {
		m.HandleFunc("/pins", logReq(NewPinsHandler(o.Pinset)))
		m.HandleFunc("/pins/status", logReq(NewPinStatusHandler(o.Pinset)))
	}
	if o.Dsync != nil {
		m.HandleFunc("/dsync", logReq(dsync.HTTPRemoteHandler(o.Dsync)))
	}

	return m
}

func logReq(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	}
}

// HealthCheckHandler is a basic "hey I'm fine" for load balancers & co
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"meta":{"code": 200,"status":"ok"},"data":null}`))
}
