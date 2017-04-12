#!/bin/bash

# Reverse of steps in setup_loadbalance.sh

export CLOUDSDK_CORE_PROJECT=gcping-1369
REGIONS="us-central1 us-east1 us-west1 europe-west1 asia-east1 asia-northeast1 asia-southeast1"

gcloud -q compute forwarding-rules    delete http-content-rule --global
gcloud -q compute target-http-proxies delete http-lb-proxy
gcloud -q compute url-maps            delete web-map
gcloud -q compute backend-services    delete backend-service --global
gcloud -q compute http-health-checks  delete http-basic-check
gcloud -q compute addresses           delete global --global

for r in $REGIONS; do
  ig=instance-group-$r
  zone=$r-b
  gcloud -q compute instance-groups unmanaged delete $ig --zone=$zone
done
