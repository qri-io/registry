package regclient

import (
	"net/http/httptest"
	"testing"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/regserver/handlers"
)

func TestSearchRequests(t *testing.T) {

	memDs := registry.NewMemDatasets()
	reg := registry.Registry{
		Profiles: registry.NewMemProfiles(),
		Datasets: memDs,
		Search:   registry.MockSearch{memDs},
	}

	srv := httptest.NewServer(handlers.NewRoutes(reg))
	c := NewClient(&Config{
		Location: srv.URL,
	})

	searchParams := &SearchParams{QueryString: "presidents", Limit: 100, Offset: 0}
	// TODO: need to add tests that actually inspect the search results
	_, err := c.Search(searchParams)
	if err != nil {
		t.Errorf("error executing search: %s", err)
	}
}
