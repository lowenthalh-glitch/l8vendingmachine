#!/usr/bin/env bash
kubectl delete -f ./vend-parser.yaml
kubectl delete -f ./vend-collector.yaml
kubectl delete -f ./vend-inv.yaml
kubectl delete -f ./vend-web.yaml
kubectl delete -f ./vend.yaml
kubectl delete -f ./vend-vnet.yaml
