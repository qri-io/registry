package regclient

import (
	"context"
	"fmt"

	"github.com/qri-io/dag"
	"github.com/qri-io/dag/dsync"

	ipld "github.com/ipfs/go-ipld-format"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
)

// DsyncSend sents an entire Manifest
func (c *Client) DsyncSend(ctx context.Context, lng ipld.NodeGetter, mfst *dag.Manifest) error {
	remote := &dsync.HTTPClient{
		URL: fmt.Sprintf("%s/dsync", c.cfg.Location),
	}

	snd, err := dsync.NewPush(lng, &dag.Info{Manifest: mfst}, remote, true)
	if err != nil {
		return err
	}

	return snd.Do(ctx)
}

// DsyncFetch fetches an entire DAG designated by root path
func (c *Client) DsyncFetch(ctx context.Context, path string, ng ipld.NodeGetter, bapi coreiface.BlockAPI) error {
	remote := &dsync.HTTPClient{
		URL: fmt.Sprintf("%s/dsync", c.cfg.Location),
	}

	fetch, err := dsync.NewPull(path, ng, bapi, remote, nil)
	if err != nil {
		return err
	}

	return fetch.Do(ctx)
}
