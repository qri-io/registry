package registry

import (
	"encoding/base64"
	// "math/rand"
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-crypto"
	"github.com/qri-io/dataset"
)

// Test Private Key. peerId: QmZePf5LeXow3RW5U1AgEiNbW46YnRGhZ7HPvm1UmPFPwt
var testPk = []byte(`CAASpgkwggSiAgEAAoIBAQC/7Q7fILQ8hc9g07a4HAiDKE4FahzL2eO8OlB1K99Ad4L1zc2dCg+gDVuGwdbOC29IngMA7O3UXijycckOSChgFyW3PafXoBF8Zg9MRBDIBo0lXRhW4TrVytm4Etzp4pQMyTeRYyWR8e2hGXeHArXM1R/A/SjzZUbjJYHhgvEE4OZy7WpcYcW6K3qqBGOU5GDMPuCcJWac2NgXzw6JeNsZuTimfVCJHupqG/dLPMnBOypR22dO7yJIaQ3d0PFLxiDG84X9YupF914RzJlopfdcuipI+6gFAgBw3vi6gbECEzcohjKf/4nqBOEvCDD6SXfl5F/MxoHurbGBYB2CJp+FAgMBAAECggEAaVOxe6Y5A5XzrxHBDtzjlwcBels3nm/fWScvjH4dMQXlavwcwPgKhy2NczDhr4X69oEw6Msd4hQiqJrlWd8juUg6vIsrl1wS/JAOCS65fuyJfV3Pw64rWbTPMwO3FOvxj+rFghZFQgjg/i45uHA2UUkM+h504M5Nzs6Arr/rgV7uPGR5e5OBw3lfiS9ZaA7QZiOq7sMy1L0qD49YO1ojqWu3b7UaMaBQx1Dty7b5IVOSYG+Y3U/dLjhTj4Hg1VtCHWRm3nMOE9cVpMJRhRzKhkq6gnZmni8obz2BBDF02X34oQLcHC/Wn8F3E8RiBjZDI66g+iZeCCUXvYz0vxWAQQKBgQDEJu6flyHPvyBPAC4EOxZAw0zh6SF/r8VgjbKO3n/8d+kZJeVmYnbsLodIEEyXQnr35o2CLqhCvR2kstsRSfRz79nMIt6aPWuwYkXNHQGE8rnCxxyJmxV4S63GczLk7SIn4KmqPlCI08AU0TXJS3zwh7O6e6kBljjPt1mnMgvr3QKBgQD6fAkdI0FRZSXwzygx4uSg47Co6X6ESZ9FDf6ph63lvSK5/eue/ugX6p/olMYq5CHXbLpgM4EJYdRfrH6pwqtBwUJhlh1xI6C48nonnw+oh8YPlFCDLxNG4tq6JVo071qH6CFXCIank3ThZeW5a3ZSe5pBZ8h4bUZ9H8pJL4C7yQKBgFb8SN/+/qCJSoOeOcnohhLMSSD56MAeK7KIxAF1jF5isr1TP+rqiYBtldKQX9bIRY3/8QslM7r88NNj+aAuIrjzSausXvkZedMrkXbHgS/7EAPflrkzTA8fyH10AsLgoj/68mKr5bz34nuY13hgAJUOKNbvFeC9RI5g6eIqYH0FAoGAVqFTXZp12rrK1nAvDKHWRLa6wJCQyxvTU8S1UNi2EgDJ492oAgNTLgJdb8kUiH0CH0lhZCgr9py5IKW94OSM6l72oF2UrS6PRafHC7D9b2IV5Al9lwFO/3MyBrMocapeeyaTcVBnkclz4Qim3OwHrhtFjF1ifhP9DwVRpuIg+dECgYANwlHxLe//tr6BM31PUUrOxP5Y/cj+ydxqM/z6papZFkK6Mvi/vMQQNQkh95GH9zqyC5Z/yLxur4ry1eNYty/9FnuZRAkEmlUSZ/DobhU0Pmj8Hep6JsTuMutref6vCk2n02jc9qYmJuD7iXkdXDSawbEG6f5C4MUkJ38z1t1OjA==`)

func init() {
	data, err := base64.StdEncoding.DecodeString(string(testPk))
	if err != nil {
		panic(err)
	}
	testPk = data
}

func TestDatasetsRegister(t *testing.T) {
	dss := NewMemDatasets()

	pk, err := crypto.UnmarshalPrivateKey(testPk)
	if err != nil {
		t.Error(err.Error())
		return
	}

	ts, err := time.Parse(time.RFC3339Nano, "2001-01-01T01:01:01.000000001Z")
	if err != nil {
		t.Errorf("invalid timestamp: %s", err.Error())
		return
	}

	// src := rand.New(rand.NewSource(0))
	// key0, _, err := crypto.GenerateSecp256k1Key(src)
	// if err != nil {
	// 	t.Error(err.Error())
	// 	return
	// }

	// mismatchSig, err := key0.Sign([]byte("bad_data"))
	// if err != nil {
	// 	t.Error(err.Error())
	// 	return
	// }

	ds1, err := NewDataset("foo", "bar", &dataset.DatasetPod{
		Path: "foo",
		Commit: &dataset.CommitPod{
			Timestamp: ts,
			Signature: "RZU/18bxxacveMoNvGxINIS9MxvNwtc4OiSCRjCGnospztHNhJfJP0PflrzKG1tqLGi+c4w94BJRmLR/I5YaVqqwm86vGkYhwDRuBEViuT4GlKCzVEFUk63fJsT9YmcUWlabqEnUW2l0O6p+RatfmumlKOleONMYy1woa5PbIzRGoITo4u9piYiV6RVRJ9bURjEU7cr8iVXcwO+YEw6qMCUBKUAok+yttjt+iYm0JLD9hPoQO14Vu4jWMFxByoLvVIEquEqnlgyuQGvelFfuApUI5goTftOcASANuTsnrOe6gq0HJxNN27kAYQujS3swspi7qVrL9X8v341YKu77fQ==",
		},
		Structure: &dataset.StructurePod{
			Checksum: "QmcCcPTqmckdXLBwPQXxfyW2BbFcUT6gqv9oGeWDkrNTyD",
		},
	}, pk.GetPublic())
	if err != nil {
		t.Error(err.Error())
		return
	}

	cases := []struct {
		ds                dataset.DatasetPod
		name, handle, err string
	}{
		{dataset.DatasetPod{}, "foo", "bar", "path is required"},
		{dataset.DatasetPod{
			Path: "QmFooPath",
			Commit: &dataset.CommitPod{
				Timestamp: ts,
				Signature: "RZU/18bxxacveMoNvGxINIS9MxvNwtc4OiSCRjCGnospztHNhJfJP0PflrzKG1tqLGi+c4w94BJRmLR/I5YaVqqwm86vGkYhwDRuBEViuT4GlKCzVEFUk63fJsT9YmcUWlabqEnUW2l0O6p+RatfmumlKOleONMYy1woa5PbIzRGoITo4u9piYiV6RVRJ9bURjEU7cr8iVXcwO+YEw6qMCUBKUAok+yttjt+iYm0JLD9hPoQO14Vu4jWMFxByoLvVIEquEqnlgyuQGvelFfuApUI5goTftOcASANuTsnrOe6gq0HJxNN27kAYQujS3swspi7qVrL9X8v341YKu77fQ==",
			},
			Structure: &dataset.StructurePod{
				Checksum: "bad",
			},
		}, "foo", "bar", "invalid signature: crypto/rsa: verification error"},
		{dataset.DatasetPod{
			Path: "QmFooPath",
			Commit: &dataset.CommitPod{
				Timestamp: ts,
				Signature: "RZU/18bxxacveMoNvGxINIS9MxvNwtc4OiSCRjCGnospztHNhJfJP0PflrzKG1tqLGi+c4w94BJRmLR/I5YaVqqwm86vGkYhwDRuBEViuT4GlKCzVEFUk63fJsT9YmcUWlabqEnUW2l0O6p+RatfmumlKOleONMYy1woa5PbIzRGoITo4u9piYiV6RVRJ9bURjEU7cr8iVXcwO+YEw6qMCUBKUAok+yttjt+iYm0JLD9hPoQO14Vu4jWMFxByoLvVIEquEqnlgyuQGvelFfuApUI5goTftOcASANuTsnrOe6gq0HJxNN27kAYQujS3swspi7qVrL9X8v341YKu77fQ==",
			},
			Structure: &dataset.StructurePod{
				Checksum: "QmcCcPTqmckdXLBwPQXxfyW2BbFcUT6gqv9oGeWDkrNTyD",
			},
		}, "foo", "bar", ""},
		{dataset.DatasetPod{
			Path: "QmFooPath2",
			Commit: &dataset.CommitPod{
				Timestamp: ts,
				Signature: "RZU/18bxxacveMoNvGxINIS9MxvNwtc4OiSCRjCGnospztHNhJfJP0PflrzKG1tqLGi+c4w94BJRmLR/I5YaVqqwm86vGkYhwDRuBEViuT4GlKCzVEFUk63fJsT9YmcUWlabqEnUW2l0O6p+RatfmumlKOleONMYy1woa5PbIzRGoITo4u9piYiV6RVRJ9bURjEU7cr8iVXcwO+YEw6qMCUBKUAok+yttjt+iYm0JLD9hPoQO14Vu4jWMFxByoLvVIEquEqnlgyuQGvelFfuApUI5goTftOcASANuTsnrOe6gq0HJxNN27kAYQujS3swspi7qVrL9X8v341YKu77fQ==",
			},
			Structure: &dataset.StructurePod{
				Checksum: "QmcCcPTqmckdXLBwPQXxfyW2BbFcUT6gqv9oGeWDkrNTyD",
			},
		}, "foo", "bar", ""},
		// {Datasets{DatasetsID: p.DatasetsID, Handle: p.Handle, Signature: p.Signature, PublicKey: "bad_data"}, "publickey base64 encoding: illegal base64 data at input byte 3"},
		// {Datasets{DatasetsID: p.DatasetsID, Handle: p.Handle, Signature: p.Signature, PublicKey: base64.StdEncoding.EncodeToString([]byte("bad_data"))}, "invalid publickey: unexpected EOF"},
		// {Datasets{DatasetsID: p.DatasetsID, Handle: p.Handle, PublicKey: p.PublicKey, Signature: "bad_data"}, "signature base64 encoding: illegal base64 data at input byte 3"},
		// {Datasets{DatasetsID: p.DatasetsID, Handle: p.Handle, PublicKey: p.PublicKey, Signature: base64.StdEncoding.EncodeToString([]byte("bad_data"))}, "invalid signature: malformed signature: no header magic"},
		// {Datasets{DatasetsID: p.DatasetsID, Handle: p.Handle, PublicKey: p.PublicKey, Signature: base64.StdEncoding.EncodeToString(mismatchSig)}, "mismatched signature"},
		// {*p, ""},
		// {*p, "handle 'key0' is taken"},
		// {*p2, ""},
	}

	for i, c := range cases {
		ds, err := NewDataset(c.name, c.handle, &c.ds, pk.GetPublic())
		if err != nil {
			t.Error(err.Error())
			return
		}

		err = dss.Register(ds)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: '%s', got: '%s'", i, c.err, err)
		}
	}

	if err := dss.Deregister(&Dataset{}); err == nil {
		t.Error("invalid dataset should error")
	}
	if err := dss.Deregister(ds1); err != nil {
		t.Errorf("error deregistering: %s", err.Error())
	}

	ds1.Commit.Signature = "bad"
	if err := dss.Deregister(ds1); err == nil {
		t.Error("unverifiable dataset should error")
	}
}

func TestDatasetsSortedRange(t *testing.T) {
	dss := NewMemDatasets()

	pk, err := crypto.UnmarshalPrivateKey(testPk)
	if err != nil {
		t.Error(err.Error())
		return
	}

	ts, err := time.Parse(time.RFC3339Nano, "2001-01-01T01:01:01.000000001Z")
	if err != nil {
		t.Errorf("invalid timestamp: %s", err.Error())
		return
	}

	handles := map[string]string{"a": "foo", "b": "bar"}
	for handle, name := range handles {
		p, err := NewDataset(handle, name, &dataset.DatasetPod{
			Path: "foo",
			Commit: &dataset.CommitPod{
				Timestamp: ts,
				Signature: "RZU/18bxxacveMoNvGxINIS9MxvNwtc4OiSCRjCGnospztHNhJfJP0PflrzKG1tqLGi+c4w94BJRmLR/I5YaVqqwm86vGkYhwDRuBEViuT4GlKCzVEFUk63fJsT9YmcUWlabqEnUW2l0O6p+RatfmumlKOleONMYy1woa5PbIzRGoITo4u9piYiV6RVRJ9bURjEU7cr8iVXcwO+YEw6qMCUBKUAok+yttjt+iYm0JLD9hPoQO14Vu4jWMFxByoLvVIEquEqnlgyuQGvelFfuApUI5goTftOcASANuTsnrOe6gq0HJxNN27kAYQujS3swspi7qVrL9X8v341YKu77fQ==",
			},
			Structure: &dataset.StructurePod{
				Checksum: "QmcCcPTqmckdXLBwPQXxfyW2BbFcUT6gqv9oGeWDkrNTyD",
			},
		}, pk.GetPublic())
		if err != nil {
			t.Error(err.Error())
			return
		}

		if err := dss.Register(p); err != nil {
			t.Error(err.Error())
			return
		}
	}

	if dss.Len() != len(handles) {
		t.Errorf("expected len to equal handle length. expected: %d, got: %d", len(handles), dss.Len())
		return
	}

	if _, ok := dss.Load("foo/a"); !ok {
		t.Errorf("expected foo/a to load")
		return
	}

	for iter := 0; iter < 100; iter++ {
		// i := 0
		// failed := false
		dss.SortedRange(func(key string, ds *Dataset) bool {
			t.Log(key, ds.Path)
			// if handles[i] != p.Handle {
			// 	t.Errorf("iter: %d sorted index %d mismatch. expected: %s, got: %s", iter, i, handles[i], p.Handle)
			// 	failed = true
			// 	return true
			// }
			// i++
			return false
		})
		break
	}
}
