package registry

import (
	"fmt"
	"strings"
)

// Searchable is an opt-in interface for registries that wish to
// support search
type Searchable interface {
	Search(p SearchParams) ([]Result, error)
}

// SearchParams encapsulates parameters provided to Searchable.Search
type SearchParams struct {
	Q             string
	Limit, Offset int
}

// Result is the interface that a search result implements
type Result struct {
	Type  string      // one of ["dataset", "profile"] for now
	ID    string      // identifier to lookup
	Value interface{} // Value returned
}

// ErrSearchNotSupported is the canonical error to indicate search
// isn't implemented
var ErrSearchNotSupported = fmt.Errorf("search not supported")

// NilSearch is a basic implementation of Searchable which returns
// an error to indicate that search is not supported
type NilSearch bool

// Search returns an error indicating that search is not supported
func (ns NilSearch) Search(p SearchParams) ([]Result, error) {
	return nil, ErrSearchNotSupported
}

// MockSearch is a very naive implementation of search that wraps
// registry.Datasets and only checks for exact substring matches of
// dataset's meta.title property. It's mainly intended for testing
// purposes.
type MockSearch struct {
	Datasets Datasets
}

// Search is a trivial search implementation used for testing
func (ms MockSearch) Search(p SearchParams) (results []Result, err error) {
	ms.Datasets.Range(func(key string, ds *Dataset) bool {
		dsname := ""
		if ds.Meta != nil {
			dsname = strings.ToLower(ds.Meta.Title)
		}
		if strings.Contains(dsname, strings.ToLower(p.Q)) {
			result := &Result{Value: ds}
			results = append(results, *result)
		}
		return false
	})
	return results, nil
}
