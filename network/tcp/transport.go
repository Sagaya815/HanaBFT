// 通过这个包传输的数据有Command和Reply
package tcp

import (
	"encoding/gob"
	"fmt"
	"hanaBFT/hlog"
	"hanaBFT/messages"
	"net"
	"net/url"
)

type Transport struct {
	send chan messages.Message
	recv chan messages.Message
	uri  *url.URL
}

func init() {
	gob.Register(messages.Message{})
	gob.Register(messages.Command{})
}

func NewTransport(address string) *Transport {
	// address: tcp://192.168.10.127:8080
	uri, err := url.Parse(address)
	if err != nil {
		hlog.Fatalf("When I parse the url address, an error occurred: %s", err)
	}
	transport := &Transport{
		send: make(chan messages.Message, messages.ChanMsgSize),
		recv: make(chan messages.Message, messages.ChanMsgSize),
		uri:  uri,
	}
	return transport
}

func (t *Transport) Listen() {
	hlog.Debugf("Start listening port: %s", t.uri.Port())
	listener, err := net.Listen("tcp", ":"+t.uri.Port())
	if err != nil {
		hlog.Fatalf("When I listen port, an error occurred: %s", err)
	}
	go func(listener net.Listener) {
		defer listener.Close()
		for {
			conn, err := listener.Accept()
			if err != nil {
				hlog.Errorf("Get an new conn failed, error is: %s", err)
				continue
			}
			go func(conn net.Conn) {
				decoder := gob.NewDecoder(conn)
				defer conn.Close()
				for {
					var m messages.Message
					err := decoder.Decode(&m)
					if err != nil {
						hlog.Errorf("Decode the coming message failed, error is: %s", err)
						continue
					}
					t.recv <- m
				}
			}(conn)
		}
	}(listener)
}

func (t *Transport) Dial() error {
	conn, err := net.Dial("tcp", t.uri.Host)
	if err != nil {
		return err
	}
	hlog.Debugf("Successfully dial to %", t.uri.Host)
	go func(conn net.Conn) {
		encoder := gob.NewEncoder(conn)
		defer conn.Close()
		for m := range t.send {
			err := encoder.Encode(&m)
			if err != nil {
				hlog.Errorf("Encode the message failed, error is: %s", err)
			}
		}
	}(conn)
	return nil
}

func (t *Transport) Send(msg messages.Message) {
	t.send <- msg
}

func (t *Transport) HandleRecvMsg() {
	for {
		select {
		case msg := <-t.recv:
			fmt.Println(msg)
			fmt.Println(msg.Content)
		}
	}
}
