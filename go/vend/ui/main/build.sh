#!/usr/bin/env bash
set -e
docker build --no-cache --platform=linux/amd64 -t saichler/vendmachine-web:latest .
docker push saichler/vendmachine-web:latest
