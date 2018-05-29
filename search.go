package registry

// SearchParams encapsulates parameters provided to Searchable.Search
type SearchParams struct {
	Q             string
	Limit, Offset int
}

// Searchable is an interface for supporting search
type Searchable interface {
	Search(p SearchParams) ([]Result, error)
}

// TODO: Define an error 'searchNotSupported'

// Result is the interface that a search result implements
type Result struct {
	Type  string      // one of ["dataset", "profile"] for now
	ID    string      // identifier to lookup
	Value interface{} // Value returned
}
