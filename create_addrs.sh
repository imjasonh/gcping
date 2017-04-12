#!/bin/bash

export CLOUDSDK_CORE_PROJECT=gcping-1369

# TODO: Delete addresses?

# Create static addresses for each region.
while read r; do
  gcloud compute addresses create "$r" --region=$r
done < regions.txt

gcloud compute addresses list
