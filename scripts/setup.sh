#!/bin/bash
set -e

echo "=== Creating k3d cluster ==="
if k3d cluster list | grep -q "distributed-cache"; then
    echo "Cluster 'distributed-cache' already exists, skipping creation"
else
    k3d cluster create distributed-cache --servers 1 --agents 1
fi

echo "=== Deploying DragonFly ==="
kubectl apply -f k8s/dragonfly-headless-service.yaml
kubectl apply -f k8s/dragonfly-statefulset.yaml

echo "=== Waiting for pods to be ready ==="
kubectl rollout status statefulset/dragonfly --timeout=120s

echo "=== Setting up port-forwards ==="
kubectl port-forward pod/dragonfly-0 6379:6379 &
kubectl port-forward pod/dragonfly-1 6380:6379 &
kubectl port-forward pod/dragonfly-2 6381:6379 &

sleep 2

echo "=== Verifying connectivity ==="
redis-cli -p 6379 ping
redis-cli -p 6380 ping
redis-cli -p 6381 ping

echo ""
echo "=== Ready! Run: go run main.go ==="
