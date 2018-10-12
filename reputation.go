package registry

import (
	"fmt"
)

// Reputation is record of the peers reputation on the network
// This is a stub to be filled in later

type Reputation struct {
	ProfileID string
	Rep       int
}

// NewReputation creates a new reputation. Reputations start at 0 for now
func NewReputation(id string) *Reputation {
	return &Reputation{
		ProfileID: id,
		Rep:       0,
	}
}

// Validate is a sanity check that all required values are present
func (r *Reputation) Validate() error {
	if r.ProfileID == "" {
		return fmt.Errorf("profileID is required")
	}

	return nil
}

// SetReputation sets the reputation of a given Reputation
func (r *Reputation) SetReputation(reputation int) {
	r.Rep = reputation
}

// Reputation gets the rep of a given Reputation
func (r *Reputation) Reputation() int {
	return r.Rep
}
