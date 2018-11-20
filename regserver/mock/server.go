// Package mock provides a mock registry server for testing purposes
// it mocks the behaviour of a registry server with in-memory storage
package mock

import (
	"net/http/httptest"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/pinset"
	"github.com/qri-io/registry/regclient"
	"github.com/qri-io/registry/regserver/handlers"
)

// NewMockServer creates an in-memory mock server (with a pinset) without any access protection and
// a registry client to match
func NewMockServer() (*regclient.Client, *httptest.Server) {
	reg := NewMemRegistry()
	ps := &pinset.MemPinset{Profiles: reg.Profiles}
	return NewMockServerRegistryPinset(reg, ps)
}

// NewMockServerRegistry creates a mock server & client with a passed-in registry
func NewMockServerRegistry(reg registry.Registry) (*regclient.Client, *httptest.Server) {
	s := httptest.NewServer(handlers.NewRoutes(handlers.NewNoopProtector(), reg))
	c := regclient.NewClient(&regclient.Config{Location: s.URL})
	return c, s
}

// NewMockServerRegistryPinset creates a mock server & client with a passed-in registry and pinset
func NewMockServerRegistryPinset(reg registry.Registry, ps pinset.Pinset) (*regclient.Client, *httptest.Server) {
	s := httptest.NewServer(handlers.NewRoutesPinset(handlers.NewNoopProtector(), reg, ps))
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
