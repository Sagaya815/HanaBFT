package messages

import (
	"encoding/gob"
	"fmt"
)

type Reply struct {
	SenderID  string
	CommandID int
}

func (r Reply) String() string {
	return fmt.Sprintf("<Reply>{SenderID: %s, CommandID: %d}", r.SenderID, r.CommandID)
}

func init() {
	gob.Register(Reply{})
}
