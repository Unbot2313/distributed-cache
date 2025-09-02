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
	AddNode(Id string)
	GetNode(Id string) *Node
	DeleteNode(id string) error
}

func NewRing(Hasher hash.Hasher, VirtualNodes int) *consistentHash {
	return &consistentHash{
		Nodes:      Nodes{},
		Hasher:    Hasher,
		virtualNodes:  VirtualNodes,
	}
}

func (r *consistentHash) AddNode(id string) {

	// por cada iteracion se guarda un nodo virtual con un hash diferente debido a su virtual id
	for i := 0; i < r.virtualNodes; i++ {
		Node := NewVirtualNode(r.Hasher, id, i)
		r.Nodes = append(r.Nodes, Node)
	}

	// sort para busqueda binaria al momento de ingresar un registro y buscar el nodo cercano
	sort.Sort(r.Nodes)
}

