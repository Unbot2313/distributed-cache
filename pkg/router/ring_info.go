package router

import "fmt"

// RingInfo holds distribution statistics for the ring.
type RingInfo struct {
	Nodes        []string       `json:"nodes"`
	TotalNodes   int            `json:"total_nodes"`
	Distribution map[string]int `json:"sample_distribution"`
	SampleSize   int            `json:"sample_size"`
}

// GetRingInfo samples N keys and counts how many route to each physical node.
func (cr *CacheRouter) GetRingInfo(sampleSize int) RingInfo {
	dist := make(map[string]int)
	nodeIDs := cr.GetNodeIDs()

	for _, id := range nodeIDs {
		dist[id] = 0
	}

	for i := 0; i < sampleSize; i++ {
		key := fmt.Sprintf("sample-key-%d", i)
		node := cr.ring.GetNode(key)
		if node != nil {
			dist[node.PhysicalId]++
		}
	}

	return RingInfo{
		Nodes:        nodeIDs,
		TotalNodes:   len(nodeIDs),
		Distribution: dist,
		SampleSize:   sampleSize,
	}
}
