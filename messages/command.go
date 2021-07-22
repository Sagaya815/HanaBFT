package messages

import (
	"encoding/gob"
	"fmt"
)

type Command struct {
	ProposerID   string
	ProposerName string
	CommandID    int
	Key          key
	Value        value
}

func (c Command) String() string {
	return fmt.Sprintf("<Command>{ProposerID: %s, ProposerName: %s, CommandID: %d, key: %v, value: %v",
		c.ProposerID, c.ProposerName, c.CommandID, c.Key, c.Value)
}

func init() {
	gob.Register(Command{})
}
