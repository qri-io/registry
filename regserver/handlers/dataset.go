package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/qri-io/apiutil"
	"github.com/qri-io/registry"
	"github.com/qri-io/registry/ns"
)

// DefaultLimit is the default limit of datasets that will get sent back on
// a dataset list request
const DefaultLimit = 25

// NewDatasetsHandler creates a datasets handler function that operates
// on a *registry.Datasets
func NewDatasetsHandler(datasets registry.Datasets, idxr registry.Indexer) http.HandlerFunc {
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

			// add datasets to search index if it's defined
			if idxr != nil {
				if err := idxr.IndexDatasets(ps); err != nil {
					apiutil.WriteErrResponse(w, http.StatusInternalServerError, err)
					return
				}
			}

			fallthrough
		case "GET":
			p := apiutil.PageFromRequest(r)
			offset := p.Offset()
			limit := p.Limit()
			ds := make([]*registry.Dataset, 0, limit)

			datasets.SortedRange(func(key string, d *registry.Dataset) bool {
				if offset > 0 {
					offset--
					return false
				}
				if len(ds) == limit {
					return true
				}
				ds = append(ds, d)
				return false
			})

			apiutil.WriteResponse(w, ds)
		}
	}
}

// NewDatasetHandler creates a dataset handler func that operats on
// a *registry.Datasets
func NewDatasetHandler(datasets registry.Datasets, idxr registry.Indexer) http.HandlerFunc {
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
			datasets.Range(func(key string, ds *registry.Dataset) bool {
				if ds.Path == ref.Path || (ref.Name == ds.Name && ref.Peername == ds.Handle) {
					*p = *ds
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
			if idxr != nil {
				if err := idxr.IndexDatasets([]*registry.Dataset{p}); err != nil {
					apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
					return
				}
			}
		case "DELETE":
			if err := registry.DeregisterDataset(datasets, p); err != nil {
				apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
				return
			}
			if idxr != nil {
				if err := idxr.UnindexDatasets([]*registry.Dataset{p}); err != nil {
					apiutil.WriteErrResponse(w, http.StatusBadRequest, err)
					return
				}
			}
		default:
			apiutil.NotFoundHandler(w, r)
			return
		}

		apiutil.WriteResponse(w, p)
	}
}
