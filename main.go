package main

import (
	"context"
	"flag"
	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"hanaBFT/hlog"
	"hanaBFT/messages"
	p2p2 "hanaBFT/network/p2p"
	"hanaBFT/peer"
	"strings"
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

	// init routers
	p2p2.Init()

	// add self ip address to routers
	ipInfo := strings.Split(h.Addrs()[0].String(), "/")
	IP := ipInfo[2]
	Port := ipInfo[4]
	p2p2.MR.Add(h.ID().Pretty(), IP, Port)

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		hlog.Fatalf("When create a gossip router, an error occurred: %s", err)
	}

	// setup local mDNS discovery
	err = p2p2.SetupDiscovery(ctx, h)
	if err != nil {
		hlog.Fatalf("When setup local mDNS discovery, an error occurred: %s", err)
	}

	// join the organization
	org := p2p2.JoinOrganization(ctx, ps, h.ID(), *selfName, *orgName)

	node := peer.NewNode(h.ID(), org)

	node.Run()

	for i := 0; i < 100; i++ {
		k, v := messages.NextKeyValue(i)
		message := messages.Message{
			SenderID:    h.ID().Pretty(),
			SenderName:  *selfName,
			Timestamp:   time.Now().UnixNano(),
			ContentType: "command",
			Content: messages.Command{
				ProposerID:   h.ID().Pretty(),
				ProposerName: *selfName,
				CommandID:    i,
				Key:          k,
				Value:        v,
			},
		}
		node.Broadcast(message)
		time.Sleep(time.Second * 2)
	}
}
