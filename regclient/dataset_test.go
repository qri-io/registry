package regclient

import (
	"testing"
	"time"

	"github.com/qri-io/dataset"
	regmock "github.com/qri-io/registry/regserver/mock"
)

func TestDatasetRequests(t *testing.T) {
	handle := "b5"
	srv := regmock.NewMockServer()
	c := NewClient(&Config{
		Location: srv.URL,
	})

	// err = c.GetDataset()
	// if err == nil {
	// 	t.Errorf("expected empty get to error")
	// } else if err.Error() != "error 404: " {
	// 	t.Errorf("error mistmatch. expected: %s, got: %s", "error 404: ", err.Error())
	// }

	ts, err := time.Parse(time.RFC3339Nano, "2001-01-01T01:01:01.000000001Z")
	if err != nil {
		t.Errorf("invalid timestamp: %s", err.Error())
		return
	}

	ds := &dataset.DatasetPod{
		Path: "/map/QmYXMg6gqMAT8seUFhgAagknFvfs71auFWbnSfVcg1NTd8",
		Commit: &dataset.CommitPod{
			Timestamp: ts,
			Signature: "RZU/18bxxacveMoNvGxINIS9MxvNwtc4OiSCRjCGnospztHNhJfJP0PflrzKG1tqLGi+c4w94BJRmLR/I5YaVqqwm86vGkYhwDRuBEViuT4GlKCzVEFUk63fJsT9YmcUWlabqEnUW2l0O6p+RatfmumlKOleONMYy1woa5PbIzRGoITo4u9piYiV6RVRJ9bURjEU7cr8iVXcwO+YEw6qMCUBKUAok+yttjt+iYm0JLD9hPoQO14Vu4jWMFxByoLvVIEquEqnlgyuQGvelFfuApUI5goTftOcASANuTsnrOe6gq0HJxNN27kAYQujS3swspi7qVrL9X8v341YKu77fQ==",
		},
		Structure: &dataset.StructurePod{
			Checksum: "QmcCcPTqmckdXLBwPQXxfyW2BbFcUT6gqv9oGeWDkrNTyD",
		},
	}

	err = c.PutDataset(handle, "b", ds, pk1.GetPublic())
	if err != nil {
		t.Error(err.Error())
	}

	// err = c.GetDataset()
	// if err != nil {
	// 	t.Error(err)
	// }

	err = c.DeleteDataset(handle, "b", ds, pk1.GetPublic())
	if err != nil {
		t.Error(err.Error())
	}
}
