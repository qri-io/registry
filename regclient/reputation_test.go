package regclient

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/qri-io/registry"
	"github.com/qri-io/registry/regserver/handlers"
)

func TestReputationRequests(t *testing.T) {
	profileID := "my_id"
	profileID2 := "my_second_id"
	memRs := registry.NewMemReputations()
	newRep := registry.NewReputation(profileID)
	newRep.SetReputation(-1)
	err := memRs.Add(newRep)
	if err != nil {
		t.Error(err)
	}

	reg := registry.Registry{
		Reputations: memRs,
	}
	ts := httptest.NewServer(handlers.NewRoutes(handlers.NewNoopProtector(), reg))
	c := NewClient(&Config{
		Location: ts.URL,
	})

	rep, err := c.GetReputation(profileID)
	if err != nil {
		t.Error(err)
		return
	}
	if -1 != rep.Reputation() {
		t.Error(fmt.Errorf("reputation value not equal: expect -1, got %d", rep.Reputation()))
	}

	rep, err = c.GetReputation(profileID2)
	if err != nil {
		t.Error(err)
		return
	}
	if 1 != rep.Reputation() {
		t.Errorf("reputation value not equal: expect 1, got %d", rep.Reputation())
	}
	if memRs.Len() != 2 {
		t.Errorf("reputations list should equal 2, got %d", memRs.Len())
	}
}
