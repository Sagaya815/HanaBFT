package p2p

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"hanaBFT/hlog"
	"hanaBFT/messages"
	"time"
)

const MessagesBufSize = 1024

type Organization struct {
	Messages         chan *messages.Message
	ctx              context.Context
	ps               *pubsub.PubSub
	topic            *pubsub.Topic
	sub              *pubsub.Subscription
	organizationName string
	selfID           peer.ID
	selfName         string
}

func topicName(organizationName string) string {
	return "HanaBFT-" + organizationName
}

func JoinOrganization(ctx context.Context, ps *pubsub.PubSub, selfID peer.ID, selfName, organizationName string) *Organization {
	topic, err := ps.Join(topicName(organizationName))
	if err != nil {
		hlog.Fatalf("When we joined %s topic, an error occurred: %s", topicName(organizationName), err)
	}
	sub, err := topic.Subscribe()
	if err != nil {
		hlog.Fatalf("When we subscribed to %s topic, an error occurred: %s", topicName(organizationName), err)
	}
	org := &Organization{
		Messages:         make(chan *messages.Message, MessagesBufSize),
		ctx:              ctx,
		ps:               ps,
		topic:            topic,
		sub:              sub,
		organizationName: organizationName,
		selfID:           selfID,
		selfName:         selfName,
	}
	go org.recvMsgLoop()
	go org.handleRecvMessage()
	go org.refreshPeers()
	return org
}

func (org *Organization) recvMsgLoop() {
	for {
		msg, err := org.sub.Next(org.ctx)
		if err != nil {
			hlog.Errorf("When organization received a message, an error occurred: %s", err)
			close(org.Messages)
			return
		}
		if msg.ReceivedFrom == org.selfID {
			continue
		}
		message := new(messages.Message)
		err = json.Unmarshal(msg.Data, message)
		if err != nil {
			hlog.Errorf("When organization unmarshal coming message, an error occurred: %s", err)
			continue
		}
		org.Messages <- message
	}
}

func (org *Organization) Publish(msg messages.Message) {
	hlog.Debugf("%s broadcast message %v", org.selfName, msg)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		hlog.Errorf("When I marshal a message, an error occurred: %s", err)
		return
	}
	err = org.topic.Publish(org.ctx, msgBytes)
	if err != nil {
		hlog.Errorf("When I broadcast a message, an error occurred: %s", err)
		return
	}
}

func (org *Organization) ListPeers() []peer.ID {
	return org.ps.ListPeers(topicName(org.organizationName))
}

func (org *Organization) refreshPeers() {
	peerRefreshTicker := time.NewTicker(time.Second * 3)
	defer peerRefreshTicker.Stop()
	reuploadMapRouter := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-peerRefreshTicker.C:
			ids := org.ListPeers()
			MR.Update(ids)
		case <-reuploadMapRouter.C:
			MR.upload()
		}
	}
}

func (org *Organization) handleRecvMessage() {
	for {
		select {
		case msg := <-org.Messages:
			fmt.Println(msg)
			if msg.ContentType == "command" {
				command := new(messages.Command)
				commandBytes, _ := json.Marshal(msg.Content)
				json.Unmarshal(commandBytes, command)
				fmt.Println(command)
			}
		}
	}
}

func realHandleMsg(msg messages.Message) {
	switch msg.ContentType {
	case "preprepare":

	case "prepare":

	case "commit":

	}
}
