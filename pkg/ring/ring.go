package ring

import (
	"sort"

	"github.com/unbot2313/distributed-cache/pkg/hash"
)

// en lugar de nodos se podria implementar un arbol binario pero no encontre diferencias en su Big(O)
type consistentHash struct {
	Nodes Nodes
	Hasher hash.Hasher
	virtualNodes int // Number of virtual nodes per physical node
}

type Ring interface {
	AddNode(Id string)
	GetNode(Id string) *Node
	DeleteNode(id string) error
	redistribute()
}

func NewRing(Hasher hash.Hasher, VirtualNodes int) Ring {
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

func (r *consistentHash) GetNode(key string) *Node {
	hash := r.Hasher.Hash(key)
	// sort devuelve el valor maximo posible osea 2**32 que es el numero maximo del hash en caso de no encontrar nada
	idx := sort.Search(len(r.Nodes), func(i int) bool {
		return r.Nodes[i].HashId >= hash
	})
	// si no hay un nodo con valor superior para seguir el orden de reloj voy al primero
	if idx == len(r.Nodes) {
		idx = 0
	}

	return r.Nodes[idx]
}

func (r *consistentHash) DeleteNode(id string) error {

	// funcion de redistribucion
	

	
	return nil
}

func (r *consistentHash) redistribute() {}
