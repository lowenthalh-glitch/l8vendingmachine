#!/usr/bin/env bash
set -e
docker build --no-cache --platform=linux/amd64 -t saichler/vendmachine-inv:latest .
docker push saichler/vendmachine-inv:latest
