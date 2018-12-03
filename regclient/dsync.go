package regclient

import (
	"context"
	"fmt"

	"github.com/qri-io/dag"
	"github.com/qri-io/dag/dsync"

	ipld "gx/ipfs/QmR7TcHkR9nxkUorfi8XMTAMLUK7GiP64TWWBzY3aacc1o/go-ipld-format"
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
