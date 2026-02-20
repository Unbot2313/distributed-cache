package gossip

import (
	"encoding/json"
	"log/slog"
	"net"
)

type gossipCommunicationServiceImp struct {}

type GossipCommunicationService interface {
	SendNodesList(transportAddress string, nodesList NodesList) (err error)
}

func NewGossipCommunicationService() GossipCommunicationService {
	return &gossipCommunicationServiceImp{}
}

func (g *gossipCommunicationServiceImp) SendNodesList(transportAddress string, nodesList NodesList) (err error) {
	//TODO: Implement the logic to send the nodes list to other nodes in the cluster

	// resolve with tcp to the transportAddress and send the nodesList

	addr, err := net.ResolveUDPAddr("udp", transportAddress)

	if err != nil {
		slog.Error("Error resolving address", "error", err)
		return err
	}

	// connect with udp to the transportAddress and send the nodesList
	conn, err := net.DialUDP("udp", nil, addr)
    if err != nil {
        slog.Error("Connection failed", "error", err)
        return err
    }

	// close the connection after sending the nodesList
    defer conn.Close()


	// parse nodelist to bytes json

	byteBody, err := json.Marshal(nodesList)
	if err != nil {
		slog.Error("Error marshalling nodes list", "error", err)
		return err
	}

	// send the byteBody to the transportAddress using tcp

	_, err = conn.Write(byteBody)
    if err != nil {
        slog.Error("Send failed", "error", err)
        return err
    }

	// no response is expected, so we can return nil

	return nil

}
