package ring

import (
	"fmt"

	"github.com/unbot2313/distributed-cache/pkg/hash"
)

type Nodes []*Node

type Node struct {
    PhysicalId string  // "server-1" 
    VirtualId  string  // "server-1:0", "server-1:1", etc
    HashId     uint32  // hash("server-1:0")
}

func NewVirtualNode(hasher hash.Hasher, physicalId string, virtualIndex int) *Node {
    virtualId := fmt.Sprintf("%s:%d", physicalId, virtualIndex)
    return &Node{
        PhysicalId: physicalId,
        VirtualId:  virtualId,
        HashId:     hasher.Hash(virtualId),
    }
}