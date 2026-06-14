#!/usr/bin/env bash
set -e
docker build --no-cache --platform=linux/amd64 -t saichler/vend-logs-vnet:latest .
docker push saichler/vend-logs-vnet:latest
