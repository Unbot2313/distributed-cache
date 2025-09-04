package ring

import (
	"github.com/unbot2313/distributed-cache/pkg/hash"
	"github.com/unbot2313/distributed-cache/pkg/utils"
)

type Nodes []*Node

type Node struct {
    PhysicalId string  // "server-1" 
    VirtualId  string  // "server-1:0", "server-1:1", etc
    HashId     uint32  // hash("server-1:0")
}

func NewVirtualNode(hasher hash.Hasher, physicalId string, virtualIndex int) *Node {
    virtualId := utils.GetVirtualId(physicalId, virtualIndex)
    return &Node{
        PhysicalId: physicalId,
        VirtualId:  virtualId,
        HashId:     hasher.Hash(virtualId),
    }
}