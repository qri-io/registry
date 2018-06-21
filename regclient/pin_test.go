package regclient

import (
	"net/http/httptest"
	"testing"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/regserver/handlers"
)

func TestPinRequests(t *testing.T) {
	ps := registry.NewMemProfiles()
	pins := &registry.MemPinset{Profiles: ps}
	reg := registry.Registry{
		Profiles: ps,
		Datasets: registry.NewMemDatasets(),
		Pinset:   pins,
	}
	ts := httptest.NewServer(handlers.NewRoutes(handlers.NewNoopProtector(), reg))
	c := NewClient(&Config{
		Location: ts.URL,
	})

	handle := "b5"
	err := c.PutProfile(handle, pk1)
	if err != nil {
		t.Error(err.Error())
	}

	path := "/foo"
	pinned, err := c.GetPinned(path)
	if err != nil {
		t.Error(err.Error())
	}
	if pinned != false {
		t.Errorf("expected pinned '%s' to equal false", path)
	}

	if err := c.Pin(path, pk1, nil); err != nil {
		t.Error(err.Error())
	}

	pinned, err = c.GetPinned(path)
	if err != nil {
		t.Error(err.Error())
	}
	if !pinned {
		t.Errorf("expected pinned '%s' to equal true", path)
	}

	if err := c.Unpin(path, pk1); err != nil {
		t.Error(err.Error())
	}

	pinned, err = c.GetPinned(path)
	if err != nil {
		t.Error(err.Error())
	}
	if pinned != false {
		t.Errorf("expected pinned '%s' to equal false", path)
	}
}
