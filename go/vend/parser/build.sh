#!/usr/bin/env bash
set -e
docker build --no-cache --platform=linux/amd64 -t saichler/vendmachine-parser:latest .
docker push saichler/vendmachine-parser:latest
