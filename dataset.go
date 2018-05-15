package registry

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-crypto"
	"github.com/qri-io/dataset"
)

// Dataset is a registry's version of a dataset
type Dataset struct {
	dataset.DatasetPod
	Handle    string `json:",omitempty"`
	Name      string `json:",omitempty"`
	PublicKey string `json:",omitempty"`
}

// NewDataset creates a new dataset instance
func NewDataset(handle, name string, cds *dataset.DatasetPod, pubkey crypto.PubKey) (*Dataset, error) {
	pubb, err := pubkey.Bytes()
	if err != nil {
		return nil, err
	}

	return &Dataset{
		DatasetPod: *cds,
		PublicKey:  base64.StdEncoding.EncodeToString(pubb),
		Name:       name,
		Handle:     handle,
	}, nil
}

// NewDatasetRef creates a dataset with any known reference detail strings
func NewDatasetRef(peername, name, profileID, path string) *Dataset {
	return &Dataset{
		DatasetPod: dataset.DatasetPod{
			Path: path,
		},
		Handle: peername,
		Name:   name,
	}
}

// Validate is a sanity check that all required values are present
func (d *Dataset) Validate() error {
	if d.Handle == "" {
		return fmt.Errorf("handle is required")
	}
	if d.Name == "" {
		return fmt.Errorf("name is required")
	}
	if d.PublicKey == "" {
		return fmt.Errorf("publicKey is required")
	}
	if d.Path == "" {
		return fmt.Errorf("path is required")
	}
	if d.Commit == nil {
		return fmt.Errorf("commit is required")
	}
	if d.Structure == nil {
		return fmt.Errorf("structure is required")
	}
	return nil
}

// Key gives the string this dataset value should be keyed to
func (d *Dataset) Key() string {
	return fmt.Sprintf("%s/%s", d.Handle, d.Name)
}

// sigBytes gives the signable bytes from a dataset
func (d *Dataset) sigBytes() []byte {
	return []byte(fmt.Sprintf("%s\n%s", d.Commit.Timestamp.UTC().Format(time.RFC3339), d.Structure.Checksum))
}

// Verify checks a profile's proof of key ownership
// Registree's must prove they have control of the private key by signing the desired handle,
// which is validated with a provided public key. Public key, handle, and date of
func (d *Dataset) Verify() error {
	return verify(d.PublicKey, d.Commit.Signature, d.sigBytes())
}
