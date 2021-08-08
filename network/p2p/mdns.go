// -author: Iridescent -time: 2021/8/8
package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"time"
)

const MAXPEERSNUM = 256

type MDNS struct {
	peerChan chan peer.AddrInfo
	frequency time.Duration
	ctx context.Context
	host host.Host
}

func (mdns *MDNS) HandlePeerFound(p peer.AddrInfo) {
	fmt.Println("new peer:", p)
	mdns.peerChan <- p
	err := mdns.host.Connect(mdns.ctx, p)
	if err != nil {
		panic(err)
	}
}

func InitMDNS(frequency time.Duration, ctx context.Context, host host.Host, site string) *MDNS {
	mdns := &MDNS{
		peerChan: make(chan peer.AddrInfo, MAXPEERSNUM),
		frequency: frequency,
		ctx: ctx,
		host: host,
	}
	service, err := discovery.NewMdnsService(ctx, host, frequency, site)
	if err != nil {
		panic(err)
	}
	service.RegisterNotifee(mdns)
	return mdns
}