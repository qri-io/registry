package regclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p-crypto"
	"github.com/qri-io/registry"
	"github.com/qri-io/registry/pinset"
)

// Status checks if a given path is pinned to this registry
func (c Client) Status(path string) (s pinset.PinStatus, err error) {
	var (
		req *http.Request
		res *http.Response
	)

	req, err = http.NewRequest("GET", fmt.Sprintf("%s/pins/status?path=%s", c.cfg.Location, path), nil)
	if err != nil {
		return
	}

	res, err = c.httpClient.Do(req)
	if err != nil {
		return
	}

	rs := struct {
		Data pinset.PinStatus
	}{}
	if err = json.NewDecoder(res.Body).Decode(&rs); err != nil {
		return
	}

	return rs.Data, nil
}

// Pin requests a dataset be replicated on the registry
func (c Client) Pin(path string, privKey crypto.PrivKey, addrs []string) error {
	req, err := pinset.NewPinRequest(path, privKey, addrs)
	if err != nil {
		return err
	}
	status, err := c.doJSONPinReq("POST", req)
	if err != nil {
		return err
	}

	if status.Pinned {
		return nil
	} else if status.Error != "" {
		return fmt.Errorf(status.Error)
	}

	// poll, checking for pinned == true
	updates := c.statusPoll(path, stdPollInterval)

	for status := range updates {
		if status.Pinned {
			break
		} else if status.Error != "" {
			return fmt.Errorf(status.Error)
		}
	}
	return nil
}

// Unpin requests a dataset not be replicated to the registry
func (c Client) Unpin(path string, privKey crypto.PrivKey) error {
	req, err := pinset.NewPinRequest(path, privKey, nil)
	if err != nil {
		return err
	}
	_, err = c.doJSONPinReq("DELETE", req)
	return nil
}

// doJSONPinReq is a common wrapper for /pin endpoint requests
func (c Client) doJSONPinReq(method string, pr *pinset.PinRequest) (*pinset.PinStatus, error) {
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
		if strings.Contains(err.Error(), "no such host") {
			return nil, ErrNoRegistry
		}
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, registry.ErrPinsetNotSupported
	}

	// add response to an envelope
	env := struct {
		Data *pinset.PinStatus
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
	if env.Data == nil {
		return &pinset.PinStatus{}, nil
	}

	return env.Data, nil
}

const stdPollInterval = time.Duration(time.Second)

func (c Client) statusPoll(path string, interval time.Duration) (updates chan pinset.PinStatus) {
	tick := time.NewTicker(interval)
	updates = make(chan pinset.PinStatus)
	stopTick := make(chan bool)
	done := func() {
		tick.Stop()
		stopTick <- true
		close(updates)
	}

	go func() {
		for {
			select {
			case <-stopTick:
				return
			case <-tick.C:
				status, err := c.Status(path)
				if err != nil {
					status.Error = err.Error()
					go done()
				}
				updates <- status
				if status.Pinned {
					go done()
				}
			}
		}
	}()

	return updates
}
