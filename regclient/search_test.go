package regclient

import (
	"net/http/httptest"
	"testing"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/regserver/handlers"
)

func TestSearchRequests(t *testing.T) {

	searchParams := &SearchParams{QueryString: "presidents", Limit: 100, Offset: 0}

	srv := httptest.NewServer(handlers.NewRoutes(handlers.NewNoopProtector(), registry.NewMemProfiles(), registry.NewMemDatasets(), nilSearch))
	c := NewClient(&Config{
		Location: srv.URL,
	})

	expectedErr := "error 400: search not supported"
	_, err := c.Search(searchParams)
	if err == nil {
		t.Errorf("expected nilSearch to return an error")
	} else if err.Error() != expectedErr {
		t.Errorf("error mismatch. expected: %s, got: %s", expectedErr, err.Error())
	}
}
