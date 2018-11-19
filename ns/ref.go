// Package ns defines the qri dataset naming system
// This package is currently an experiment while we try to settle the distinction
// between a dataset.datasetPod and a reference. This code is extracted from
//  github.com/qri-io/qri/repo/ref.go
//
// This package will either move to a new location or be removed entirely
// before 2019-04-01
//
// initial RFC is here: https://github.com/qri-io/rfcs/blob/master/text/0006-dataset_naming.md
package ns

import (
	"fmt"
	"strings"

	"github.com/mr-tron/base58"
	multihash "github.com/multiformats/go-multihash"
)

var (
	// ErrEmptyRef indicates that the given reference is empty
	ErrEmptyRef = fmt.Errorf("repo: empty dataset reference")
)

// Ref encapsulates a reference to a dataset. This needs to exist to bind
// ways of referring to a dataset to a dataset itself, as datasets can't easily
// contain their own hash information, and names are unique on a per-repository
// basis.
// It's tempting to think this needs to be "bigger", supporting more fields,
// keep in mind that if the information is important at all, it should
// be stored as metadata within the dataset itself.
type Ref struct {
	// Peername of dataset owner
	Peername string `json:"peername,omitempty"`
	// ProfileID of dataset owner
	ProfileID string `json:"profileID,omitempty"`
	// Unique name reference for this dataset
	Name string `json:"name,omitempty"`
	// Content-addressed path for this dataset
	Path string `json:"path,omitempty"`
}

// String implements the Stringer interface for Ref
func (r Ref) String() (s string) {
	s = r.AliasString()
	if r.ProfileID != "" || r.Path != "" {
		s += "@"
	}
	if r.ProfileID != "" {
		s += r.ProfileID
	}
	if r.Path != "" {
		s += r.Path
	}
	return
}

// AliasString returns the alias components of a Ref as a string
func (r Ref) AliasString() (s string) {
	s = r.Peername
	if r.Name != "" {
		s += "/" + r.Name
	}
	return
}

// Match checks returns true if Peername and Name are equal,
// and/or path is equal
func (r Ref) Match(b Ref) bool {
	// fmt.Printf("\nr.Peername: %s b.Peername: %s\n", r.Peername, b.Peername)
	// fmt.Printf("\nr.Name: %s b.Name: %s\n", r.Name, b.Name)
	return (r.Path != "" && b.Path != "" && r.Path == b.Path) || (r.ProfileID == b.ProfileID || r.Peername == b.Peername) && r.Name == b.Name
}

// Equal returns true only if Peername Name and Path are equal
func (r Ref) Equal(b Ref) bool {
	return r.Peername == b.Peername && r.ProfileID == b.ProfileID && r.Name == b.Name && r.Path == b.Path
}

// IsPeerRef returns true if only Peername is set
func (r Ref) IsPeerRef() bool {
	return (r.Peername != "" || r.ProfileID != "") && r.Name == "" && r.Path == ""
}

// IsEmpty returns true if none of it's fields are set
func (r Ref) IsEmpty() bool {
	return r.Equal(Ref{})
}

// MustParseRef panics if the reference is invalid. Useful for testing
func MustParseRef(refstr string) Ref {
	ref, err := ParseRef(refstr)
	if err != nil {
		panic(err)
	}
	return ref
}

// ParseRef decodes a dataset reference from a string value
// It’s possible to refer to a dataset in a number of ways.
// The full definition of a dataset reference is as follows:
//     dataset_reference = peer_name/dataset_name@peer_id/network/hash
//
// we swap in defaults as follows, all of which are represented as
// empty strings:
//     network - defaults to /ipfs/
//     hash - tip of version history (latest known commit)
//
// these defaults are currently enforced by convention.
// TODO - make Dataset Ref parsing the responisiblity of the repo.Repo
// interface, replacing empty strings with actual defaults
//
// dataset names & hashes are disambiguated by checking if the input
// parses to a valid multihash after base58 decoding.
// through defaults & base58 checking the following should all parse:
//     peer_name/dataset_name
//     /network/hash
//     peername
//     peer_id
//     @peer_id
//     @peer_id/network/hash
//
// see tests for more exmples
//
// TODO - add validation that prevents peernames from being
// valid base58 multihashes and makes sure hashes are actually valid base58 multihashes
// TODO - figure out how IPFS CID's play into this
func ParseRef(ref string) (Ref, error) {
	if ref == "" {
		return Ref{}, ErrEmptyRef
	}

	var (
		// nameRefString string
		dsr = Ref{}
		err error
	)

	// if there is an @ symbol, we are dealing with a Ref
	// with an identifier
	atIndex := strings.Index(ref, "@")

	if atIndex != -1 {

		dsr.Peername, dsr.Name = parseAlias(ref[:atIndex])
		dsr.ProfileID, dsr.Path, err = parseIdentifiers(ref[atIndex+1:])

	} else {

		var peername, datasetname, pid bool
		toks := strings.Split(ref, "/")

		for i, tok := range toks {
			if isBase58Multihash(tok) {
				// first hash we encounter is a peerID
				if !pid {
					dsr.ProfileID = tok
					pid = true
					continue
				}

				if !isBase58Multihash(toks[i-1]) {
					dsr.Path = fmt.Sprintf("/%s/%s", toks[i-1], strings.Join(toks[i:], "/"))
				} else {
					dsr.Path = fmt.Sprintf("/ipfs/%s", strings.Join(toks[i:], "/"))
				}
				break
			}

			if !peername {
				dsr.Peername = tok
				peername = true
				continue
			}

			if !datasetname {
				dsr.Name = tok
				datasetname = true
				continue
			}

			dsr.Path = strings.Join(toks[i:], "/")
			break
		}
	}

	if dsr.ProfileID == "" && dsr.Peername == "" && dsr.Name == "" && dsr.Path == "" {
		err = fmt.Errorf("malformed Ref string: %s", ref)
		return dsr, err
	}

	// if dsr.ProfileID != "" {
	// 	if !isBase58Multihash(dsr.ProfileID) {
	// 		err = fmt.Errorf("invalid ProfileID: '%s'", dsr.ProfileID)
	// 		return dsr, err
	// 	}
	// }

	return dsr, err
}

func parseAlias(alias string) (peer, dataset string) {
	for i, tok := range strings.Split(alias, "/") {
		switch i {
		case 0:
			peer = tok
		case 1:
			dataset = tok
		}
	}
	return
}

func parseIdentifiers(ids string) (profileID, path string, err error) {

	toks := strings.Split(ids, "/")
	switch len(toks) {
	case 0:
		err = fmt.Errorf("malformed Ref identifier: %s", ids)
	case 1:
		if toks[0] != "" {
			profileID = toks[0]

			return
		}
	case 2:
		profileID = toks[0]

		if isBase58Multihash(toks[0]) && isBase58Multihash(toks[1]) {
			toks[1] = fmt.Sprintf("/ipfs/%s", toks[1])
		}

		path = toks[1]
	default:
		profileID = toks[0]

		path = fmt.Sprintf("/%s/%s", toks[1], toks[2])
		return
	}

	return
}

// TODO - this could be more robust?
func stripProtocol(ref string) string {
	if strings.HasPrefix(ref, "/ipfs/") {
		return ref[len("/ipfs/"):]
	}
	return ref
}

func isBase58Multihash(hash string) bool {
	data, err := base58.Decode(hash)
	if err != nil {
		return false
	}
	if _, err := multihash.Decode(data); err != nil {
		return false
	}

	return true
}

// CompareRef compares two Dataset Refs, returning an error
// describing any difference between the two references
func CompareRef(a, b Ref) error {
	if a.ProfileID != b.ProfileID {
		return fmt.Errorf("PeerID mismatch. %s != %s", a.ProfileID, b.ProfileID)
	}
	if a.Peername != b.Peername {
		return fmt.Errorf("Peername mismatch. %s != %s", a.Peername, b.Peername)
	}
	if a.Name != b.Name {
		return fmt.Errorf("Name mismatch. %s != %s", a.Name, b.Name)
	}
	if a.Path != b.Path {
		return fmt.Errorf("Path mismatch. %s != %s", a.Path, b.Path)
	}
	return nil
}
