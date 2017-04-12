#!/bin/bash

export CLOUDSDK_CORE_PROJECT=gcping-1369
REGIONS="us-central1 us-east1 us-west1 europe-west1 asia-east1 asia-northeast1 asia-southeast1"

if [[ -n $CREATE_ADDRS ]]; then
  # Create static addresses
  for r in $REGIONS; do
    gcloud compute addresses create "$r" --region=$r
  done
fi

# Create VMs
for r in $REGIONS; do
  zone=$r-b

  gcloud -q compute instances delete "$r" --zone=$zone || true

  addr=$(gcloud compute addresses describe "$r" --region=$r | grep "address:" | cut -d' ' -f2)
  echo $r $addr

  gcloud compute instances create "$r" \
    --zone=$zone \
    --machine-type=f1-micro \
    --metadata-from-file startup-script=startupscript.sh \
    --network=network \
    --subnet=subnet \
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
