package main

import (
	"hanaBFT/peer"
	"time"
)

func main() {
	peer.Setup()
	//go node.Run()
	//node.AsReplica()

	time.Sleep(time.Hour)
}
