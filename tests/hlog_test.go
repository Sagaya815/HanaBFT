package tests

import (
	"context"
	"github.com/libp2p/go-libp2p"
	"hanaBFT/hlog"
	"testing"
)

func TestHlogSetup(t *testing.T) {
	ctx := context.Background()
	h, err := libp2p.New(ctx, libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}
	hlog.Setup(h.ID())
}
