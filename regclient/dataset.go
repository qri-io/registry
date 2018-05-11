package regclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/libp2p/go-libp2p-crypto"
	"github.com/qri-io/dataset"
	"github.com/qri-io/registry"
)

// GetDataset fetches a dataset from a registry
// func (c Client) GetDataset(peername, dsname, profileID, hash string) (*dataset.Dataset, error) {

// }

// PutDataset adds a dataset to a registry
func (c Client) PutDataset(peername, dsname string, ds *dataset.DatasetPod, pubKey crypto.PubKey) error {
	d, err := registry.NewDataset(peername, dsname, ds, pubKey)
	if err != nil {
		return err
	}

	_, err = c.doJSONDatasetReq("POST", d)
	return err
}

// DeleteDataset removes a dataset from the registry
func (c Client) DeleteDataset(peername, dsname string, ds *dataset.DatasetPod, pubKey crypto.PubKey) error {
	d, err := registry.NewDataset(peername, dsname, ds, pubKey)
	if err != nil {
		return err
	}

	_, err = c.doJSONDatasetReq("DELETE", d)
	return err
}

func (c Client) doJSONDatasetReq(method string, d *registry.Dataset) (*registry.Dataset, error) {
	if c.cfg.Location == "" {
		return nil, ErrNoRegistry
	}

	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s/dataset", c.cfg.Location), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	env := struct {
		Data *registry.Dataset
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
