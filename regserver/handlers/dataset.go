package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/datatogether/api/apiutil"
	"github.com/qri-io/registry"
	"github.com/qri-io/registry/ns"
)

// DefaultLimit is the default limit of datasets that will get sent back on
// a dataset list request
const DefaultLimit = 25

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
			var limit, offset int
			if i, err := apiutil.ReqParamInt("limit", r); err == nil {
				limit = i
			}
			if i, err := apiutil.ReqParamInt("offset", r); err == nil {
				offset = i
			}
			if offset < 0 {
				offset = 0
			}
			if limit <= 0 {
				limit = DefaultLimit
			}

			ps := make([]*registry.Dataset, limit)

			i := 0

			datasets.SortedRange(func(key string, p *registry.Dataset) bool {
				if i < offset {
					i++
					return false
				}
				if i-offset == limit {
					return true
				}
				ps[i] = p
				i++
				return false
			})
			apiutil.WriteResponse(w, ps[:datasets.Len()-offset])
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
			if r.Method != "GET" {
				err := fmt.Errorf("Content-Type must be application/json")
				apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
				return
			}
		}

		switch r.Method {
		case "GET":
			if !strings.HasPrefix(r.URL.Path, "/dataset/") {
				err := fmt.Errorf("no reference provided")
				apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
				return
			}
			path := ns.HTTPPathToQriPath(r.URL.Path[len("/dataset/"):])
			ref, err := ns.ParseRef(path)
			if err != nil {
				apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
				return
			}

			if ref.IsEmpty() {
				err := fmt.Errorf("no reference provided")
				apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
				return
			}

			var ok bool
			datasets.Range(func(key string, dataset *registry.Dataset) bool {
				if dataset.Path == ref.Path || (ref.Name == dataset.Name && ref.Peername == dataset.Handle) {
					*p = *dataset
					ok = true
					return true
				}
				return false
			})

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
