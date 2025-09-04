package ring

import (
	"sort"

	"github.com/unbot2313/distributed-cache/pkg/hash"
)

// en lugar de nodos se podria implementar un arbol binario pero no encontre diferencias en su Big(O)
type consistentHash struct {
	nodes Nodes
	hasher hash.Hasher
	virtualnodes int // Number of virtual nodes per physical node
}

type Ring interface {
	AddNode(Id string)
	GetNode(Id string) *Node
	DeleteNode(id string) error
	getLength() int
}

func NewRing(Hasher hash.Hasher, Virtualnodes int) Ring {
	return &consistentHash{
		nodes:      Nodes{},
		hasher:    Hasher,
		virtualnodes:  Virtualnodes,
	}
}

func (r *consistentHash) getLength() int{
	return len(r.nodes)
}

func (r *consistentHash) AddNode(id string) {

	// por cada iteracion se guarda un nodo virtual con un hash diferente debido a su virtual id
	for i := 0; i < r.virtualnodes; i++ {
		Node := NewVirtualNode(r.hasher, id, i)
		r.nodes = append(r.nodes, Node)
	}

	// sort para busqueda binaria al momento de ingresar un registro y buscar el nodo cercano
	sort.Sort(r.nodes)
}

func (r *consistentHash) GetNode(key string) *Node {
	hash := r.hasher.Hash(key)
	// sort devuelve el valor maximo posible osea 2**32 que es el numero maximo del hash en caso de no encontrar nada
	idx := sort.Search(r.getLength(), func(i int) bool {
		return r.nodes[i].HashId >= hash
	})
	// si no hay un nodo con valor superior para seguir el orden de reloj voy al primero
	if idx == len(r.nodes) {
		idx = 0
	}

	return r.nodes[idx]
}

func (r *consistentHash) DeleteNode(id string) error {
    var newnodes Nodes
    
    // Filtro normal con complejidad O(n)
    for _, node := range r.nodes {
        if node.PhysicalId != id {
            newnodes = append(newnodes, node)
        }
    }
    
    r.nodes = newnodes
    return nil
}

// func (r *consistentHash) DeleteNode(id string) error {

// complejidad O(n*logN)
// 	for i := 0; i < r.virtualnodes; i++{
// 		virtualId := utils.GetVirtualId(id, i)
// 		hashedVirtualId := r.Hasher.Hash(virtualId)
// 		idx := sort.Search(r.getLength(), func(i int) bool {
// 		return r.nodes[i].HashId == hashedVirtualId
// 	})
// 		r.nodes = append(r.nodes[:idx], r.nodes[idx+1:]...)
// 	}
// 	return nil
// }
