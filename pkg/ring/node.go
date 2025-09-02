package ring

import (
	"github.com/unbot2313/distributed-cache/pkg/hash"
)

type Nodes []*Node

type Node struct {
	Id     int
	HashId uint32
}

func NewNode(Hasher hash.Hasher, id int) *Node{
	  return &Node{
		Id:     id,
		HashId: Hasher.Hash(id),
	}
}