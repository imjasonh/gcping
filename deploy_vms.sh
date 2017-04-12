#!/bin/bash

set -x

export CLOUDSDK_CORE_PROJECT=gcping-1369
REGIONS="us-central1 us-east1 us-west1 europe-west1 asia-east1 asia-northeast1 asia-southeast1"

# Delete network and recreate it with subnets for each region.
FWR_NAME=network-allow-http
NETWORK_NAME=network
SUBNET_NAME=subnet
gcloud -q compute firewall-rules delete $FWR_NAME || true
gcloud -q compute networks delete $NETWORK_NAME || true
gcloud compute networks create $NETWORK_NAME \
  --mode=custom \
  --description="Non-default network"
gcloud compute firewall-rules create network-allow-http \
  --network=network \
  --allow=tcp:80 \
  --target-tags=http-server

part=100
while read r; do
  gcloud compute networks subnets create $SUBNET_NAME \
    --region=$r \
    --network=$NETWORK_NAME \
    --range="10.$part.0.0/20"
  part=$((part+2))
done < regions.txt

# Delete and recreate VMs.
while read r; do
  # b-zones just happen to exist in every region. Let's hope that doesn't
  # change...
  zone=$r-b

  gcloud -q compute instances delete "$r" --zone=$zone || true

  addr=$(gcloud compute addresses describe "$r" --region=$r | grep "address:" | cut -d' ' -f2)
  echo $r $addr

  gcloud compute instances create "$r" \
    --zone=$zone \
    --machine-type=f1-micro \
    --metadata-from-file startup-script=startupscript.sh \
    --network=$NETWORK_NAME \
    --subnet=$SUBNET_NAME \
    --address=$addr \
    --tags=http-server \
    --maintenance-policy=MIGRATE \
    --image-family=ubuntu-1604-lts \
    --image-project=ubuntu-os-cloud \
    --boot-disk-size=10 \
    --boot-disk-type=pd-standard \
    --boot-disk-device-name="$r" \
    --no-scopes \
    --no-service-account
done < regions.txt
