package handlers

import (
	"net/http"
	"time"

	"github.com/datatogether/api/apiutil"
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
	m.HandleFunc("/status", apiutil.HealthCheckHandler)

	if ps := reg.Profiles; ps != nil {
		m.HandleFunc("/profile", logReq(NewProfileHandler(ps)))
		m.HandleFunc("/profiles", pro.ProtectMethods("POST")(logReq(NewProfilesHandler(ps))))
	}
	if ds := reg.Datasets; ds != nil {
		m.HandleFunc("/dataset", logReq(NewDatasetHandler(ds)))
		m.HandleFunc("/datasets", pro.ProtectMethods("POST")(logReq(NewDatasetsHandler(ds))))
	}
	if s := reg.Search; s != nil {
		m.HandleFunc("/search", logReq(NewSearchHandler(s)))
	}
	if pinset := reg.Pinset; pinset != nil {
		m.HandleFunc("/pins", logReq(NewPinsHandler(pinset)))
	}

	return m
}

func logReq(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	}
}
