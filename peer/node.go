package peer

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"hanaBFT/hlog"
	"hanaBFT/messages"
	"hanaBFT/network/p2p"
	"hanaBFT/network/tcp"
	"hanaBFT/utils"
	"strconv"
)

type Node struct {
	PeersTransport map[string]*tcp.Transport
	MyTransport    *tcp.Transport
	ID             peer.ID

	*p2p.Organization
}

func NewNode(id peer.ID, organization *p2p.Organization) *Node {
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
	peerTransport := n.PeersTransport[id]
	peerTransport.Send(msg)
}

func (n *Node) Broadcast(msg messages.Message) {
	for id, _ := range n.PeersTransport {
		n.Send(msg, id)
	}
}

func (n *Node) Run() {
	go n.MyTransport.HandleRecvMsg()
}
