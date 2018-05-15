package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/datatogether/api/apiutil"
	"github.com/qri-io/registry"
)

// NewDatasetsHandler creates a datasets handler function that operates
// on a *registry.Datasets
func NewDatasetsHandler(datasets registry.Datasets) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			ps := []*registry.Dataset{}
			switch r.Header.Get("Content-Type") {
			case "application/json":
				if err := json.NewDecoder(r.Body).Decode(&ps); err != nil {
					apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
					return
				}
			default:
				err := fmt.Errorf("Content-Type must be application/json")
				apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
				return
			}

			for _, pro := range ps {
				datasets.Store(pro.Handle, pro)
			}
			fallthrough
		case "GET":
			ps := make([]*registry.Dataset, datasets.Len())

			i := 0
			datasets.SortedRange(func(key string, p *registry.Dataset) bool {
				ps[i] = p
				i++
				return false
			})

			apiutil.WriteResponse(w, ps)
		}
	}
}

// NewDatasetHandler creates a dataset handler func that operats on
// a *registry.Datasets
func NewDatasetHandler(datasets registry.Datasets) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &registry.Dataset{}
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
			var ok bool
			if p.Path != "" {
				datasets.Range(func(key string, dataset *registry.Dataset) bool {
					if dataset.Path == p.Path {
						*p = *dataset
						ok = true
						return true
					}
					return false
				})
			} else if p.Key() != "" {
				p, ok = datasets.Load(p.Key())
			}

			if !ok {
				apiutil.NotFoundHandler(w, r)
				return
			}
		case "PUT", "POST":
			if err := registry.RegisterDataset(datasets, p); err != nil {
				apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
				return
			}
		case "DELETE":
			if err := registry.DeregisterDataset(datasets, p); err != nil {
				apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
				return
			}
		default:
			apiutil.NotFoundHandler(w, r)
			return
		}

		apiutil.WriteResponse(w, p)
	}
}
