package ipfstests

import (
  "context"

  core "github.com/ipfs/go-ipfs/core"
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

  err = ipfs.Bootstrap(core.DefaultBootstrapConfig)
  return ipfs, err
}
