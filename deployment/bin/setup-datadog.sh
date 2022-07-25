#!/usr/bin/env bash

helm repo add datadog https://helm.datadoghq.com
helm repo update

kubectl create namespace datadog

helm install dd-agent \
    --set datadog.site='datadoghq.eu' \
    --set datadog.apiKey="${DATADOG_API_KEY}" \
    --set datadog.logs.enabled=true \
    --set datadog.apm.portEnabled=true \
    --set datadog.logs.containerCollectAll=true \
    datadog/datadog --namespace datadog
