package registry

import (
	"encoding/base64"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p-crypto"
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

// Register adds a profile to the list if it's valid and the desired handle isn't taken.
// Registree's must prove they have control of the private key by signing the desired handle,
// which is validated with a provided public key. Public key, handle, and date of
func (ps *Profiles) Register(p *Profile) error {
	if err := p.Validate(); err != nil {
		return err
	}

	pkbytes, err := base64.StdEncoding.DecodeString(p.PublicKey)
	if err != nil {
		return fmt.Errorf("publickey base64 encoding: %s", err.Error())
	}

	pubkey, err := crypto.UnmarshalPublicKey(pkbytes)
	if err != nil {
		return fmt.Errorf("invalid publickey: %s", err.Error())
	}

	sigbytes, err := base64.StdEncoding.DecodeString(p.Signature)
	if err != nil {
		return fmt.Errorf("signature base64 encoding: %s", err.Error())
	}

	valid, err := pubkey.Verify([]byte(p.Handle), sigbytes)
	if err != nil {
		return fmt.Errorf("invalid signature: %s", err.Error())
	}

	if !valid {
		return fmt.Errorf("mismatched signature")
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
