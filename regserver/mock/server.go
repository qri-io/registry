// Package mock provides a mock registry server for testing purposes
// it mocks the behaviour of a registry server with in-memory storage
package mock

import (
	"net/http/httptest"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/regclient"
	"github.com/qri-io/registry/regserver/handlers"
)

// var nilSearch registry.NilSearch
var mockSearch = &registry.MockSearch{}

// NewMockServer creates an in-memory mock server without any access protection and
// a registry client to match
func NewMockServer() (*regclient.Client, *httptest.Server) {
	s := httptest.NewServer(handlers.NewRoutes(handlers.NewNoopProtector(), registry.NewMemProfiles(), registry.NewMemDatasets(), mockSearch))
	c := regclient.NewClient(&regclient.Config{Location: s.URL})
	return c, s
}
