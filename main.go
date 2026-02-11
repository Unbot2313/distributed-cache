package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/unbot2313/distributed-cache/pkg/handlers"
	"github.com/unbot2313/distributed-cache/pkg/hash"
	"github.com/unbot2313/distributed-cache/pkg/ring"
	"github.com/unbot2313/distributed-cache/pkg/router"
	"github.com/unbot2313/distributed-cache/pkg/services"
)

func main() {
	hasher := hash.NewXXH3Hasher()

	// aprox con 100 nodos tendra una desviacion de 10% como maximo, con 200 un 5%
	r := ring.NewRing(hasher, 100)

	// dragonfly-0 -> localhost:6379
	// dragonfly-1 -> localhost:6380
	// dragonfly-2 -> localhost:6381
	nodes := []router.NodeConfig{
		{
			ID: "dragonfly-0",
			CacheConfig: &services.CacheConfig{
				Addr:         "localhost:6379",
				Password:     "",
				DB:           0,
				MaxRetries:   3,
				Timeout:      5 * time.Second,
				PoolSize:     10,
				MinIdleConns: 2,
			},
		},
		{
			ID: "dragonfly-1",
			CacheConfig: &services.CacheConfig{
				Addr:         "localhost:6380",
				Password:     "",
				DB:           0,
				MaxRetries:   3,
				Timeout:      5 * time.Second,
				PoolSize:     10,
				MinIdleConns: 2,
			},
		},
		{
			ID: "dragonfly-2",
			CacheConfig: &services.CacheConfig{
				Addr:         "localhost:6381",
				Password:     "",
				DB:           0,
				MaxRetries:   3,
				Timeout:      5 * time.Second,
				PoolSize:     10,
				MinIdleConns: 2,
			},
		},
	}

	cr, err := router.NewCacheRouter(r, nodes)
	if err != nil {
		log.Fatalf("failed to create cache router: %v", err)
	}
	defer cr.Close()

	// Registrar handlers HTTP en el mux
	mux := http.NewServeMux()
	handlers.RegisterCacheHandlers(mux, cr)

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\nShutting down...")
		cr.Close()
		os.Exit(0)
	}()

	// Iniciar servidor HTTP en :3211 con el mux que contiene los handlers
	addr := ":3211"
	fmt.Printf("Distributed Cache server listening on %s\n", addr)
	fmt.Println("Nodes: dragonfly-0 (:6379), dragonfly-1 (:6380), dragonfly-2 (:6381)")

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
