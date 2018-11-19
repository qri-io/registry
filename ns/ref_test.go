package ns

import (
	"testing"
)

var cases = []struct {
	ref         Ref
	String      string
	AliasString string
}{
	{Ref{
		Peername: "peername",
	}, "peername", "peername"},
	{Ref{
		Peername: "peername",
		Name:     "datasetname",
	}, "peername/datasetname", "peername/datasetname"},

	{Ref{
		ProfileID: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
	}, "@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", ""},
	{Ref{
		Path: "/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1",
	}, "@/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", ""},

	{Ref{
		Peername:  "peername",
		Name:      "datasetname",
		ProfileID: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
	}, "peername/datasetname@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", "peername/datasetname"},
	{Ref{
		Peername:  "peername",
		Name:      "datasetname",
		ProfileID: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
		Path:      "/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1",
	}, "peername/datasetname@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "peername/datasetname"},

	{Ref{
		ProfileID: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
		Peername:  "lucille",
	}, "lucille@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", "lucille"},
	{Ref{
		ProfileID: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
		Peername:  "lucille",
		Name:      "ball",
		Path:      "/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1",
	}, "lucille/ball@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "lucille/ball"},

	{Ref{
		ProfileID: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
		Peername:  "bad_name",
	}, "bad_name@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", "bad_name"},
	// TODO - this used to be me@badId, which isn't very useful, but at least provided coding parity
	// might be worth revisiting
	{Ref{
		ProfileID: "badID",
		Peername:  "me",
	}, "me@badID", "me"},
}

func TestRefString(t *testing.T) {
	for i, c := range cases {
		if c.ref.String() != c.String {
			t.Errorf("case %d:\n%s\n%s", i, c.ref.String(), c.String)
			continue
		}
	}
}

func TestRefAliasString(t *testing.T) {
	for i, c := range cases {
		if c.ref.AliasString() != c.AliasString {
			t.Errorf("case %d:\n%s\n%s", i, c.ref.AliasString(), c.AliasString)
			continue
		}
	}
}

func TestParseRef(t *testing.T) {
	peernameRef := Ref{
		Peername: "peername",
	}

	nameRef := Ref{
		Peername: "peername",
		Name:     "datasetname",
	}

	peerIDRef := Ref{
		ProfileID: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
	}

	idNameRef := Ref{
		ProfileID: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
		Name:      "datasetname",
	}

	idFullRef := Ref{
		ProfileID: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
		Name:      "datasetname",
		Path:      "/network/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
	}

	idFullIPFSRef := Ref{
		Name:      "datasetname",
		ProfileID: "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
		Path:      "/ipfs/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
	}

	fullRef := Ref{
		Peername: "peername",
		Name:     "datasetname",
		Path:     "/network/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
	}

	fullIPFSRef := Ref{
		Peername: "peername",
		Name:     "datasetname",
		Path:     "/ipfs/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
	}

	pathOnlyRef := Ref{
		Path: "/network/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
	}

	ipfsOnlyRef := Ref{
		Path: "/ipfs/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y",
	}

	mapRef := Ref{
		Path: "/map/QmcQsi93yUryyWvw6mPyDNoKRb7FcBx8QGBAeJ25kXQjnC",
	}

	cases := []struct {
		input  string
		expect Ref
		err    string
	}{
		{"", Ref{}, "repo: empty dataset reference"},
		{"peername/", peernameRef, ""},
		{"peername", peernameRef, ""},

		{"QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/", peerIDRef, ""},
		{"/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", peerIDRef, ""},

		{"peername/datasetname/", nameRef, ""},
		{"peername/datasetname", nameRef, ""},
		{"peername/datasetname/@", nameRef, ""},
		{"peername/datasetname@", nameRef, ""},

		{"/datasetname@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", idNameRef, ""},
		{"/datasetname@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/", idNameRef, ""},
		{"/datasetname/@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", idNameRef, ""},

		{"peername/datasetname/@/network/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", fullRef, ""},
		{"peername/datasetname@/network/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", fullRef, ""},

		{"/datasetname@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/network/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", idFullRef, ""},
		{"/datasetname@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/network/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", idFullRef, ""}, // 15
		{"/datasetname/@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/ipfs/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", idFullIPFSRef, ""},
		{"/datasetname@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/network/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", idFullRef, ""},

		{"@/network/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", pathOnlyRef, ""},
		{"@/ipfs/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y", ipfsOnlyRef, ""},
		{"@/map/QmcQsi93yUryyWvw6mPyDNoKRb7FcBx8QGBAeJ25kXQjnC", mapRef, ""},

		{"peername/datasetname/@/network/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/junk/junk/...", fullRef, ""},
		{"peername/datasetname/@/ipfs/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/junk/junk/...", fullIPFSRef, ""},

		// TODO - restore. These have been removed b/c I didn't have time to make dem work properly - @b5
		// {"peername/datasetname@/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/junk/junk/...", fullIPFSRef, ""},
		// {"peername/datasetname@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/junk/junk/...", fullIPFSRef, ""},
		// {"@/network/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/junk/junk/...", pathOnlyRef, ""},
		// {"@network/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/junk/junk/...", pathOnlyRef, ""},
		// {"@/QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/junk/junk/...", ipfsOnlyRef, ""},
		// {"@QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D96w1L5qAhUM5Y/junk/junk/...", ipfsOnlyRef, ""},

		// {"peername/datasetname@network/bad_hash", Ref{}, "invalid ProfileID: 'network'"},
		// {"peername/datasetname@bad_hash/junk/junk..", Ref{}, "invalid ProfileID: 'bad_hash'"},
		// {"peername/datasetname@bad_hash", Ref{}, "invalid ProfileID: 'bad_hash'"},

		// {"@///*(*)/", Ref{}, "malformed Ref string: @///*(*)/"},
		// {"///*(*)/", Ref{}, "malformed Ref string: ///*(*)/"},
		// {"@", Ref{}, ""},
		// {"///@////", Ref{}, ""},
	}

	for i, c := range cases {
		got, err := ParseRef(c.input)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mismatch. expected: '%s', got: '%s'", i, c.err, err)
			continue
		}

		if err := CompareRef(got, c.expect); err != nil {
			t.Errorf("case %d: %s", i, err.Error())
		}
	}
}

func TestMatch(t *testing.T) {
	cases := []struct {
		a, b  string
		match bool
	}{
		{"a/b@/b/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "a/b@/b/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", true},
		{"a/b@/b/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "a/b@/b/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", true},
		{"QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/b@/b/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/b@/b/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", true},

		{"a/different_name@/b/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", "a/b@/b/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", true},
		{"different_peername/b@/b/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", "a/b@/b/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", true},
		{"different_peername/b@/b/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", "QmYCvbfNbCwFR45HiNP45rwJgvatpiW38D961L5qAhUM5Y/b@/b/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", true},
	}

	for i, c := range cases {
		a, err := ParseRef(c.a)
		if err != nil {
			t.Errorf("case %d error parsing dataset ref a: %s", i, err.Error())
			continue
		}
		b, err := ParseRef(c.b)
		if err != nil {
			t.Errorf("case %d error parsing dataset ref b: %s", i, err.Error())
			continue
		}

		gotA := a.Match(b)
		if gotA != c.match {
			t.Errorf("case %d a.Match", i)
			continue
		}

		gotB := b.Match(a)
		if gotB != c.match {
			t.Errorf("case %d b.Match", i)
			continue
		}
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		a, b  string
		equal bool
	}{
		{"a/b@/b/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "a/b@/b/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", true},
		{"a/b@/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "a/b@/ipfs/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", false},

		{"QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1/b@/b/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1/b@/b/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", true},
		{"QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1/b@/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1/b@/ipfs/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", false},

		{"a/different_name@/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "a/b@/ipfs/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", false},
		{"different_peername/b@/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "a/b@/ipfs/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", false},

		{"QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL/different_name@/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL/b@/ipfs/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", false},
		{"QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL/b@/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "a/b@/ipfs/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", false},
		{"QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL/b@/ipfs/QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1", "QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1/b@/ipfs/QmdJgfxj4rocm88PLeEididS7V2cc9nQosA46RpvAnWvDL", false},
	}

	for i, c := range cases {
		a, err := ParseRef(c.a)
		if err != nil {
			t.Errorf("case %d error parsing dataset ref a: %s", i, err.Error())
			continue
		}
		b, err := ParseRef(c.b)
		if err != nil {
			t.Errorf("case %d error parsing dataset ref b: %s", i, err.Error())
			continue
		}

		gotA := a.Equal(b)
		if gotA != c.equal {
			t.Errorf("case %d a.Equal", i)
			continue
		}

		gotB := b.Equal(a)
		if gotB != c.equal {
			t.Errorf("case %d b.Equal", i)
			continue
		}
	}
}

func TestIsEmpty(t *testing.T) {
	cases := []struct {
		ref   Ref
		empty bool
	}{
		{Ref{}, true},
		{Ref{Peername: "a"}, false},
		{Ref{Name: "a"}, false},
		{Ref{Path: "a"}, false},
		{Ref{ProfileID: "a"}, false},
	}

	for i, c := range cases {
		got := c.ref.IsEmpty()
		if got != c.empty {
			t.Errorf("case %d: %s", i, c.ref)
			continue
		}
	}
}

func TestCompareRefs(t *testing.T) {
	cases := []struct {
		a, b Ref
		err  string
	}{
		{Ref{}, Ref{}, ""},
		{Ref{Name: "a"}, Ref{}, "Name mismatch. a != "},
		{Ref{Peername: "a"}, Ref{}, "Peername mismatch. a != "},
		{Ref{Path: "a"}, Ref{}, "Path mismatch. a != "},
		{Ref{ProfileID: "QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1"}, Ref{}, "PeerID mismatch. QmRdexT18WuAKVX3vPusqmJTWLeNSeJgjmMbaF5QLGHna1 != "},
	}

	for i, c := range cases {
		err := CompareRef(c.a, c.b)
		if !(err == nil && c.err == "" || err != nil && err.Error() == c.err) {
			t.Errorf("case %d error mistmatch. expected: '%s', got: '%s'", i, c.err, err)
			continue
		}
	}
}
