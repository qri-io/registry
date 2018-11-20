package mock

import (
	"testing"

	"github.com/qri-io/registry/pinset"
)

func TestMockServer(t *testing.T) {
	NewMockServer()
	NewMockServerRegistry(NewMemRegistry())

	reg := NewMemRegistry()
	ps := &pinset.MemPinset{Profiles: reg.Profiles}
	NewMockServerRegistryPinset(reg, ps)
}
