package gossip
	
import (
    "log/slog"
    "math/rand/v2"
	"strings"
	"sync"
	"net"
	"fmt"
	"encoding/json"
	"time"
)

type gossipImp struct {
	selfID string
	nodesList NodesList
	transport GossipCommunicationService
	m sync.RWMutex
}

type Gossip interface {
	StartListener(port string)
	StartGossiping(interval time.Duration)
	UpdateNodesList(nodesList NodesList)
	SetSeedNodes(seeds string)
	SendNodesList(nodesList NodesList)
	UpdateHeartbeatCounter()
}

func NewGossip(selfID string, transport GossipCommunicationService) Gossip {
	// Initialize the nodes list with the self node
	// TODO: agregar el transport address del nodo actual al nodes list
	//FIX EL LOCALHOST ESTA QUEMADO
	//FIX EL LOCALHOST ESTA QUEMADO
	//FIX EL LOCALHOST ESTA QUEMADO
	//FIX EL LOCALHOST ESTA QUEMADO
	//FIX EL LOCALHOST ESTA QUEMADO
	//FIX EL LOCALHOST ESTA QUEMADO

	NodesMap := make(map[string]NodesAttr)
	NodesMap[selfID] = NodesAttr{
		HeartbeatCounter: 0,
		Timestamp:        time.Now().Unix(),
		TransportAddress: selfID,
	}
	return &gossipImp{
		selfID: selfID,
		nodesList: NodesMap,
		transport: transport,
	}
}

func (g *gossipImp) StartListener(port string) {
	// Set up UDP address
    addr, err := net.ResolveUDPAddr("udp", "localhost:"+port)
    if err != nil {
        slog.Error("Couldnt resolve address", "error", err)
		return
    }
	slog.Info("Starting gossip listener on port", "port", port)

    // Start listening
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        slog.Error("Listen failed:", "error", err)
        return
    }
    defer conn.Close()

	// Listen for incoming messages
	
    // Buffer for incoming data
	// each node from the map uses ~90 bytes, so the buffer can hold up to 10 nodes
	// TODO: implement a more robust way to handle incoming messages, maybe using a channel and a goroutine to process the messages
    buffer := make([]byte, 1024)
    for {
        // Read client message
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            slog.Info("Read error", "error", err)
            continue
        }
        fmt.Printf("Got message from %s: %s\n", clientAddr, string(buffer[:n]))

		var nodesList NodesList
		err = json.Unmarshal(buffer[:n], &nodesList)
		if err != nil {
			slog.Error("Error unmarshalling nodes list", "error", err)
			continue
		}
		g.UpdateNodesList(nodesList)
    }
}

func (g *gossipImp) StartGossiping(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		g.UpdateHeartbeatCounter()
		g.SendNodesList(g.nodesList)
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

func(g *gossipImp) SetSeedNodes(seeds string) {
	if seeds == "" {
		return
	}

	g.m.Lock()
	defer g.m.Unlock()

	seedsList := strings.Split(seeds, ",")
	for _, seed := range seedsList {
		g.nodesList[seed] = NodesAttr{
			HeartbeatCounter: 0,
			Timestamp:        time.Now().Unix(),
			TransportAddress: seed,
		}
	}

}

func (g *gossipImp) SendNodesList(nodesList NodesList) {

	nodes := []string{}
	for key := range nodesList {
		if key == g.selfID {
			continue
		}
		actualNode := g.nodesList[key]
		nodes = append(nodes, actualNode.TransportAddress)
	}

	 if len(nodes) == 0 {
      // no hay peers, no enviar nada
      return
	}

	//RANDOM NODE MOVE TO SERVICE 
	// generate random number between 0 and len(nodes)-1

	randomIndex := rand.IntN(len(nodes))
	randomNodeAddrs := nodes[randomIndex]

	slog.Info("Sending nodes list to node", "address", randomNodeAddrs)


	if err := g.transport.SendNodesList(randomNodeAddrs, nodesList); err != nil {
		slog.Error("Error sending nodes list", "error", err)
		return
	}
	slog.Info("Nodes list sent to node", "address", randomNodeAddrs)
}

func (g *gossipImp) UpdateHeartbeatCounter() {
	g.m.Lock()
	defer g.m.Unlock()

	actualNode := g.nodesList[g.selfID]
	actualNode.HeartbeatCounter++
	actualNode.Timestamp = time.Now().Unix()
	g.nodesList[g.selfID] = actualNode
}