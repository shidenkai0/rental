#!/usr/bin/env bash

# This script is meant to be called from the project Makefile.

cd deployment/helm/api || exit

if [ -n "$CI" ]; then
    if [ -z "$VALUES_PROD_BASE64" ]; then
        echo "VALUES_PROD_BASE64 is not set. Exiting."
        exit 1
    fi
    echo "${VALUES_PROD_BASE64}" | base64 --decode > /tmp/values-prod.yaml
	helm upgrade api --install -f /tmp/values-prod.yaml \
	 --set image.tag="${COMMIT_SHA}" --set image.repository="${DOCKER_REGISTRY}/${IMAGE_NAME}" \
	 . || true # Continue to make sure we remove values-prod.yaml even if the helm command fails.
    rm /tmp/values-prod.yaml
else
	helm upgrade api --install -f values-prod.yaml \
	 --set image.tag="${COMMIT_SHA}" --set image.repository="${DOCKER_REGISTRY}/${IMAGE_NAME}" \
	 .
fi
