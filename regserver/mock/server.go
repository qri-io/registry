// Package mock provides a mock registry server for testing purposes
// it mocks the behaviour of a registry server with in-memory storage
package mock

import (
  "net/http/httptest"

  "github.com/qri-io/registry"
  "github.com/qri-io/registry/regserver/handlers"
)

// NewMockServer creates an in-memory mock server without any access protection
func NewMockServer() *httptest.Server {
  return httptest.NewServer(handlers.NewRoutes(handlers.NewNoopProtector(), registry.NewMemProfiles(), registry.NewMemDatasets()))
}