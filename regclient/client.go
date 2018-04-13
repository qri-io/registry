// Package regclient defines a client for interacting with a registry server
package regclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/libp2p/go-libp2p-crypto"
	"github.com/qri-io/registry"
)

// ErrNoRegistry indicates that no registry has been specified
// all client methods MUST return ErrNoRegistry for all method calls
// when config.Registry.Location is an empty string
var ErrNoRegistry = fmt.Errorf("registry: no registry specified")

// Client wraps a registry configuration with methods for interacting
// with the configured registry
type Client struct {
	cfg *Config
}

// Config encapsulates options for working with a registry
type Config struct {
	// Location is the URL base to call to
	Location string
}

// NewClient creates a registry from a provided Registry configuration
func NewClient(cfg *Config) *Client {
	return &Client{cfg}
}

// GetProfile fills in missing fields in p with registry data
func (r Client) GetProfile(p *registry.Profile) error {
	pro, err := r.doJSONProfileReq("GET", p)
	if err != nil {
		return err
	}
	*p = *pro
	return nil
}

// PutProfile adds a profile to the registry
func (r Client) PutProfile(handle string, privKey crypto.PrivKey) error {
	p, err := registry.ProfileFromPrivateKey(handle, privKey)
	if err != nil {
		return err
	}
	_, err = r.doJSONProfileReq("POST", p)
	return err
}

// DeleteProfile removes a profile from the registry
func (r Client) DeleteProfile(handle string, privKey crypto.PrivKey) error {
	p, err := registry.ProfileFromPrivateKey(handle, privKey)
	if err != nil {
		return err
	}
	_, err = r.doJSONProfileReq("DELETE", p)
	return nil
}

// doJSONProfileReq is a common wrapper for /profile endpoint requests
func (r Client) doJSONProfileReq(method string, p *registry.Profile) (*registry.Profile, error) {
	if r.cfg.Location == "" {
		return nil, ErrNoRegistry
	}

	data, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/profile", r.cfg.Location), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	env := struct {
		Data *registry.Profile
		Meta struct {
			Error  string
			Status string
			Code   int
		}
	}{}

	if err := json.NewDecoder(res.Body).Decode(&env); err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error %d: %s", res.StatusCode, env.Meta.Error)
	}

	return env.Data, nil
}
