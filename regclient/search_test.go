package regclient

import (
	"net/http/httptest"
	"testing"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/regserver/handlers"
)

func TestSearchRequests(t *testing.T) {

	searchParams := &SearchParams{QueryString: "presidents", Limit: 100, Offset: 0}
	memDs := registry.NewMemDatasets()
	srv := httptest.NewServer(handlers.NewRoutes(handlers.NewNoopProtector(), registry.NewMemProfiles(), memDs, &registry.MockSearch{memDs}))
	c := NewClient(&Config{
		Location: srv.URL,
	})
	// TODO: need to add tests that actually inspect the search results
	_, err := c.Search(searchParams)
	if err != nil {
		t.Errorf("error executing search: %s", err)
	}
}
