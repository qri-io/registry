package regclient

import (
	"context"
	"fmt"

	"github.com/qri-io/dag"
	"github.com/qri-io/dag/dsync"

	ipld "gx/ipfs/QmR7TcHkR9nxkUorfi8XMTAMLUK7GiP64TWWBzY3aacc1o/go-ipld-format"
	coreiface "gx/ipfs/QmUJYo4etAQqFfSS2rarFAE97eNGB8ej64YkRT2SmsYD4r/go-ipfs/core/coreapi/interface"
)

// DsyncSend pushes a DAG to the registry
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
