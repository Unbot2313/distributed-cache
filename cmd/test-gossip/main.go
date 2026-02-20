package main

import (
	"os"
	"time"
	"github.com/unbot2313/distributed-cache/pkg/gossip"
)

func main() {
	addrUDPArg := os.Args[2]

	// gossip communication service
	gossipService := gossip.NewGossipCommunicationService()
	gossip := gossip.NewGossip(addrUDPArg, gossipService)
	go gossip.StartListener(addrUDPArg)
	gossip.StartGossiping(5 * time.Second)
}