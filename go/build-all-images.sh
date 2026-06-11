#!/usr/bin/env bash
set -e

echo "Building all L8 Vending Machine images..."

echo "1/6 Building vnet..."
cd vend/vnet && ./build.sh && cd ../..

echo "2/6 Building main service..."
cd vend/main && ./build.sh && cd ../..

echo "3/6 Building web UI..."
cd vend/ui/main && ./build.sh && cd ../../..

echo "4/6 Building collector..."
cd vend/collector && ./build.sh && cd ../..

echo "5/6 Building parser..."
cd vend/parser && ./build.sh && cd ../..

echo "6/6 Building inventory cache..."
cd vend/inv_vend && ./build.sh && cd ../..

echo "All images built successfully!"
