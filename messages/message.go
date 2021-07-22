package messages

import "fmt"

const ChanMsgSize = 1024

type Message struct {
	SenderID    string
	SenderName  string
	Timestamp   int64
	ContentType string
	Content     interface{}
}

func (m Message) String() string {
	return fmt.Sprintf("<Message>{SenderID: %s, SenderName: %s, ContentType: %s, Timestamp: %d", m.SenderID, m.SenderName, m.ContentType, m.Timestamp)
}
