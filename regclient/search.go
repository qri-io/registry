package regclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qri-io/registry"
)

// SearchFilter stores various types of filters that may be applied
// to a search
type SearchFilter struct {
	// Type denotes the ype of search filter
	Type string
	// Relation indicates the relation between the key and value
	// supported options include ["eq"|"neq"|"gt"|"gte"|"lt"|"lte"]
	Relation string
	// Key corresponds to the name of the index mapping that we wish to
	// apply the filter to
	Key string
	// Value is the predicate of the subject-relation-predicate triple
	// eg. [key=timestamp] [gte] [value=[today]]
	Value interface{}
}

// SearchParams contains the parameters that are passed to a
// Client.Search method
type SearchParams struct {
	QueryString string
	Filters     []SearchFilter
	Limit       int
	Offset      int
}

// Search makes a registry search request
func (c Client) Search(p *SearchParams) ([]*registry.Result, error) {
	params := &registry.SearchParams{
		Q: p.QueryString,
		//Filters: p.Filters,
		Limit:  p.Limit,
		Offset: p.Offset,
	}
	results, err := c.doJSONSearchReq("GET", params)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (c Client) doJSONSearchReq(method string, s *registry.SearchParams) ([]*registry.Result, error) {
	if c.cfg.Location == "" {
		return nil, ErrNoRegistry
	}

	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/search", c.cfg.Location), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	// add response to an envelope
	env := struct {
		Data []*registry.Result
		Meta struct {
			Error  string
			Status string
			Code   int
		}
	}{}

	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error %d: %s", res.StatusCode, env.Meta.Error)
	}

	return env.Data, nil
}
