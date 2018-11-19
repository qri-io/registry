package ns

import (
	"bytes"
	"net/http"
	"testing"
)

func TestRefFromReq(t *testing.T) {
	cases := []struct {
		url      string
		expected Ref
		err      string
	}{
		{"http://localhost:2503/peername", Ref{Peername: "peername"}, ""},
		{"http://localhost:2503/peername?limit=10&offset=2", Ref{Peername: "peername"}, ""},
		{"http://localhost:2503/peername/datasetname", Ref{Peername: "peername", Name: "datasetname"}, ""},
		{"http://localhost:2503/peername/datasetname/at/ipfs/QmdWJ7RnFj3SdWW85mR4AYP17C8dRPD9eUPyTqUxVyGMgD", Ref{Peername: "peername", Name: "datasetname", Path: "/ipfs/QmdWJ7RnFj3SdWW85mR4AYP17C8dRPD9eUPyTqUxVyGMgD"}, ""},
		{"http://localhost:2503/peername/datasetname/at/ntwk/QmdWJ7RnFj3SdWW85mR4AYP17C8dRPD9eUPyTqUxVyGMgD", Ref{Peername: "peername", Name: "datasetname", Path: "/ntwk/QmdWJ7RnFj3SdWW85mR4AYP17C8dRPD9eUPyTqUxVyGMgD"}, ""},
		{"http://localhost:2503/peername/datasetname/at/map/QmdWJ7RnFj3SdWW85mR4AYP17C8dRPD9eUPyTqUxVyGMgD/dataset.json", Ref{Peername: "peername", Name: "datasetname", Path: "/map/QmdWJ7RnFj3SdWW85mR4AYP17C8dRPD9eUPyTqUxVyGMgD"}, ""},
		{"http://localhost:2503/peername/datasetname/at/ipfs/QmdWJ7RnFj3SdWW85mR4AYP17C8dRPD9eUPyTqUxVyGMgD", Ref{Peername: "peername", Name: "datasetname", Path: "/ipfs/QmdWJ7RnFj3SdWW85mR4AYP17C8dRPD9eUPyTqUxVyGMgD"}, ""},
		{"http://google.com:8000/peername/datasetname/at/ipfs/QmdWJ7RnFj3SdWW85mR4AYP17C8dRPD9eUPyTqUxVyGMgD", Ref{Peername: "peername", Name: "datasetname", Path: "/ipfs/QmdWJ7RnFj3SdWW85mR4AYP17C8dRPD9eUPyTqUxVyGMgD"}, ""},
		{"http://google.com:8000/peername", Ref{Peername: "peername"}, ""},
		// {"http://google.com/peername", Ref{Peername: "peername"}, ""},
		{"/peername", Ref{Peername: "peername"}, ""},
		{"http://www.fkjhdekaldschjxilujkjkjknwjkn.org/peername/datasetname/", Ref{Peername: "peername", Name: "datasetname"}, ""},
		{"http://example.com", Ref{}, ""},
		{"", Ref{}, ""},
	}

	for i, c := range cases {
		r, err := http.NewRequest("GET", c.url, bytes.NewReader(nil))
		if err != nil {
			t.Errorf("case %d, error making request: %s", i, err)
		}
		got, err := RefFromReq(r)
		if (c.err != "" && err == nil) || (err != nil && c.err != err.Error()) {
			t.Errorf("case %d, error mismatch: expected '%s' but got '%s'", i, c.err, err)
			continue
		}
		if err := CompareRef(got, c.expected); err != nil {
			t.Errorf("case %d: %s", i, err.Error())
		}
	}
}
