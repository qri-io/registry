package handlers

import (
	"net/http"
	"time"

	"github.com/qri-io/registry"
	"github.com/sirupsen/logrus"
)

var (
	// logger
	log = logrus.New()
)

// NewRoutes allocates server handlers along standard routes
func NewRoutes(pro MethodProtector, reg registry.Registry) http.Handler {
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
	if pinset := reg.Pinset; pinset != nil {
		m.HandleFunc("/pins", logReq(NewPinsHandler(pinset)))
		m.HandleFunc("/pins/status", logReq(NewPinStatusHandler(pinset)))
	}
	if rs := reg.Reputations; rs != nil {
		m.HandleFunc("/reputation", (logReq(NewReputationHandler(rs))))
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
