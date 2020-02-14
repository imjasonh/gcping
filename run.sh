#!/usr/bin/env bash

set -euxo pipefail

export CLOUDSDK_CORE_PROJECT=gcping-1369

# Build the image.
image=$(KO_DOCKER_REPO=gcr.io/${CLOUDSDK_CORE_PROJECT} ko publish -B ./cmd/ping/)

while read region; do
  gcloud beta run deploy ping \
    --platform=managed \
    --region=${region} \
    --allow-unauthenticated \
    --update-env-vars=REGION=${region} \
    --image=${image} || echo not a Cloud Run region
done < regions.txt
