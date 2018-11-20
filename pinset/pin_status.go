package pinset

import (
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	crypto "github.com/libp2p/go-libp2p-crypto"
	multihash "github.com/multiformats/go-multihash"
)

// PinRequest is a signed request to modify the status of a pin
type PinRequest struct {
	ProfileID     string
	Signature     string
	Path          string
	PeerAddresses []string
}

// PinStatus carries state about the status of a pin process
type PinStatus struct {
	Path        string
	Pinned      bool
	TTL         time.Time
	PctComplete float32
	// optional string representing status, intended to be shown to userss
	Status string
	Error  string
}

// NewPinRequest creates a pin request from a private key & path combo
func NewPinRequest(path string, privKey crypto.PrivKey, addrs []string) (*PinRequest, error) {
	pubkeybytes, err := privKey.GetPublic().Bytes()
	if err != nil {
		return nil, fmt.Errorf("error getting pubkey bytes: %s", err.Error())
	}

	mh, err := multihash.Sum(pubkeybytes, multihash.SHA2_256, 32)
	if err != nil {
		return nil, fmt.Errorf("error summing pubkey: %s", err.Error())
	}

	sig, err := privKey.Sign([]byte(path))
	if err != nil {
		return nil, fmt.Errorf("signing path: %s", err.Error())
	}

	return &PinRequest{
		ProfileID:     mh.B58String(),
		Signature:     base64.StdEncoding.EncodeToString(sig),
		Path:          path,
		PeerAddresses: addrs,
	}, nil
}

// PinStatusStore is an in-memory key-value store of PinStatus
// using the Status "Path" field as keys
// To prevent the store from unbounded growth statuses are removed
// if the timespan since a status "Started" field exceeds the TTL
// value of the store (defaults to 4 hours). Statuses are checked
// with a ticker the defaults to calling Sweep every 20 minutes
type PinStatusStore struct {
	gc *time.Ticker
	sync.Mutex
	store map[string]PinStatus
}

// Set a pin status
func (p *PinStatusStore) Set(ps PinStatus) {
	if p.store == nil {
		p.store = map[string]PinStatus{}
		p.StartGC(time.Minute * 20)
	}

	p.Lock()
	p.store[ps.Path] = ps
	p.Unlock()
}

// Get a pinjob by path
func (p *PinStatusStore) Get(path string) *PinStatus {
	p.Lock()
	defer p.Unlock()
	if ps, ok := p.store[path]; ok {
		return &ps
	}
	return nil
}

// Delete removes a pinjob from the store by path
func (p *PinStatusStore) Delete(path string) {
	p.Lock()
	defer p.Unlock()
	delete(p.store, path)
}

// StartGC begins the garbage collection ticker
func (p *PinStatusStore) StartGC(interval time.Duration) {
	if p.gc != nil {
		p.gc.Stop()
	}

	p.gc = time.NewTicker(interval)

	go func() {
		for range p.gc.C {
			p.Sweep()
		}
	}()
}

// StopGC halts the ticker
func (p *PinStatusStore) StopGC() {
	if p.gc != nil {
		p.gc.Stop()
	}
}

// Sweep checks for stale PinStatus and removes them
func (p *PinStatusStore) Sweep() {
	var (
		remove []string
		now    = time.Now()
	)

	p.Lock()
	defer p.Unlock()

	for path, pj := range p.store {
		if pj.TTL.After(now) {
			remove = append(remove, path)
		}
	}

	for _, path := range remove {
		delete(p.store, path)
	}
}
