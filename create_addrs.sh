#!/bin/bash

export CLOUDSDK_CORE_PROJECT=gcping-1369
REGIONS="us-central1 us-east1 us-west1 europe-west1 asia-east1 asia-northeast1"

# Create static addresses for each region.
for r in $REGIONS; do
  gcloud compute addresses create "$r" --region=$r
done

gcloud compute addresses list
