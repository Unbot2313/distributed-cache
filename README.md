# 🚀 Distributed Cache

> A high-performance, horizontally-scalable distributed cache system built with Go, designed to compete with Redis Cluster while focusing on learning distributed systems concepts.

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)
![Build Status](https://img.shields.io/badge/Build-Passing-success?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)

## 🎯 Project Vision

**Distributed Cache** is a learning-focused project that implements a production-ready distributed caching system from scratch. The goal is to understand and implement core distributed systems concepts while building something that can compete with industry standards like Redis Cluster.

### Why This Project?

- 📚 **Learning Distributed Systems**: Hands-on experience with consensus algorithms, fault tolerance, and distributed computing
- ⚡ **Performance-Focused**: Built with speed as a primary concern, leveraging Go's concurrency model
- 🔄 **Modern Architecture**: Uses cutting-edge tools like DragonFly as the storage backend
- 🛠️ **Production-Ready**: Designed with real-world scalability and reliability requirements

---

## ✨ Key Features

### 🔥 Core Functionality
- **Consistent Hashing**: Efficient key distribution with virtual nodes for optimal load balancing
- **Horizontal Scaling**: Add/remove nodes seamlessly without service interruption  
- **High Performance**: Sub-millisecond response times with xxHash3 for lightning-fast hashing
- **Fault Tolerance**: Automatic failover and data replication across nodes

### 🌐 Distributed Systems Concepts
- **Gossip Protocol**: Efficient node discovery and cluster membership management
- **Quorum Consensus**: Configurable consistency levels (R + W > N) for different use cases
- **Automatic Rebalancing**: Smart data migration when cluster topology changes
- **Health Monitoring**: Real-time cluster health and performance metrics

### 🏗️ Architecture
- **Pluggable Storage**: Support for DragonFly, Redis, or any cache backend
- **Smart Client**: Client-side routing to minimize proxy overhead
- **Kubernetes Ready**: Native K8s support with StatefulSets and service discovery
- **Observability**: Built-in Prometheus metrics and distributed tracing

---

## 🏛️ Architecture Overview

```
                    ┌─────────────────┐
                    │   Load Balancer │
                    └─────────┬───────┘
                              │
                    ┌─────────▼───────┐
                    │  Proxy/Router   │
                    │ (Consistent     │
                    │  Hashing)       │
                    └─────────┬───────┘
                              │
         ┌────────────────────┼────────────────────┐
         │                    │                    │
    ┌────▼────┐         ┌────▼────┐         ┌────▼────┐
    │ Node 1  │◄────────┤ Node 2  │────────►│ Node 3  │
    │DragonFly│  Gossip │DragonFly│  Gossip │DragonFly│
    └─────────┘         └─────────┘         └─────────┘
         │                    │                    │
         └────────────────────┼────────────────────┘
                         Replication
```

### Hash Ring Distribution
```
                    Hash Space (0 - 2³²)
    ┌─────────────────────────────────────────────────┐
    │  [N1] ──── [N3] ── [N2] ──── [N1] ── [N3] ──   │
    │    ▲         ▲       ▲         ▲       ▲        │
    │ Virtual   Virtual Virtual   Virtual Virtual     │
    │  Nodes     Nodes   Nodes     Nodes   Nodes     │
    └─────────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

### Prerequisites
- **Go 1.21+**
- **Docker & Docker Compose**
- **DragonFly** or **Redis** (for storage backend)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/distributed-cache.git
cd distributed-cache

# Install dependencies
go mod download

# Run tests
go test ./...

# Start a 5-node cluster
docker-compose up
```

### Basic Usage

```go
package main

import (
    "github.com/yourusername/distributed-cache/pkg/client"
    "github.com/yourusername/distributed-cache/pkg/hash"
    "github.com/yourusername/distributed-cache/pkg/ring"
)

func main() {
    // Create a hash ring
    hasher := hash.NewXXH3Hasher()
    hashRing := ring.NewRing(hasher, 150) // 150 virtual nodes per physical node
    
    // Add cache nodes
    hashRing.AddNode("cache-node-1:8080")
    hashRing.AddNode("cache-node-2:8080") 
    hashRing.AddNode("cache-node-3:8080")
    
    // Route keys to appropriate nodes
    node := hashRing.GetNode("user:12345")
    fmt.Printf("Key 'user:12345' routes to: %s\n", node)
    
    // Create cache client
    client := client.New(hashRing)
    
    // Cache operations
    client.Set("user:12345", userData)
    data, err := client.Get("user:12345")
    client.Delete("user:12345")
}
```

---

## 📂 Project Structure

```
distributed-cache/
├── cmd/
│   ├── cache-node/          # Cache server executable
│   ├── cache-proxy/         # Proxy/router service  
│   └── benchmark/           # Performance testing tools
├── pkg/
│   ├── hash/                # Hash functions (xxHash3)
│   ├── ring/                # Consistent hashing implementation
│   ├── gossip/              # Gossip protocol for node discovery
│   ├── storage/             # Storage backends (DragonFly, Redis)
│   ├── quorum/              # Quorum consensus implementation  
│   ├── client/              # Smart client library
│   └── monitoring/          # Metrics and observability
├── docker/
│   ├── docker-compose.yml   # 5-node development cluster
│   └── Dockerfile           # Container images
├── k8s/                     # Kubernetes manifests
│   ├── statefulset.yml      # Cache nodes
│   ├── service.yml          # Service discovery
│   └── configmap.yml        # Configuration
├── monitoring/
│   ├── prometheus.yml       # Metrics collection
│   ├── grafana/             # Dashboards
│   └── alerts.yml           # Alerting rules
└── benchmarks/              # Performance tests and comparisons
    ├── redis-comparison/    # Benchmark vs Redis Cluster  
    └── load-tests/          # Stress testing scenarios
```

---

## 🎯 Implementation Roadmap

### ✅ Phase 1: Core Foundation
- [x] Consistent Hashing with Virtual Nodes
- [x] xxHash3 Implementation  
- [x] Basic Ring Operations (Add/Remove/Get)
- [x] Unit Tests & Benchmarks

### 🚧 Phase 2: Storage Integration (In Progress)
- [ ] DragonFly Backend Integration
- [ ] Redis Backend Support
- [ ] Storage Interface Abstraction
- [ ] Connection Pooling & Management

### 📋 Phase 3: Distributed Coordination  
- [ ] Gossip Protocol Implementation
- [ ] Node Discovery & Health Monitoring
- [ ] Automatic Failure Detection
- [ ] Cluster Membership Management

### 📋 Phase 4: Consensus & Consistency
- [ ] Quorum-based Operations (R + W > N)
- [ ] Configurable Consistency Levels
- [ ] Read Repair Mechanisms
- [ ] Anti-Entropy Protocols

### 📋 Phase 5: Production Features
- [ ] Data Migration & Rebalancing
- [ ] Prometheus Metrics Export
- [ ] Kubernetes Native Deployment
- [ ] Performance Optimization

### 📋 Phase 6: Advanced Features
- [ ] Multi-Datacenter Replication
- [ ] Compression & Serialization Options
- [ ] TTL & Expiration Policies  
- [ ] Admin Dashboard & Monitoring

---

## ⚡ Performance Goals

### Target Benchmarks
| Operation | Target Latency | Target Throughput |
|-----------|---------------|-------------------|
| GET       | < 1ms P99     | > 100K ops/sec    |
| SET       | < 2ms P99     | > 80K ops/sec     |
| DELETE    | < 1ms P99     | > 90K ops/sec     |

### vs Redis Cluster Comparison
- **Latency**: Aim for comparable P95/P99 latencies
- **Throughput**: Target 80%+ of Redis Cluster performance  
- **Memory Efficiency**: Competitive memory usage per key
- **Network Overhead**: Minimize inter-node communication

---

## 🧪 Testing & Development

### Running Tests
```bash
# Unit tests
go test ./pkg/...

# Integration tests
go test ./integration/...

# Benchmarks
go test -bench=. ./pkg/hash/
go test -bench=. ./pkg/ring/

# Race condition detection
go test -race ./...
```

### Development Setup
```bash
# Start development cluster
make dev-cluster

# Run with hot reload
make dev-server

# Generate test data
make test-data

# Performance profiling
make profile
```

---

## 📊 Monitoring & Observability

### Metrics Exposed
- **Cache Hit Ratio**: Percentage of successful cache retrievals
- **Request Latency**: P50, P95, P99 response times by operation
- **Node Health**: CPU, Memory, Network usage per node
- **Ring Distribution**: Key distribution balance across nodes
- **Gossip Activity**: Membership changes and failure detection

### Grafana Dashboards
- **Cluster Overview**: High-level system health and performance
- **Node Details**: Per-node metrics and resource usage
- **Cache Performance**: Hit ratios, latency distributions
- **Network Activity**: Inter-node communication patterns

---

## 🤝 Contributing

This is primarily a learning project, but contributions are welcome! 

### Development Principles
- **Code Quality**: Comprehensive tests and clean, readable code
- **Performance First**: Every feature should be benchmarked
- **Learning Focus**: Document design decisions and trade-offs
- **Production Ready**: Build as if it's going to production

### Getting Started
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Add tests for your changes
4. Ensure benchmarks pass
5. Submit a pull request

---

## 📚 Learning Resources

### Distributed Systems Concepts
- [Consistent Hashing Paper](https://www.cs.princeton.edu/courses/archive/fall09/cos518/papers/chash.pdf)
- [Gossip Protocols Overview](https://en.wikipedia.org/wiki/Gossip_protocol)  
- [CAP Theorem Explained](https://en.wikipedia.org/wiki/CAP_theorem)

### Go Performance
- [Go Memory Model](https://golang.org/ref/mem)
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Concurrency Patterns](https://blog.golang.org/concurrency-patterns)

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- **Redis Team** for inspiration and architecture patterns
- **DragonFly Team** for the modern cache backend
- **Go Community** for excellent tooling and libraries
- **Distributed Systems Research** for foundational algorithms

---


**Built with ❤️ for learning distributed systems**

[Report Bug](https://github.com/yourusername/distributed-cache/issues) • [Request Feature](https://github.com/yourusername/distributed-cache/issues) • [Documentation](https://github.com/yourusername/distributed-cache/wiki)

