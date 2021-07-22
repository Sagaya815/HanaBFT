package utils

import (
	"github.com/libp2p/go-libp2p-core/peer"
	"time"
)

func ShortPeerID(id peer.ID) string {
	return id.Pretty()[len(id.Pretty())-8:]
}

func Retry(f func() error, maxTimes int) {
	for i := 0; i < maxTimes; i++ {
		err := f()
		if err == nil {
			break
		}
		time.Sleep(time.Duration(1000 * (i + 1)))
	}
}
