// Package mock provides a mock registry server for testing purposes
// it mocks the behaviour of a registry server with in-memory storage
package mock

import (
	"fmt"

	"net/http/httptest"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/pinset"
	"github.com/qri-io/registry/regclient"
	"github.com/qri-io/registry/regserver/handlers"
)

func init() {
	// don't need verbose logging when working with mock servers
	handlers.SetLogLevel("error")
}

// NewMockServer creates an in-memory mock server (with a pinset) without any access protection and
// a registry client to match
func NewMockServer() (*regclient.Client, *httptest.Server) {
	return NewMockServerWithNumDatasets(0)
}

// NewMockServerWithNumDatasets creates an in-memory mock server containing num datasets
func NewMockServerWithNumDatasets(num int) (*regclient.Client, *httptest.Server) {
	reg := NewMemRegistry()
	for i := 0; i < num; i++ {
		name := fmt.Sprintf("ds_%d", i)
		ds := &registry.Dataset{
			Path:   fmt.Sprintf("QmAbC%d", i),
			Name:   name,
			Handle: "peer",
		}
		reg.Datasets.Store(name, ds)
	}
	ps := &pinset.MemPinset{Profiles: reg.Profiles}
	return NewMockServerRegistryPinset(reg, ps)
}

// NewMockServerRegistry creates a mock server & client with a passed-in registry
func NewMockServerRegistry(reg registry.Registry) (*regclient.Client, *httptest.Server) {
	s := httptest.NewServer(handlers.NewRoutes(reg))
	c := regclient.NewClient(&regclient.Config{Location: s.URL})
	return c, s
}

// NewMockServerRegistryPinset creates a mock server & client with a passed-in registry and pinset
func NewMockServerRegistryPinset(reg registry.Registry, ps pinset.Pinset) (*regclient.Client, *httptest.Server) {
	s := httptest.NewServer(handlers.NewRoutes(reg, handlers.AddPinset(ps)))
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
