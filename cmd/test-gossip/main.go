package main

import (
	"time"
	"flag"
	"github.com/unbot2313/distributed-cache/pkg/gossip"
)

func main() {
	// Command-line arguments
	gossipPortArg := flag.String("gossip", "7000", "UDP port to listen on for gossip communication")
	// httpPortArg := flag.String("http", "3211", "HTTP port to listen on")
	seedsArg := flag.String("seeds", "", "Comma-separated list of seed nodes in the format <ip>:<port>")

	flag.Parse()

	// gossip communication service
	gossipService := gossip.NewGossipCommunicationService()
	gossip := gossip.NewGossip("localhost:" + *gossipPortArg, gossipService)
	gossip.SetSeedNodes(*seedsArg)
	go gossip.StartListener(*gossipPortArg)
	gossip.StartGossiping(5 * time.Second)
}