package messages

import (
	"encoding/gob"
	"fmt"
)

const ChanMsgSize = 1024

type Message struct {
	SenderID    string
	SenderName  string
	ContentType string
	Content     interface{}
}

func (m Message) String() string {
	return fmt.Sprintf("<Message>{SenderID: %s, SenderName: %s, ContentType: %s", m.SenderID, m.SenderName, m.ContentType)
}

func init() {
	gob.Register(Message{})
}
