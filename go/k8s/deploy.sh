#!/usr/bin/env bash
kubectl apply -f ./vend-vnet.yaml
kubectl apply -f ./vend-logs-vnet.yaml
sleep 2
kubectl apply -f ./vend.yaml
kubectl apply -f ./vend-web.yaml
kubectl apply -f ./vend-inv.yaml
kubectl apply -f ./vend-collector.yaml
kubectl apply -f ./vend-parser.yaml
kubectl apply -f ./vend-log-agent.yaml
