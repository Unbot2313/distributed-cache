package ring

import (
	"sort"

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

func (r *consistentHash) AddNode(id int) {
	// se añade la cantidad de nodos segund los nodos virtuales
	for i := 0; i < r.virtualNodes; i++ {
		Node := NewNode(r.Hasher, id)
		r.Nodes = append(r.Nodes, Node)
	}
	sort.Sort(r.Nodes)
}

