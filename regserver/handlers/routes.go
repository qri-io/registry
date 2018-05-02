package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/datatogether/api/apiutil"
	"github.com/qri-io/registry"
	"github.com/sirupsen/logrus"
)

var (
	// logger
	log      = logrus.New()
	adminKey string
)

func init() {
	adminKey = NewAdminKey()
	log.Infof("admin key: %s", adminKey)
}

// NewRoutes allocates server handlers along standard routes
func NewRoutes(ps *registry.Profiles) http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", apiutil.HealthCheckHandler)
	m.HandleFunc("/profile", logReq(NewProfileHandler(ps)))
	m.HandleFunc("/profiles", protectedRoute([]string{"POST"}, logReq(NewProfilesHandler(ps))))

	return m
}

func logReq(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s %s", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	}
}

func protectedRoute(methods []string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, m := range methods {
			if r.Method == m {
				_, password, set := r.BasicAuth()
				if !set || password != adminKey {
					log.Infof("'%s' != '%s'", adminKey, password)
					apiutil.WriteErrResponse(w, http.StatusForbidden, errors.New("invalid key"))
					return
				}
			}
		}

		h.ServeHTTP(w, r)
	}
}
