#!/usr/bin/env bash
set -euo pipefail

REGISTRY=${REGISTRY:-quay.io/rahulk10}
IMAGE=${IMAGE:-travelportal}
TAG=${TAG:-latest}
RELEASE=${RELEASE:-travelportal}
NAMESPACE=${NAMESPACE:-default}
CHART_DIR=${CHART_DIR:-helm-charts/travelportal}

helm upgrade --install "$RELEASE" "$CHART_DIR" \
  --namespace "$NAMESPACE" --create-namespace \
  --set image.repository="$REGISTRY/$IMAGE" \
  --set image.tag="$TAG"
