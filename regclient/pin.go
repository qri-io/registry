package regclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/libp2p/go-libp2p-crypto"
	"github.com/qri-io/registry"
)

// GetPinned checks if a given path is pinned to this registry
func (c Client) GetPinned(path string) (bool, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/pins?path=%s", c.cfg.Location, path), nil)
	if err != nil {
		return false, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}

	rs := struct {
		Data struct {
			Pinned bool
		}
	}{}
	if err := json.NewDecoder(res.Body).Decode(&rs); err != nil {
		return false, err
	}

	return rs.Data.Pinned, nil
}

// Pin requests a dataset be replicated on the registry
func (c Client) Pin(path string, privKey crypto.PrivKey, addrs []string) error {
	req, err := registry.NewPinRequest(path, privKey, addrs)
	if err != nil {
		return err
	}
	_, err = c.doJSONPinReq("POST", req)
	return err
}

// Unpin requests a dataset not be replicated to the registry
func (c Client) Unpin(path string, privKey crypto.PrivKey) error {
	req, err := registry.NewPinRequest(path, privKey, nil)
	if err != nil {
		return err
	}
	_, err = c.doJSONPinReq("DELETE", req)
	return nil
}

// doJSONProfileReq is a common wrapper for /profile endpoint requests
func (c Client) doJSONPinReq(method string, pr *registry.PinRequest) (*registry.PinRequest, error) {
	if c.cfg.Location == "" {
		return nil, ErrNoRegistry
	}

	data, err := json.Marshal(pr)
	if err != nil {
		fmt.Println("marshal err:", err.Error())
		return nil, err
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/pins", c.cfg.Location), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, registry.ErrPinsetNotSupported
	}

	// add response to an envelope
	env := struct {
		Data *registry.PinRequest
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
