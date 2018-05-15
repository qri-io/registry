package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qri-io/dataset"
	"github.com/qri-io/registry"
)

func TestDataset(t *testing.T) {
	s := httptest.NewServer(NewRoutes(NewNoopProtector(), registry.NewMemProfiles(), registry.NewMemDatasets()))

	b5, err := registry.ProfileFromPrivateKey("b5", privKey1)
	if err != nil {
		t.Errorf("error generating profile: %s", err.Error())
		return
	}

	data, err := ioutil.ReadFile("testdata/cities.dataset.json")
	if err != nil {
		t.Errorf("error reading dataset file: %s", err.Error())
		return
	}
	cds := &dataset.DatasetPod{}
	if err := json.Unmarshal(data, cds); err != nil {
		t.Errorf("error unmarshaling dataset json: %s", err.Error())
		return
	}

	name := "dataset"
	ds, err := registry.NewDataset(b5.Handle, name, cds, privKey1.GetPublic())
	if err != nil {
		t.Errorf("error creating dataset: %s", err.Error())
		return
	}

	type env struct {
		Data *registry.Dataset
		Meta struct {
			Code int
		}
	}

	cases := []struct {
		method      string
		contentType string
		dataset     *registry.Dataset
		resStatus   int
		res         *env
	}{
		{"OPTIONS", "", nil, http.StatusBadRequest, nil},
		{"OPTIONS", "application/json", nil, http.StatusBadRequest, nil},
		{"OPTIONS", "application/json", &registry.Dataset{Handle: "foo"}, http.StatusNotFound, nil},
		{"POST", "", nil, http.StatusBadRequest, nil},
		{"POST", "application/json", nil, http.StatusBadRequest, nil},
		{"POST", "application/json", &registry.Dataset{Handle: b5.Handle}, http.StatusBadRequest, nil},
		{"POST", "application/json", ds, http.StatusOK, nil},
		{"GET", "application/json", &registry.Dataset{Name: name, Handle: b5.Handle}, http.StatusOK, &env{Data: ds}},
		{"GET", "application/json", registry.NewDatasetRef("", "", "", ds.Path), http.StatusOK, &env{Data: ds}},
		{"GET", "application/json", registry.NewDatasetRef("", "", "", "foo"), http.StatusNotFound, nil},
		{"DELETE", "application/json", nil, http.StatusBadRequest, nil},
		{"DELETE", "application/json", &registry.Dataset{Handle: b5.Handle, Name: name}, http.StatusBadRequest, nil},
		{"DELETE", "application/json", ds, http.StatusOK, nil},
	}

	for i, c := range cases {
		req, err := http.NewRequest(c.method, fmt.Sprintf("%s/dataset", s.URL), nil)
		if err != nil {
			t.Errorf("case %d error creating request: %s", i, err.Error())
			continue
		}

		if c.contentType != "" {
			req.Header.Set("Content-Type", c.contentType)
		}
		if c.dataset != nil {
			data, err := json.Marshal(c.dataset)
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

		if c.res != nil {
			e := &env{}
			if err := json.NewDecoder(res.Body).Decode(e); err != nil {
				t.Errorf("case %d error reading response body: %s", i, err.Error())
				continue
			}

			// if len(e.Data) != len(c.res.Data) {
			//  t.Errorf("case %d reponse body mismatch. expected %d, got: %d", i, len(e.Data), len(c.res.Data))
			//  continue
			// }
			if e.Data.Handle != c.res.Data.Handle {
				t.Errorf("case %d reponse handle mismatch. expected %s, got: %s", i, e.Data.Handle, c.res.Data.Handle)
			}

			// TODO - check each response for profile matches
		}
	}
}

func TestDatasets(t *testing.T) {
	// s := httptest.NewServer(NewRoutes(NewNoopProtector(), registry.NewMemProfiles(), registry.NewMemDatasets()))

	// p1, err := registry.DatasetFromPrivateKey("b5", privKey1)
	// if err != nil {
	// 	t.Errorf("error generating profile: %s", err.Error())
	// 	return
	// }

	// p1Rename, err := registry.DatasetFromPrivateKey("b6", privKey1)
	// if err != nil {
	// 	t.Errorf("error generating profile: %s", err.Error())
	// 	return
	// }

	// type env struct {
	// 	Data []*registry.Dataset
	// 	Meta struct {
	// 		Code int
	// 	}
	// }

	// b5 := &registry.Dataset{
	// 	ProfileID: "QmZePf5LeXow3RW5U1AgEiNbW46YnRGhZ7HPvm1UmPFPwt",
	// 	Handle:    "b5",
	// 	PublicKey: "CAASpgIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC/7Q7fILQ8hc9g07a4HAiDKE4FahzL2eO8OlB1K99Ad4L1zc2dCg+gDVuGwdbOC29IngMA7O3UXijycckOSChgFyW3PafXoBF8Zg9MRBDIBo0lXRhW4TrVytm4Etzp4pQMyTeRYyWR8e2hGXeHArXM1R/A/SjzZUbjJYHhgvEE4OZy7WpcYcW6K3qqBGOU5GDMPuCcJWac2NgXzw6JeNsZuTimfVCJHupqG/dLPMnBOypR22dO7yJIaQ3d0PFLxiDG84X9YupF914RzJlopfdcuipI+6gFAgBw3vi6gbECEzcohjKf/4nqBOEvCDD6SXfl5F/MxoHurbGBYB2CJp+FAgMBAAE=",
	// }

	// b6 := &registry.Dataset{
	// 	ProfileID: "QmZePf5LeXow3RW5U1AgEiNbW46YnRGhZ7HPvm1UmPFPwt",
	// 	Handle:    "b6",
	// 	PublicKey: "CAASpgIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC/7Q7fILQ8hc9g07a4HAiDKE4FahzL2eO8OlB1K99Ad4L1zc2dCg+gDVuGwdbOC29IngMA7O3UXijycckOSChgFyW3PafXoBF8Zg9MRBDIBo0lXRhW4TrVytm4Etzp4pQMyTeRYyWR8e2hGXeHArXM1R/A/SjzZUbjJYHhgvEE4OZy7WpcYcW6K3qqBGOU5GDMPuCcJWac2NgXzw6JeNsZuTimfVCJHupqG/dLPMnBOypR22dO7yJIaQ3d0PFLxiDG84X9YupF914RzJlopfdcuipI+6gFAgBw3vi6gbECEzcohjKf/4nqBOEvCDD6SXfl5F/MxoHurbGBYB2CJp+FAgMBAAE=",
	// }

	// cases := []struct {
	// 	method      string
	// 	endpoint    string
	// 	contentType string
	// 	profile     *registry.Dataset
	// 	resStatus   int
	// 	res         *env
	// }{
	// 	{"GET", "/datasets", "", nil, http.StatusOK, &env{}},
	// 	{"POST", "/dataset", "application/json", p1, http.StatusOK, nil},
	// 	{"GET", "/datasets", "", nil, http.StatusOK, &env{Data: []*registry.Dataset{b5}}},
	// 	{"POST", "/dataset", "application/json", p1Rename, http.StatusOK, nil},
	// 	{"GET", "/datasets", "", nil, http.StatusOK, &env{Data: []*registry.Dataset{b6}}},
	// 	{"DELETE", "/dataset", "application/json", p1Rename, http.StatusOK, nil},
	// 	{"GET", "/datasets", "", nil, http.StatusOK, &env{Data: []*registry.Dataset{}}},
	// }

	// for i, c := range cases {
	// 	req, err := http.NewRequest(c.method, fmt.Sprintf("%s%s", s.URL, c.endpoint), nil)
	// 	if err != nil {
	// 		t.Errorf("case %d error creating request: %s", i, err.Error())
	// 		continue
	// 	}

	// 	if c.contentType != "" {
	// 		req.Header.Set("Content-Type", c.contentType)
	// 	}
	// 	if c.profile != nil {
	// 		data, err := json.Marshal(c.profile)
	// 		if err != nil {
	// 			t.Errorf("error marshaling json body: %s", err.Error())
	// 			return
	// 		}
	// 		req.Body = ioutil.NopCloser(bytes.NewReader([]byte(data)))
	// 	}

	// 	res, err := http.DefaultClient.Do(req)
	// 	if res.StatusCode != c.resStatus {
	// 		t.Errorf("case %d res status mismatch. expected: %d, got: %d", i, c.resStatus, res.StatusCode)
	// 		continue
	// 	}

	// 	if c.res != nil {
	// 		e := &env{}
	// 		if err := json.NewDecoder(res.Body).Decode(e); err != nil {
	// 			t.Errorf("case %d error reading response body: %s", i, err.Error())
	// 			continue
	// 		}

	// 		if len(e.Data) != len(c.res.Data) {
	// 			t.Errorf("case %d reponse body mismatch. expected %d, got: %d", i, len(e.Data), len(c.res.Data))
	// 			continue
	// 		}

	// 		// TODO - check each response for profile matches
	// 	}
	// }
}

func TestPostDataset(t *testing.T) {
	// 	un := "username"
	// 	pw := "password"
	// 	s := httptest.NewServer(NewRoutes(NewNoopProtector(), registry.NewMemProfiles(), registry.NewMemDatasets()))

	// 	const profiles = `[
	//   {
	//     "ProfileID": "QmamJUR83rGtDMEvugcC2gtLDx2nhZUTzpzhH6MA2Pb3Md",
	//     "Handle": "EDGI",
	//     "PublicKey": "CAASpgIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCmTFRx/6dKmoxje8AG+jFv94IcGUGnjrupa7XEr12J/c4ZLn3aPrD8F0tjRbstt1y/J+bO7Qb69DGiu2iSIqyE21nl2oex5+14jtxbupRq9jRTbpUHRj+y9I7uUDwl0E2FS1IQpBBfEGzDPIBVavxbhguC3O3XA7Aq7vea2lpJ1tWpr0GDRYSNmJAybkHS6k7dz1eVXFK+JE8FGFJi/AThQZKWRijvWFdlZvb8RyNFRHzpbr9fh38bRMTqhZpw/YGO5Ly8PNSiOOE4Y5cNUHLEYwG2/lpT4l53iKScsaOazlRkJ6NmkM1il7riCa55fcIAQZDtaAx+CT5ZKfmek4P5AgMBAAE=",
	//     "Created": "2018-05-01T22:31:18.288004308Z"
	//   },
	//   {
	//     "ProfileID": "QmSyDX5LYTiwQi861F5NAwdHrrnd1iRGsoEvCyzQMUyZ4W",
	//     "Handle": "b5",
	//     "PublicKey": "CAASpgIwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQC/W17VPFX+pjyM1MHI1CsvIe+JYQX45MJNITTd7hZZDX2rRWVavGXhsccmVDGU6ubeN3t6ewcBlgCxvyewwKhmZKCAs3/0xNGKXK/YMyZpRVjTWw9yPU9gOzjd9GuNJtL7d1Hl7dPt9oECa7WBCh0W9u2IoHTda4g8B2mK92awLOZTjXeA7vbhKKX+QVHKDxEI0U2/ooLYJUVxEoHRc+DUYNPahX5qRgJ1ZDP4ep1RRRoZR+HjGhwgJP+IwnAnO5cRCWUbZvE1UBJUZDvYMqW3QvDp+TtXwqUWVvt69tp8EnlBgfyXU91A58IEQtLgZ7klOzdSEJDP+S8HIwhG/vbTAgMBAAE=",
	//     "Created": "2018-04-19T22:10:49.909268968Z"
	//   }
	// ]`

	// 	req, err := http.NewRequest("POST", fmt.Sprintf("%s/datasets", s.URL), strings.NewReader(profiles))
	// 	if err != nil {
	// 		t.Error(err.Error())
	// 		return
	// 	}

	// 	req.Header.Set("Content-Type", "application/json")
	// 	req.SetBasicAuth(un, pw)

	// 	res, err := http.DefaultClient.Do(req)
	// 	if err != nil {
	// 		t.Error(err.Error())
	// 	}

	// 	if res.StatusCode != 200 {
	// 		t.Errorf("response status mismatch. expected 200, got: %d", res.StatusCode)
	// 	}

}
