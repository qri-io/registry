package regclient

import (
	"net/http/httptest"
	"testing"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/pinset"
	"github.com/qri-io/registry/regserver/handlers"
)

func TestPinRequests(t *testing.T) {
	ps := registry.NewMemProfiles()
	pins := &pinset.MemPinset{Profiles: ps}
	reg := registry.Registry{
		Profiles: ps,
		Datasets: registry.NewMemDatasets(),
	}
	ts := httptest.NewServer(handlers.NewRoutes(reg, handlers.AddPinset(pins)))
	c := NewClient(&Config{
		Location: ts.URL,
	})

	handle := "b5"
	err := c.PutProfile(handle, pk1)
	if err != nil {
		t.Error(err.Error())
	}

	path := "foo"
	status, err := c.Status(path)
	if err != nil {
		t.Error(err.Error())
	}
	if status.Pinned != false {
		t.Errorf("expected pinned '%s' to equal false", path)
	}

	if err := c.Pin(path, pk1, nil); err != nil {
		t.Error(err.Error())
	}

	status, err = c.Status(path)
	if err != nil {
		t.Error(err.Error())
	}
	if !status.Pinned {
		t.Errorf("expected pinned '%s' to equal true", path)
	}

	if err := c.Unpin(path, pk1); err != nil {
		t.Error(err.Error())
	}

	status, err = c.Status(path)
	if err != nil {
		t.Error(err.Error())
	}
	if status.Pinned != false {
		t.Errorf("expected pinned '%s' to equal false", path)
	}
}
