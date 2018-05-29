package handlers

import (
  "bytes"
  "testing"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "encoding/json"

  "github.com/qri-io/registry"
)

func TestSearch(t *testing.T) {
  s := httptest.NewServer(NewRoutes(NewNoopProtector(), registry.NewMemProfiles(), registry.NewMemDatasets(), nilSearch))

  cases := []struct {
    method string
    endpoint string
    contentType string
    params *registry.SearchParams
    resStatus int

  }{
    {"GET", "/search", "application/json", &registry.SearchParams{"abc", 0, 100}, 400},
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
    if res.StatusCode != c.resStatus {
      t.Errorf("case %d res status mismatch. expected: %d, got: %d", i, c.resStatus, res.StatusCode)
      continue
    }
  }
}