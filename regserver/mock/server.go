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
	return NewMockServerRegistry(NewMemRegistry())
}

// NewMockServerRegistry creates a mock server & client with a passed-in registry
func NewMockServerRegistry(reg registry.Registry) (*regclient.Client, *httptest.Server) {
	s := httptest.NewServer(handlers.NewRoutes(handlers.NewNoopProtector(), reg))
	c := regclient.NewClient(&regclient.Config{Location: s.URL})
	return c, s
}

// NewMemRegistry creates a new in-memory registry
func NewMemRegistry() registry.Registry {
	return registry.Registry{
		Profiles: registry.NewMemProfiles(),
		Datasets: registry.NewMemDatasets(),
	}
}

// NewMemRegistryPinset creates an in-memory registry without any access protection,
// a registry client, and an in-memory pinset
func NewMemRegistryPinset() registry.Registry {
	prof := registry.NewMemProfiles()
	ds := registry.NewMemDatasets()
	return registry.Registry{
		Profiles: prof,
		Datasets: ds,
		Pinset:   &registry.MemPinset{Profiles: prof},
		Search:   registry.MockSearch{Datasets: ds},
	}
}
