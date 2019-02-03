package registry

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-crypto"
	"github.com/qri-io/dataset"
)

func TestDatasetValidate(t *testing.T) {
	cases := []struct {
		dataset *Dataset
		err     string
	}{
		{&Dataset{}, "handle is required"},
		{&Dataset{Handle: "foo"}, "name is required"},
		{&Dataset{Handle: "foo", Name: "bar"}, "publicKey is required"},
		{&Dataset{Handle: "foo", Name: "bar", PublicKey: "baz"}, "path is required"},
		{&Dataset{Handle: "foo", Name: "bar", PublicKey: "baz", Path: "bat"}, "commit is required"},
		{&Dataset{Handle: "foo", Name: "bar", PublicKey: "baz", Path: "bat", Commit: &dataset.Commit{}}, "structure is required"},
		{&Dataset{Handle: "foo", Name: "bar", PublicKey: "baz", Path: "bat", Commit: &dataset.Commit{}, Structure: &dataset.Structure{}}, ""},
	}

	for i, c := range cases {
		err := c.dataset.Validate()
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d mismatch. expected: '%s', got: '%s'", i, c.err, err)
			continue
		}
	}
}

func TestNewDataset(t *testing.T) {
	pubB58 := "CAISIQPpsxS5rdL3TVPQ+JWHa4cCGyjjuuY3chmSfa+Cw9V1aA=="
	data, err := base64.StdEncoding.DecodeString(pubB58)
	if err != nil {
		t.Error(err.Error())
		return
	}

	pub, err := crypto.UnmarshalPublicKey(data)
	if err != nil {
		t.Error(err.Error())
		return
	}

	ds, err := NewDataset("foo", "bar", &dataset.Dataset{}, pub)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	if ds.PublicKey != pubB58 {
		t.Errorf("publicKey mistmatch. expected '%s', got: '%s'", pubB58, ds.PublicKey)
	}
}

func TestNewDatasetRef(t *testing.T) {
	ref := NewDatasetRef("peername", "name", "profileID", "path")
	if ref.Handle != "peername" {
		t.Errorf("expected handle to equal peername. got: %s", ref.Handle)
	}
	if ref.Name != "name" {
		t.Errorf("expected name to equal name, got: %s", ref.Name)
	}
	if ref.Path != "path" {
		t.Errorf("expected path to equal path, got: %s", ref.Path)
	}
}

func TestDatasetVerify(t *testing.T) {
	ts, err := time.Parse(time.RFC3339Nano, "2001-01-01T01:01:01.000000001Z")
	if err != nil {
		t.Errorf("invalid timestamp: %s", err.Error())
		return
	}

	ds := &Dataset{
		PublicKey: "CAASpgIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC/7Q7fILQ8hc9g07a4HAiDKE4FahzL2eO8OlB1K99Ad4L1zc2dCg+gDVuGwdbOC29IngMA7O3UXijycckOSChgFyW3PafXoBF8Zg9MRBDIBo0lXRhW4TrVytm4Etzp4pQMyTeRYyWR8e2hGXeHArXM1R/A/SjzZUbjJYHhgvEE4OZy7WpcYcW6K3qqBGOU5GDMPuCcJWac2NgXzw6JeNsZuTimfVCJHupqG/dLPMnBOypR22dO7yJIaQ3d0PFLxiDG84X9YupF914RzJlopfdcuipI+6gFAgBw3vi6gbECEzcohjKf/4nqBOEvCDD6SXfl5F/MxoHurbGBYB2CJp+FAgMBAAE=",
		Commit: &dataset.Commit{
			Timestamp: ts,
			Signature: "RZU/18bxxacveMoNvGxINIS9MxvNwtc4OiSCRjCGnospztHNhJfJP0PflrzKG1tqLGi+c4w94BJRmLR/I5YaVqqwm86vGkYhwDRuBEViuT4GlKCzVEFUk63fJsT9YmcUWlabqEnUW2l0O6p+RatfmumlKOleONMYy1woa5PbIzRGoITo4u9piYiV6RVRJ9bURjEU7cr8iVXcwO+YEw6qMCUBKUAok+yttjt+iYm0JLD9hPoQO14Vu4jWMFxByoLvVIEquEqnlgyuQGvelFfuApUI5goTftOcASANuTsnrOe6gq0HJxNN27kAYQujS3swspi7qVrL9X8v341YKu77fQ==",
		},
		Structure: &dataset.Structure{
			Checksum: "QmcCcPTqmckdXLBwPQXxfyW2BbFcUT6gqv9oGeWDkrNTyD",
		},
	}

	if err := ds.Verify(); err != nil {
		t.Errorf("unexpected error: '%s'", err.Error())
	}
}
