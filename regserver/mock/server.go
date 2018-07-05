// Package mock provides a mock registry server for testing purposes
// it mocks the behaviour of a registry server with in-memory storage
package mock

import (
	"net/http/httptest"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/regclient"
	"github.com/qri-io/registry/regserver/handlers"
)

// NewMockServer creates an in-memory mock server without any access protection and
// a registry client to match
func NewMockServer() (*regclient.Client, *httptest.Server) {
	ds := registry.NewMemDatasets()
	reg := registry.Registry{
		Profiles: registry.NewMemProfiles(),
		Datasets: ds,
		Search:   registry.MockSearch{ds},
	}
	s := httptest.NewServer(handlers.NewRoutes(handlers.NewNoopProtector(), reg))
	c := regclient.NewClient(&regclient.Config{Location: s.URL})
	return c, s
}

// NewMockServerWithMemPinset creates an in-memory mock server without any access protection and
// a registry client to match, but also adds an in-memory pinset to test the /pin endpoint
func NewMockServerWithMemPinset() (*regclient.Client, *httptest.Server) {
	protek := handlers.NewNoopProtector()
	prof := registry.NewMemProfiles()
	ds := registry.NewMemDatasets()
	reg := registry.Registry{
		Profiles: prof,
		Datasets: ds,
		Pinset:   &registry.MemPinset{Profiles: prof},
		Search:   registry.MockSearch{ds},
	}
	s := httptest.NewServer(handlers.NewRoutes(protek, reg))
	c := regclient.NewClient(&regclient.Config{Location: s.URL})
	return c, s
}
