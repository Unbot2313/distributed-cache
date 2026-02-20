package gossip

type NodesAttr struct {
	HeartbeatCounter uint64
	Timestamp        int64
	TransportAddress string
}

type NodesList map[string]NodesAttr