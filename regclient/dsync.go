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
func (c *Client) DsyncSend(ctx context.Context, ng ipld.NodeGetter, mfst *dag.Manifest) error {
	remote := &dsync.HTTPRemote{
		URL: fmt.Sprintf("%s/dsync", c.cfg.Location),
	}

	snd, err := dsync.NewSend(ctx, ng, mfst, remote)
	if err != nil {
		return err
	}

	return snd.Do()
}

// DsyncFetch fetches an entire DAG designated by root path
func (c *Client) DsyncFetch(ctx context.Context, path string, ng ipld.NodeGetter, bapi coreiface.BlockAPI) error {
	remote := &dsync.HTTPRemote{
		URL: fmt.Sprintf("%s/dsync", c.cfg.Location),
	}

	fetch, err := dsync.NewFetch(ctx, path, ng, bapi, remote)
	if err != nil {
		return err
	}

	return fetch.Do()
}
