package gossip
	
import (
    "fmt"
    "math/rand/v2"
	 "sync"
)

type gossipImp struct {
	selfID string
	nodesList NodesList
	transport GossipCommunicationService
	m sync.RWMutex
}

type Gossip interface {
	UpdateNodesList(nodesList NodesList)
	SendNodesList(transportAddress string, nodesList NodesList)
}

func NewGossip(selfID string, transport GossipCommunicationService) Gossip {
	return &gossipImp{
		selfID: selfID,
		nodesList: make(map[string]NodesAttr),
		transport: transport,
	}
}

func (g *gossipImp) UpdateNodesList(nodesList NodesList) {

	g.m.Lock()
	defer g.m.Unlock()

	for key := range nodesList {
		actualNode := g.nodesList[key]
		if actualNode.HeartbeatCounter < nodesList[key].HeartbeatCounter { 
			g.nodesList[key] = nodesList[key]
		}
	}
}

func (g *gossipImp) SendNodesList(transportAddress string, nodesList NodesList) {

	nodes := []string{}
	for key := range nodesList {
		if key == g.selfID {
			continue
		}
		actualNode := g.nodesList[key]
		nodes = append(nodes, actualNode.TransportAddress)
	}

	//RANDOM NODE MOVE TO SERVICE 
	// generate random number between 0 and len(nodes)-1

	randomIndex := rand.IntN(len(nodes))
	randomNodeAddrs := nodes[randomIndex]

	fmt.Printf("Sending nodes list to node: %s\n", randomNodeAddrs)


	g.transport.SendNodesList(randomNodeAddrs, nodesList)
}