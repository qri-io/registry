package registry

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

var (
	// nowFunc is an internal function for getting timestamps
	nowFunc = func() time.Time { return time.Now() }
)

// Profiles is a map of profile data safe for concurrent use
// heavily inspired by sync.Map
type Profiles struct {
	sync.RWMutex
	internal map[string]*Profile
}

// NewProfiles allocates a new *Profiles map
func NewProfiles() *Profiles {
	return &Profiles{
		internal: make(map[string]*Profile),
	}
}

// Len returns the number of records in the map
func (ps *Profiles) Len() int {
	return len(ps.internal)
}

// Load fetches a profile from the list by key
func (ps *Profiles) Load(key string) (value *Profile, ok bool) {
	ps.RLock()
	result, ok := ps.internal[key]
	ps.RUnlock()
	return result, ok
}

// Range calls an iteration fuction on each element in the map until
// the end of the list is reached or iter returns true
func (ps *Profiles) Range(iter func(key string, p *Profile) (brk bool)) {
	ps.RLock()
	defer ps.RUnlock()
	for key, p := range ps.internal {
		if iter(key, p) {
			break
		}
	}
}

// SortedRange is like range but with deterministic key ordering
func (ps *Profiles) SortedRange(iter func(key string, p *Profile) (brk bool)) {
	ps.RLock()
	defer ps.RUnlock()
	keys := make([]string, len(ps.internal))
	i := 0
	for key := range ps.internal {
		keys[i] = key
		i++
	}
	sort.StringSlice(keys).Sort()
	for _, key := range keys {
		if iter(key, ps.internal[key]) {
			break
		}
	}
}

// Delete removes a record from Profiles at key
func (ps *Profiles) Delete(key string) {
	ps.Lock()
	delete(ps.internal, key)
	ps.Unlock()
}

// store adds an entry
func (ps *Profiles) store(key string, value *Profile) {
	ps.Lock()
	ps.internal[key] = value
	ps.Unlock()
}

// Register adds a profile to the list if it's valid and the desired handle isn't taken
func (ps *Profiles) Register(p *Profile) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if err := p.Verify(); err != nil {
		return err
	}

	if _, ok := ps.Load(p.Handle); ok {
		return fmt.Errorf("handle '%s' is taken", p.Handle)
	}

	prev := ""
	ps.Range(func(key string, profile *Profile) bool {
		if profile.ProfileID == p.ProfileID {
			prev = key
			return true
		}
		return false
	})

	if prev != "" {
		ps.Delete(prev)
	}

	p.Signature = ""
	p.Created = nowFunc()
	ps.store(p.Handle, p)
	return nil
}

// Deregister removes a profile from the registry if it exists
func (ps *Profiles) Deregister(p *Profile) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if err := p.Verify(); err != nil {
		return err
	}

	ps.Delete(p.Handle)
	return nil
}
