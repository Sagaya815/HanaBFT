package messages

import (
	"encoding/gob"
	"fmt"
)

type Command struct {
	SenderID   string
	SenderName string
	CommandID  int
	Key        key
	Value      value
}

func (c Command) String() string {
	return fmt.Sprintf("<Command>{SenderID: %s, SenderName: %s, CommandID: %d, key: %v, value: %v",
		c.SenderID, c.SenderName, c.CommandID, c.Key, c.Value)
}

func init() {
	gob.Register(Command{})
}
