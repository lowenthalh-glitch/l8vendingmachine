#!/usr/bin/env bash

set -e

wget https://raw.githubusercontent.com/saichler/l8types/refs/heads/main/proto/api.proto
wget https://raw.githubusercontent.com/saichler/l8common/refs/heads/main/proto/l8common.proto

# Use the protoc image to run protoc.sh and generate the bindings.
# Note: l8common.proto is downloaded for import resolution only — its Go code
# comes from the vendored l8common dependency, not from local generation.

# Vending Machine Management
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-common.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-fleet.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-inventory.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-sales.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-payment.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-temperature.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-maintenance.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-route.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-analytics.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-access.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-dex.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-warehouse.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-dashboard.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-compliance.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-reports.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest
docker run --user "$(id -u):$(id -g)" -e PROTO=vend-retention.proto --mount type=bind,source="$PWD",target=/home/proto/ -i saichler/protoc:latest

rm api.proto l8common.proto

# Now move the generated bindings to the models directory and clean up
rm -rf ../go/types
mkdir -p ../go/types
# Remove l8common generated types if present — those come from vendored dependency
rm -rf ./types/l8common
mv ./types/* ../go/types/.
rm -rf ./types

rm -rf *.rs

cd ../go
find . -name "*.go" -type f -exec sed -i 's|"./types/l8services"|"github.com/saichler/l8types/go/types/l8services"|g' {} +
find . -name "*.go" -type f -exec sed -i 's|"./types/l8api"|"github.com/saichler/l8types/go/types/l8api"|g' {} +
find . -name "*.go" -type f -exec sed -i 's|"./types/l8common"|"github.com/saichler/l8common/go/types/l8common"|g' {} +
find . -name "*.go" -type f -exec sed -i 's|"./types/vend"|"github.com/saichler/l8vendingmachine/go/types/vend"|g' {} +
