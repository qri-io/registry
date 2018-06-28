package registry

import (
	"encoding/base64"
	"fmt"

	"github.com/jbenet/go-multihash"
	"github.com/libp2p/go-libp2p-crypto"
)

// ErrPinsetNotSupported is a cannonical error for a repository that does not support pinning
var ErrPinsetNotSupported = fmt.Errorf("pinset is not supported")

// Pinset is the interface for acting as a remote pinning service
type Pinset interface {
	Pinned(path string) (bool, error)
	Pin(req *PinRequest) (chan Progress, error)
	Unpin(req *PinRequest) error
}

// PinRequest is a signed request to modify the status of a pin
type PinRequest struct {
	ProfileID     string
	Signature     string
	Path          string
	Pinned        bool
	PeerAddresses []string
}

// Progress is an update on the status of task completion
type Progress struct {
	Pct float32
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
	Profiles Profiles
	Pins     []string
}

// Pin a dataset
func (m *MemPinset) Pin(req *PinRequest) (chan Progress, error) {
	m.Pins = append(m.Pins, req.Path)
	return nil, nil
}

// Unpin a dataset
func (m *MemPinset) Unpin(req *PinRequest) error {
	for i, p := range m.Pins {
		if p == req.Path {
			m.Pins = append(m.Pins[:i], m.Pins[i+1:]...)
			return nil
		}
	}
	return nil
}

// Pinned gets the pin status of a path
func (m *MemPinset) Pinned(path string) (pinned bool, err error) {
	for _, p := range m.Pins {
		if p == path {
			return true, nil
		}
	}
	return false, nil
}
