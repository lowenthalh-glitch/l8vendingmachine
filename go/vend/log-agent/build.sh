#!/usr/bin/env bash
set -e
docker build --no-cache --platform=linux/amd64 -t saichler/vend-log-agent:latest .
docker push saichler/vend-log-agent:latest
