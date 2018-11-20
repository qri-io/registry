package pinset

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"testing"

	"github.com/libp2p/go-libp2p-crypto"
)

func ensurePinRequestEqual(a, b *PinRequest) error {
	if a.ProfileID != b.ProfileID {
		return fmt.Errorf("ProfileID mismatch: %s != %s", a.ProfileID, b.ProfileID)
	}
	if a.Signature != b.Signature {
		return fmt.Errorf("Signature mismatch: %s != %s", a.Signature, b.Signature)
	}
	if a.Path != b.Path {
		return fmt.Errorf("Path mismatch: %s != %s", a.Path, b.Path)
	}
	if len(a.PeerAddresses) != len(b.PeerAddresses) {
		return fmt.Errorf("PeerAddresses length mismatch: %d != %d", len(a.PeerAddresses), len(b.PeerAddresses))
	}
	return nil
}

func TestNewPinRequest(t *testing.T) {
	cases := []struct {
		path, base64Key string
		addrs           []string
		res             PinRequest
		err             string
	}{
		// TODO
	}

	for i, c := range cases {
		data, err := base64.StdEncoding.DecodeString(c.base64Key)
		if err != nil {
			t.Fatalf("case %d key base64 encoding error: %s", i, err.Error())
		}
		pk, err := crypto.UnmarshalPrivateKey(data)
		if err != nil {
			t.Fatalf("case %d pk unmarshall error: %s", i, err.Error())
		}
		pr, err := NewPinRequest(c.path, pk, c.addrs)

		if !(err == nil || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: '%s', got: '%s'", i, c.err, err)
			continue
		}

		if err = ensurePinRequestEqual(&c.res, pr); err != nil {
			t.Errorf("case %d: %s", i, err.Error())
			continue
		}
	}
}

func TestInsertSorted(t *testing.T) {
	cases := []struct {
		list   []string
		elem   string
		expect []string
	}{
		{[]string{}, "e", []string{"e"}},
		{[]string{"b"}, "e", []string{"b","e"}},
		{[]string{"m"}, "e", []string{"e","m"}},
		{[]string{"b","d","m"}, "e", []string{"b","d","e","m"}},
		{[]string{"m","p","x"}, "e", []string{"e","m","p","x"}},
	}

	for i, c := range cases {
		got := insertSorted(c.list, c.elem)
		if !reflect.DeepEqual(got, c.expect) {
			t.Errorf("case %d failed, got: %s, expected: %s", i, got, c.expect)
		}
	}
}
