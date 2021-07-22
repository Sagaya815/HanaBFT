package p2p

import (
	"context"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"hanaBFT/hlog"
	"strings"
	"time"
)

const DiscoveryInterval = time.Second

const DiscoveryServiceTag = "HanaBFT-example"

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	// fmt.Printf("discovered new peer %s\n", pi.ID.Pretty())
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		// fmt.Printf("error connecting to peer %s : %s\n", pi.ID.Pretty(), err)
		hlog.Fatalf("When connect to a new peer, an error occurred: %s", err)
	}
	ID := pi.ID.Pretty()
	ipInfo := strings.Split(pi.Addrs[0].String(), "/")
	IP := ipInfo[2]
	Port := ipInfo[4]
	MR.Add(ID, IP, Port)
}

// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func SetupDiscovery(ctx context.Context, h host.Host) error {
	// setup mDNS discovery to find local peers
	disc, err := discovery.NewMdnsService(ctx, h, DiscoveryInterval, DiscoveryServiceTag)
	if err != nil {
		return err
	}
	n := discoveryNotifee{h: h}
	disc.RegisterNotifee(&n)
	return nil
}
