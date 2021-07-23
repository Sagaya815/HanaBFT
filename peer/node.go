package peer

import (
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"hanaBFT/hlog"
	"hanaBFT/messages"
	"hanaBFT/network/p2p"
	"hanaBFT/network/tcp"
	"hanaBFT/utils"
	"strconv"
	"strings"
)

type Node struct {
	PeersTransport map[string]*tcp.Transport
	MyTransport    *tcp.Transport
	ID             peer.ID
	Name           string
	*p2p.Organization
}

func NewNode(id peer.ID, name string, organization *p2p.Organization) *Node {
	myPortStr := p2p.MR.PeerRoutes()[id.Pretty()].Port
	myPortNum, err := strconv.Atoi(myPortStr)
	if err != nil {
		hlog.Fatalf("When I new node, an error occurred: %s", err)
	}
	myPort := myPortNum - 1
	port := strconv.Itoa(myPort)
	address := "tcp://:" + port
	myTransport := tcp.NewTransport(address)
	node := &Node{
		PeersTransport: make(map[string]*tcp.Transport),
		MyTransport:    myTransport,
		ID:             id,
		Name:           name,
		Organization:   organization,
	}
	node.MyTransport.Listen()
	go node.UpdatePeersTransport()
	return node
}

func (n *Node) UpdatePeersTransport() {
	for {
		select {
		case <-p2p.MR.Refresh:
			n.PeersTransport = make(map[string]*tcp.Transport)
			for id, route := range p2p.MR.PeerRoutes() {
				if n.ID.Pretty() == id {
					continue
				}
				peerPortStr := route.Port
				peerPortNum, err := strconv.Atoi(peerPortStr)
				if err != nil {
					hlog.Errorf("When I get %s's port, an error occurred: %s", id, err)
					continue
				}
				port := peerPortNum - 1
				address := "tcp://" + route.IP + ":" + strconv.Itoa(port)
				peerTransport := tcp.NewTransport(address)
				utils.Retry(peerTransport.Dial, 10)
				n.PeersTransport[id] = peerTransport
			}
		}
	}
}

func (n *Node) Send(msg messages.Message, id string) {
	hlog.Debugf("Send message to %s", id)
	peerTransport, ok := n.PeersTransport[id]
	if ok {
		peerTransport.Send(msg)
	} else {
		fmt.Println("no exists transport")
	}
}

func (n *Node) Broadcast(msg messages.Message) {
	for id, _ := range n.PeersTransport {
		n.Send(msg, id)
	}
}

func (n *Node) Run() {
	go n.MyTransport.HandleRecvMsg()
}

func (n *Node) AsClient() {
	n.MyTransport.RepliesChan = make(chan messages.Reply, 1024)
}

func (n *Node) AsReplica() {
	n.MyTransport.CommandsChan = make(chan messages.Command, 1024)
	go func() {
		for {
			select {
			case command := <-n.MyTransport.CommandsChan:
				reply := messages.Reply{
					SenderID:  n.ID.Pretty(),
					CommandID: command.CommandID,
				}
				message := messages.Message{
					SenderID:    n.ID.Pretty(),
					SenderName:  n.Name,
					ContentType: "reply",
					Content:     reply,
				}
				n.Send(message, command.ProposerID)
			}
		}
	}()
}

func Setup() *Node {
	selfName := flag.String("name", "", "peer's name")
	orgName := flag.String("org", "default-org", "organization's name")
	flag.Parse()

	ctx := context.Background()
	// create a new libp2p Host that listens on a random TCP port
	host, err := libp2p.New(ctx, libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
	if err != nil {
		panic(err)
	}

	hlog.Setup(host.ID())

	// init routes
	p2p.Init()

	// add self ip address to routes
	ipInfo := strings.Split(host.Addrs()[0].String(), "/")
	IP := ipInfo[2]
	Port := ipInfo[4]
	p2p.MR.Add(host.ID().Pretty(), IP, Port)

	// create a new PubSub service using the GossipSub router
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		hlog.Fatalf("When create a gossip router, an error occurred: %s", err)
	}

	// setup local mDNS discovery
	err = p2p.SetupDiscovery(ctx, host)
	if err != nil {
		hlog.Fatalf("When setup local mDNS discovery, an error occurred: %s", err)
	}

	// join the organization
	org := p2p.JoinOrganization(ctx, ps, host.ID(), *selfName, *orgName)

	node := NewNode(host.ID(), *selfName, org)

	node.Run()

	node.AsReplica()

	node.AsClient()

	return node
}
