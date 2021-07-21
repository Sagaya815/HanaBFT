package main

import (
	"context"
	"flag"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"hanaBFT/hlog"
	"hanaBFT/p2p"
	"time"
)

func main() {
	selfName := flag.String("name", "", "peer's name")
	orgName := flag.String("org", "default-org", "organization's name")
	flag.Parse()

	ctx := context.Background()
	// create a new libp2p Host that listens on a random TCP port
	h, err := libp2p.New(ctx, libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}

	hlog.Setup(h.ID())

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		hlog.Fatalf("When create a gossip router, an error occurred: %s", err)
	}

	// setup local mDNS discovery
	err = p2p.SetupDiscovery(ctx, h)
	if err != nil {
		hlog.Fatalf("When setup local mDNS discovery, an error occurred: %s", err)
	}

	// join the organization
	p2p.JoinOrganization(ctx, ps, h.ID(), *selfName, *orgName)

	time.Sleep(time.Hour)
}
