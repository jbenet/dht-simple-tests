package dhttests

import (
  "context"
  "fmt"
  "math/rand"
  "testing"
  "time"
)

var BgCtx context.Context

func init() {
  rand.Seed(time.Now().Unix())
  BgCtx = context.Background()
}

func TNode() *Node {
  n, err := NewNode()
  if err != nil {
    panic(err)
  }
  return n
}

func TNodes(_ *testing.T, n int) []*Node {
  ch := make(chan *Node, n)
  for i := 0; i < n; i++ {
    go func() {
      ch <- TNode()
    }()
  }

  ns := make([]*Node, n)
  for i := 0; i < n; i++ {
    ns[i] = <-ch
  }
  fmt.Printf("bootstrapped %d nodes\n", n)
  return ns
}

func Timed(_ *testing.T, s string, f func()) time.Duration {
  t1 := time.Now()
  f()
  td := time.Since(t1)
  fmt.Printf("%s runtime: %v\n", s, td)
  return td
}

func logDurationStats(t *testing.T, ds []time.Duration) {
  max := ds[0]
  min := ds[0]
  avg := time.Duration(0)

  for _, d := range ds {
    if d < min {
      min = d
    }
    if d > max {
      max = d
    }
    avg += d
  }
  avg = time.Duration(int(avg) / len(ds))
  t.Logf("max: %v, min: %v, avg: %v", max, min, avg)
}

func TestFindPeer1(t *testing.T) {
  var err error
  var ns []*Node

  Timed(t, "setup", func() {
    ns = TNodes(t, 2)
  })

  Timed(t, "query", func() {
    _, err = ns[0].DHT.FindPeer(BgCtx, ns[1].Host.ID())
  })

  if err != nil {
    t.Error("n0 failed to find n1", err)
  }
}

func TestFindPeer2(t *testing.T) {
  var err error
  var ns []*Node
  n := 25

  Timed(t, "setup", func() {
    ns = TNodes(t, n)
  })

  dsch := make(chan time.Duration, n)
  for i := 1; i < n; i++ {
    go func(i int) {
      d := Timed(t, fmt.Sprintf("n0 -> n%d", i), func() {
        _, err = ns[0].DHT.FindPeer(BgCtx, ns[i].Host.ID())
      })
      if err != nil {
        t.Errorf("n0 failed to find n%d. %v", i, err)
      }
      dsch <- d
    }(i)
  }

  var ds []time.Duration
  for i := 0; i < n; i++ {
    ds = append(ds, <-dsch)
  }
  logDurationStats(t, ds)
}
