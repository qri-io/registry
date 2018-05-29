package registry

import (
	"fmt"
)

// SearchParams encapsulates parameters provided to Searchable.Search
type SearchParams struct {
	Q             string
	Limit, Offset int
}

// Searchable is an interface for supporting search
type Searchable interface {
	Search(p SearchParams) ([]Result, error)
}

// ErrSearchNotSupported is returned for collections that are not searchable
var ErrSearchNotSupported = fmt.Errorf("search not supported")

// NilSearch is a stub for collections that don't support search
type NilSearch bool

// Search returns an error indicating that search is not supported
func (ns NilSearch) Search(p SearchParams) ([]Result, error) {
	return nil, ErrSearchNotSupported
}

// Result is the interface that a search result implements
type Result struct {
	Type  string      // one of ["dataset", "profile"] for now
	ID    string      // identifier to lookup
	Value interface{} // Value returned
}
