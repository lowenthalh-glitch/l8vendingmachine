#!/usr/bin/env bash
set -e

echo "Building all L8 Vending Machine images..."

echo "1/8 Building vnet..."
cd vend/vnet && ./build.sh && cd ../..

echo "2/8 Building logs-vnet..."
cd vend/log-vnet && ./build.sh && cd ../..

echo "3/8 Building log-agent..."
cd vend/log-agent && ./build.sh && cd ../..

echo "4/8 Building main service..."
cd vend/main && ./build.sh && cd ../..

echo "5/8 Building web UI..."
cd vend/ui/main && ./build.sh && cd ../../..

echo "6/8 Building collector..."
cd vend/collector && ./build.sh && cd ../..

echo "7/8 Building parser..."
cd vend/parser && ./build.sh && cd ../..

echo "8/8 Building inventory cache..."
cd vend/inv_vend && ./build.sh && cd ../..

echo "All images built successfully!"
