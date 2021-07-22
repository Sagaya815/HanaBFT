### libp2p的使用方法
1. 第一步，创建一个用于在`p2p`网络中用来和别人进行交流的`host`
```go
import github.com/libp2p/go-libp2p"
ctx := context.Background()
// create a new libp2p Host that listens on a random TCP port
h, err := libp2p.New(ctx, libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/0"))
```
2. 第二步，利用在第一步中得到的`host`，创建一个基于`gossip`通信协议的推送和订阅服务
```go
import pubsub "github.com/libp2p/go-libp2p-pubsub"
// ctx还是第一步里的ctx
ps, err := pubsub.NewGossipSub(ctx, h)
```
3. 第三步，利用第二步里得到的`ps`，创建一个`topic`，并加入它，然后订阅该`topic`
```go
// topicName是一个string
topic, err := ps.Join(topicName(topicName))
sub, err := topic.Subscribe()
```
4. 第四步，定义一个结构体，然后让这个结构体实现`HandlePeerFound(pi peer.AddrInfo)`方法，从而实现接口`Notifee`
```go
type discoveryNotifee struct {
	h host.Host
}
//HandlePeerFound方法其实就是在发现一个新的peer后，我们该去做哪些处理
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	// fmt.Printf("discovered new peer %s\n", pi.ID.Pretty())
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		// fmt.Printf("error connecting to peer %s : %s\n", pi.ID.Pretty(), err)
		hlog.Fatalf("When connect to a new peer, an error occurred: %s", err)
	}
}
```
5. 第五步，在第四步里，我们定义了`discoveryNotifee`结构体，并让它实现了`Notifee`接口，现在定义一个函数`SetupDiscovery(ctx context.Context, h host.Host)`，这个函数的作用是能为我们在第一步里得到的`host`创建一个`MDNS`探测服务，它可以自动地发现与自己在同一网络中的其他`peer`
```go
func SetupDiscovery(ctx context.Context, h host.Host) error {
	// DiscoveryInterval是time.Duration，代表dns的刷新时间间隔
	disc, err := discovery.NewMdnsService(ctx, h, DiscoveryInterval, DiscoveryServiceTag)
	if err != nil {
		return err
	}
	// discoveryNotifee是在第四步里定义的结构体
	n := discoveryNotifee{h: h}
	disc.RegisterNotifee(&n)
	return nil
}
```

### 小结
Command是由client发给主节点的和从节点的，所以这里可以利用net.Dial的方式建立tcp连接，然后给每个节点发送Command，主从节点之间和从节点之间通过p2p网络转发PBFT的三种消息：{PrePrepare Prepare Commit}