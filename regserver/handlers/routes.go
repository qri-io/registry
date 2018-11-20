package handlers

import (
	"net/http"
	"time"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/pinset"
	"github.com/sirupsen/logrus"
)

var (
	// logger
	log = logrus.New()
)

// NewRoutes allocates server handlers along standard routes
func NewRoutes(pro MethodProtector, reg registry.Registry) *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/", HealthCheckHandler)

	if ps := reg.Profiles; ps != nil {
		m.HandleFunc("/profile", logReq(NewProfileHandler(ps)))
		m.HandleFunc("/profiles", pro.ProtectMethods("POST")(logReq(NewProfilesHandler(ps))))
	}
	if ds := reg.Datasets; ds != nil {
		m.HandleFunc("/dataset", logReq(NewDatasetHandler(ds)))
		m.HandleFunc("/dataset/", logReq(NewDatasetHandler(ds)))
		m.HandleFunc("/datasets", pro.ProtectMethods("POST")(logReq(NewDatasetsHandler(ds))))
	}
	if s := reg.Search; s != nil {
		m.HandleFunc("/search", logReq(NewSearchHandler(s)))
	}
	if rs := reg.Reputations; rs != nil {
		m.HandleFunc("/reputation", (logReq(NewReputationHandler(rs))))
	}
	return m
}

// NewRoutesPinset adds standard routes and pinset routes
func NewRoutesPinset(pro MethodProtector, reg registry.Registry, ps pinset.Pinset) *http.ServeMux {
	m := NewRoutes(pro, reg)
	if ps != nil {
		m.HandleFunc("/pins", logReq(NewPinsHandler(ps)))
		m.HandleFunc("/pins/status", logReq(NewPinStatusHandler(ps)))
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
