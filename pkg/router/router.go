package router

import (
	"context"
	"fmt"
	"sync"

	"github.com/unbot2313/distributed-cache/pkg/ring"
	"github.com/unbot2313/distributed-cache/pkg/services"
)

// NodeConfig holds the configuration for a single cache node.
type NodeConfig struct {
	ID          string              // PhysicalId used in the ring (e.g., "dragonfly-0")
	CacheConfig *services.CacheConfig
}

// CacheRouter routes cache operations to the correct DragonFly instance
// based on the consistent hash ring.
type CacheRouter struct {
	ring     ring.Ring
	services map[string]services.CacheService // PhysicalId -> CacheService
	mu       sync.RWMutex
}

// NewCacheRouter creates a CacheRouter from a Ring and a slice of NodeConfigs.
// It creates one CacheService per node, adds each node to the ring,
// and maps PhysicalId -> CacheService.
func NewCacheRouter(r ring.Ring, nodes []NodeConfig) (*CacheRouter, error) {
	svcMap := make(map[string]services.CacheService, len(nodes))

	for _, n := range nodes {
		client := services.CreateCacheClient(n.CacheConfig)
		svc := services.NewCacheService(client)
		svcMap[n.ID] = svc
		r.AddNode(n.ID)
	}

	return &CacheRouter{
		ring:     r,
		services: svcMap,
	}, nil
}

// resolve returns the CacheService responsible for the given cache key.
func (cr *CacheRouter) resolve(key string) (services.CacheService, string, error) {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	node := cr.ring.GetNode(key)
	if node == nil {
		return nil, "", fmt.Errorf("no node found for key: %s", key)
	}

	svc, ok := cr.services[node.PhysicalId]
	if !ok {
		return nil, "", fmt.Errorf("no cache service for node: %s", node.PhysicalId)
	}

	return svc, node.PhysicalId, nil
}

// Get retrieves a value by key from the appropriate node.
func (cr *CacheRouter) Get(ctx context.Context, key string) (string, bool, error) {
	svc, _, err := cr.resolve(key)
	if err != nil {
		return "", false, err
	}
	return svc.GetKey(ctx, key)
}

// Set stores a key-value pair on the appropriate node.
func (cr *CacheRouter) Set(ctx context.Context, key string, value string) error {
	svc, _, err := cr.resolve(key)
	if err != nil {
		return err
	}
	return svc.SetKey(ctx, key, value)
}

// Delete removes a key from the appropriate node.
func (cr *CacheRouter) Delete(ctx context.Context, key string) error {
	svc, _, err := cr.resolve(key)
	if err != nil {
		return err
	}
	return svc.DeleteKey(ctx, key)
}

// GetNodeForKey returns which physical node a key would be routed to.
func (cr *CacheRouter) GetNodeForKey(key string) (string, error) {
	_, nodeID, err := cr.resolve(key)
	return nodeID, err
}

// PingAll pings every cache node and returns a map of nodeID -> error (nil if healthy).
func (cr *CacheRouter) PingAll(ctx context.Context) map[string]error {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	results := make(map[string]error, len(cr.services))
	for id, svc := range cr.services {
		results[id] = svc.Ping(ctx)
	}
	return results
}

// GetNodeIDs returns the list of physical node IDs registered in the router.
func (cr *CacheRouter) GetNodeIDs() []string {
	cr.mu.RLock()
	defer cr.mu.RUnlock()

	ids := make([]string, 0, len(cr.services))
	for id := range cr.services {
		ids = append(ids, id)
	}
	return ids
}

// Close closes all underlying CacheService connections.
func (cr *CacheRouter) Close() error {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	var firstErr error
	for id, svc := range cr.services {
		if err := svc.Close(); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("failed to close service %s: %w", id, err)
		}
	}
	return firstErr
}
