package registry

import (
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/jbenet/go-multihash"
	"github.com/libp2p/go-libp2p-crypto"
)

// ErrPinsetNotSupported is a cannonical error for a repository that does not
// support pinning
var ErrPinsetNotSupported = fmt.Errorf("pinset is not supported")

// Pinset is the interface for acting as a remote pinning service.
// Pinset implementations are expected to keep a store of PinStatus
// that callers can use to probe the progress of a request
type Pinset interface {
	// Pin can take a while, so Pin returns a channel of PinStatus to
	// deliver updates structs that should all share an ID for the request that
	// the Pinset implementation will provide
	Pin(req *PinRequest) (chan PinStatus, error)
	// Unpin removes a pin
	Unpin(req *PinRequest) error
	// Status gets the current pin state value for a given PinRequest
	Status(req *PinRequest) (PinStatus, error)
	// Pins lists pins within the range defoined by limit & offset in
	// lexographical order
	Pins(limit, offset int) ([]string, error)
	// PinLen returns the number of pins in the set
	PinLen() (int, error)
}

// PinRequest is a signed request to modify the status of a pin
type PinRequest struct {
	ProfileID     string
	Signature     string
	Path          string
	PeerAddresses []string
}

// PinStatus carries state about the status of a pin
type PinStatus struct {
	Path        string
	Pinned      bool
	Started     time.Time
	PctComplete float32
	// optional string representing status, intended to be shown
	// to userss
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

// MemPinset is a completely ficticious implementation of a pinset
// it shouldn't be used anywhere, by anyone, ever.
type MemPinset struct {
	pk       PinStatusStore
	Profiles Profiles
	pins     []string
}

// Pin a dataset
func (m *MemPinset) Pin(req *PinRequest) (chan PinStatus, error) {
	if len(m.pins) == 0 {
		m.pins = append(m.pins, req.Path)
	} else {
		for i, p := range m.pins {
			if req.Path > p {
				m.pins = append(append(m.pins[:i], req.Path), m.pins[i:]...)
				break
			}
		}
	}

	pc := make(chan PinStatus)
	go func() {
		prog := PinStatus{
			Path:        req.Path,
			PctComplete: 1.0,
			Pinned:      true,
		}

		m.pk.Set(prog)
		pc <- prog
		close(pc)
	}()
	return pc, nil
}

// PinProgress gives a hydrated, up-to-date progress struct for a given request ID
func (m *MemPinset) Status(req *PinRequest) (PinStatus, error) {
	ps := m.pk.Get(req.Path)
	if ps == nil {
		return PinStatus{}, fmt.Errorf("not found")
	}

	return *ps, nil
}

// Unpin a dataset
func (m *MemPinset) Unpin(req *PinRequest) error {
	for i, p := range m.pins {
		if p == req.Path {
			m.pk.Delete(req.Path)
			m.pins = append(m.pins[:i], m.pins[i+1:]...)
			return nil
		}
	}
	return nil
}

// Pinned gets the pin status of a path
func (m *MemPinset) Pinned(path string) (pinned bool, err error) {
	for _, p := range m.pins {
		if p == path {
			return true, nil
		}
	}
	return false, nil
}

// Pins reads from the list present in the pinset
func (m *MemPinset) Pins(limit, offset int) (pins []string, err error) {
	for i, p := range m.pins {
		if i > offset {
			continue
		}
		pins = append(pins, p)
		if len(pins) == limit {
			break
		}
	}
	return
}

// PinLen returns the number of pins in the pinset
func (m *MemPinset) PinLen() (int, error) {
	return len(m.pins), nil
}

// PinStatusStore is an in-memory key-value store of PinStatus
// using the Status "Path" field as keys
// To prevent the store from unbounded growth statuses are removed
// if the timespan since a status "Started" field exceeds the Ttl
// value of the store (defaults to 4 hours). Statuses are checked
// with a ticker the defaults to calling Sweep every 20 minutes
type PinStatusStore struct {
	Ttl time.Duration
	gc  *time.Ticker
	sync.Mutex
	store map[string]PinStatus
}

// Set a pin status
func (p *PinStatusStore) Set(ps PinStatus) {
	if p.store == nil {
		p.store = map[string]PinStatus{}
		if p.Ttl == 0 {
			p.Ttl = time.Hour * 4
		}
		p.StartGC(time.Minute * 20)
	}
	if ps.Started.IsZero() {
		ps.Started = time.Now()
	}

	p.Lock()
	p.store[ps.Path] = ps
	p.Unlock()
}

// Get a pin status
func (p *PinStatusStore) Get(path string) *PinStatus {
	p.Lock()
	defer p.Unlock()
	if ps, ok := p.store[path]; ok {
		return &ps
	}
	return nil
}

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

	for path, ps := range p.store {
		if ps.Started.Add(p.Ttl).Before(now) {
			remove = append(remove, path)
		}
	}

	for _, path := range remove {
		delete(p.store, path)
	}
}
