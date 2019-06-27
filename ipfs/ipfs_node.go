package ipfstests

import (
  "context"

  core "github.com/ipfs/go-ipfs/core"
  bootstrap "github.com/ipfs/go-ipfs/core/bootstrap"
)

type Node = core.IpfsNode

func NewIpfsNode() (*Node, error) {

  ctx := context.Background()
  ipfs, err := core.NewNode(ctx, &core.BuildCfg{
    Online: true,
  })
  if err != nil {
    return nil, err
  }

  err = ipfs.Bootstrap(bootstrap.DefaultBootstrapConfig)
  return ipfs, err
}
