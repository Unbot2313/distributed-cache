package ring

import (
	"github.com/unbot2313/distributed-cache/pkg/hash"
)

type consistentHash struct {
	Nodes Nodes
	Hasher hash.Hasher
	virtualNodes int // Number of virtual nodes per physical node
}

type Ring interface {
	AddNode(Id int)
	GetNode(Id int) *Node
	DeleteNode(id int) error
}

func NewRing(Hasher hash.Hasher, VirtualNodes int) *consistentHash {
	return &consistentHash{
		Nodes:      Nodes{},
		Hasher:    Hasher,
		virtualNodes:  VirtualNodes,
	}
}

// func (r *Ring) AddNode(id string) {
// 	node := NewNode(r.Hasher, id)
// 	r.Nodes = append(r.Nodes, *node)
// 	sort.Sort(r.Nodes)
// }