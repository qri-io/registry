package registry

import (
	"fmt"
)

// Searchable is an interface for supporting search
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
