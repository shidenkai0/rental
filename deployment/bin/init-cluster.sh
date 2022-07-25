#!/usr/bin/env bash

parent_path=$( cd "$(dirname "${BASH_SOURCE[0]}")" || return ; pwd -P )

cd "$parent_path" || exit

KUBE_CONTEXT=$(kubectl config current-context)

echo "Initializing kubernetes cluster with context $KUBE_CONTEXT"

./image-pull-secret.sh
./setup-ingress.sh
./setup-datadog.sh
