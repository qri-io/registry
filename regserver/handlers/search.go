package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/datatogether/api/apiutil"
	"github.com/qri-io/registry"
)

// NewSearchHandler creates a search handler function taht operates on a *registry.Searchable
func NewSearchHandler(s registry.Searchable) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &registry.SearchParams{}
		switch r.Header.Get("Content-Type") {
		case "application/json":
			if err := json.NewDecoder(r.Body).Decode(p); err != nil {
				apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
				return
			}
		default:
			err := fmt.Errorf("Content-Type must be application/json")
			apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
			return
		}
		switch r.Method {
		case "GET":
			if p.Q != "" {
				results, err := s.Search(*p)
				if err != nil {
					apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
					return
				}
				apiutil.WriteResponse(w, results)
				return
			}
		}
	}
}
