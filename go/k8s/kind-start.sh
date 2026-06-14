#!/usr/bin/env bash
set -e

CLUSTER_NAME="vend"

# Check if kind is installed
if ! command -v kind &> /dev/null; then
    echo "kind not found, installing..."
    go install sigs.k8s.io/kind@latest
fi

# Create cluster with 1 control-plane + 1 worker
echo "Creating KIND cluster '$CLUSTER_NAME'..."
cat <<EOF | kind create cluster --name "$CLUSTER_NAME" --config=-
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
  - role: control-plane
  - role: worker
EOF

echo "Loading images into KIND cluster..."
IMAGES=(
    "saichler/vendmachine-vnet:latest"
    "saichler/vend-logs-vnet:latest"
    "saichler/vendmachine:latest"
    "saichler/vendmachine-web:latest"
    "saichler/vendmachine-inv:latest"
    "saichler/vendmachine-collector:latest"
    "saichler/vendmachine-parser:latest"
    "saichler/vend-log-agent:latest"
)

for img in "${IMAGES[@]}"; do
    echo "  Loading $img..."
    kind load docker-image "$img" --name "$CLUSTER_NAME" 2>/dev/null || echo "  Warning: $img not found locally, skipping"
done

echo "Deploying vend-kind.yaml..."

# Phase 1: Namespace + vnet + logs-vnet
kubectl apply -f vend-kind.yaml
echo "Waiting for vnet rollout..."
kubectl rollout status statefulset/vend-vnet -n vend-vnet --timeout=120s
kubectl rollout status statefulset/vend-logs-vnet -n vend-logs-vnet --timeout=120s

# Phase 2: Core services
echo "Waiting for core services rollout..."
kubectl rollout status statefulset/vend -n vend --timeout=120s
kubectl rollout status statefulset/vend-web -n vend-web --timeout=120s
kubectl rollout status statefulset/vend-inv -n vend-inv --timeout=120s
kubectl rollout status statefulset/vend-collector -n vend-collector --timeout=120s
kubectl rollout status statefulset/vend-parser -n vend-parser --timeout=120s
kubectl rollout status statefulset/vend-log-agent -n vend-log-agent --timeout=120s

echo ""
echo "KIND cluster '$CLUSTER_NAME' is ready!"
echo "To access: kubectl cluster-info --context kind-$CLUSTER_NAME"
