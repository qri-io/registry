package regclient

import (
	"net/http/httptest"
	"testing"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/regserver/handlers"
)

func TestSearchRequests(t *testing.T) {

	searchParams := &SearchParams{QueryString: "presidents", Limit: 100, Offset: 0}

	srv := httptest.NewServer(handlers.NewRoutes(handlers.NewNoopProtector(), registry.Registry{Profiles: registry.NewMemProfiles(), Datasets: registry.NewMemDatasets(), Search: nilSearch}))
	c := NewClient(&Config{
		Location: srv.URL,
	})
	// TODO: need to add tests that actually inspect the search results
	_, err := c.Search(searchParams)
	if err != nil {
		t.Errorf("error executing search: %s", err)
	}
}
