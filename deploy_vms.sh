#!/bin/bash

export CLOUDSDK_CORE_PROJECT=gcping-1369
REGIONS="us-central1 us-east1 us-west1 europe-west1 asia-east1 asia-northeast1 asia-southeast-1"

# Delete network and recreate it with subnets for each region.
NETWORK_NAME=network
SUBNET_NAME=subnet
gcloud compute networks delete $NETWORK_NAME
gcloud compute networks create $NETWORK_NAME \
  --mode=custom \
  --description="Non-default network"

part=100
for r in $REGIONS; do
  gcloud compute networks subnet create $SUBNET_NAME \
    --network=$NETWORK \
    --range="10.$part.0.0/0" \
  part=$((part+2))
done

# Delete and recreate VMs.
for r in $REGIONS; do
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
done
