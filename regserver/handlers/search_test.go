package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qri-io/registry"
)

func TestSearch(t *testing.T) {
	ds := registry.NewMemDatasets()
	s := httptest.NewServer(NewRoutes(NewNoopProtector(), registry.NewMemProfiles(), ds, &registry.MockSearch{ds}))

	cases := []struct {
		method      string
		endpoint    string
		contentType string
		params      *registry.SearchParams
		resStatus   int
	}{
		{"GET", "/search", "application/json", &registry.SearchParams{"abc", 0, 100}, 200},
	}

	for i, c := range cases {
		req, err := http.NewRequest(c.method, fmt.Sprintf("%s%s", s.URL, c.endpoint), nil)
		if err != nil {
			t.Errorf("case %d error creating request: %s", i, err.Error())
			continue
		}
		if c.contentType != "" {
			req.Header.Set("Content-Type", c.contentType)
		}
		if c.params != nil {
			data, err := json.Marshal(c.params)
			if err != nil {
				t.Errorf("error marshaling json body: %s", err.Error())
				return
			}
			req.Body = ioutil.NopCloser(bytes.NewReader([]byte(data)))
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("case %d unexpected error: %s", i, err)
			continue
		}
		if res.StatusCode != c.resStatus {
			t.Errorf("case %d res status mismatch. expected: %d, got: %d", i, c.resStatus, res.StatusCode)
			continue
		}
	}
}
