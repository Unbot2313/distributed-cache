package gossip

type gossipCommunicationServiceImp struct {}

type GossipCommunicationService interface {
	SendNodesList(transportAddress string, nodesList NodesList)
}

func NewGossipCommunicationService() GossipCommunicationService {
	return &gossipCommunicationServiceImp{}
}

func (g *gossipCommunicationServiceImp) SendNodesList(transportAddress string, nodesList NodesList) {
	//TODO: Implement the logic to send the nodes list to other nodes in the cluster
}
