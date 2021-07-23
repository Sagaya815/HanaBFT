// TODO
package main

import (
	"fmt"
	"hanaBFT/messages"
	"hanaBFT/network/p2p"
	"hanaBFT/peer"
	"sync"
	"time"
)

type client struct {
	*peer.Node
	latencies []time.Duration
}

func NewClient(node *peer.Node) *client {
	return &client{
		Node:      node,
		latencies: make([]time.Duration, 0),
	}
}

func (c *client) benchmark() {
	quorum := &Quorum{Vote: make(map[*messages.Command]int), lock: sync.RWMutex{}}
	done := make(chan bool)
	for i := 0; i < 10000; i++ {
		k, v := messages.NextKeyValue(i)
		command := &messages.Command{
			ProposerID:   c.ID.Pretty(),
			ProposerName: c.Name,
			CommandID:    i,
			Timestamp:    time.Now().UnixNano(),
			Key:          k,
			Value:        v,
		}
		message := messages.Message{
			SenderID:    c.ID.Pretty(),
			SenderName:  c.Name,
			ContentType: "command",
			Content:     *command,
		}
		quorum.Vote[command] = 0
		c.Broadcast(message)
	loop:
		for {
			select {
			case reply := <-c.MyTransport.RepliesChan:
				quorum.lock.Lock()
				ok := quorum.add(reply.CommandID)
				quorum.lock.Unlock()
				if ok {
					nowTime := time.Now().UnixNano()
					quorum.lock.RLock()
					c.latencies = append(c.latencies, time.Duration(nowTime-quorum.getCommand(reply.CommandID).Timestamp))
					quorum.lock.RUnlock()
					if len(c.latencies) == 10000 {
						close(done)
					}
					break loop
				}
			}
		}
	}
	<-done
	fmt.Println("总共发送了", len(c.latencies), "个任务")
	sum := 0
	for _, lantency := range c.latencies {
		sum += int(lantency)
	}
	fmt.Printf("平均时延是: %f\n", float64(sum)/float64(10000*1000000000))
	fmt.Printf("总的时延是: %f\n", float64(sum)/float64(1000000000))
}

func main() {
	node := peer.Setup()
	client := NewClient(node)
	time.Sleep(time.Second * 5)
	client.benchmark()
}

/*
定义一些临时变量
*/

type Quorum struct {
	Vote map[*messages.Command]int
	lock sync.RWMutex
}

func (q *Quorum) add(commandID int) bool {
	for command, _ := range q.Vote {
		if command.CommandID == commandID {
			q.Vote[command]++
			if q.Vote[command] == p2p.MR.PeerNum()-1 {
				return true
			}
		}
	}
	return false
}

func (q *Quorum) getCommand(commandID int) *messages.Command {
	for command, _ := range q.Vote {
		if command.CommandID == commandID {
			return command
		}
	}
	return nil
}
