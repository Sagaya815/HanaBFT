package utils

import "github.com/libp2p/go-libp2p-core/peer"

func ShortPeerID(id peer.ID) string {
	return id.Pretty()[len(id.Pretty())-8:]
}
