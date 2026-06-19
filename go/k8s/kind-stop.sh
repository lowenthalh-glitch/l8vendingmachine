#!/usr/bin/env bash
set -e

CLUSTER_NAME="vend"

echo "Deleting KIND cluster '$CLUSTER_NAME'..."
kind delete cluster --name "$CLUSTER_NAME"
echo "KIND cluster '$CLUSTER_NAME' deleted."
