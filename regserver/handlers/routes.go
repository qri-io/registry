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
func NewRoutes(pro MethodProtector, ps registry.Profiles, ds registry.Datasets, searchable registry.Searchable) http.Handler {

	m := http.NewServeMux()
	m.HandleFunc("/", apiutil.HealthCheckHandler)
	m.HandleFunc("/profile", logReq(NewProfileHandler(ps)))
	m.HandleFunc("/profiles", pro.ProtectMethods("POST")(logReq(NewProfilesHandler(ps))))

	m.HandleFunc("/dataset", logReq(NewDatasetHandler(ds)))
	m.HandleFunc("/datasets", pro.ProtectMethods("POST")(logReq(NewDatasetsHandler(ds))))
	m.HandleFunc("/search", logReq(NewSearchHandler(searchable)))

	return m
}

func logReq(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	}
}
