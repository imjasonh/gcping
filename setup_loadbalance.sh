#!/bin/bash

export CLOUDSDK_CORE_PROJECT=gcping-1369
REGIONS="us-central1 us-east1 us-west1 europe-west1 asia-east1 asia-northeast1 asia-southeast1"

lb_addr=$(gcloud compute addresses describe global --global | grep "address: " | cut -d' ' -f2)
if [[ -n $lb_addr ]]; then
  echo "No IP address found, creating"
  lb_addr=$(gcloud compute addresses create global --global | grep "address: " | cut -d' ' -f2)
fi
echo "Load Balance IP address:" $lb_addr

# Create health check that hits /ping
gcloud compute http-health-checks create http-basic-check \
  --request-path=/ping

# Create backend service using that health check
gcloud compute backend-services create backend-service \
  --protocol=HTTP \
  --http-health-checks=http-basic-check \
  --global \
  --enable-cdn

# Create URL map to map all incoming requests to all instances
gcloud compute url-maps create web-map \
  --default-service=backend-service

# Create target HTTP proxy to route requests to URL map
gcloud compute target-http-proxies create http-lb-proxy \
  --url-map=web-map

# Create global forwarding rule to route requests to HTTP proxy
gcloud compute forwarding-rules create http-content-rule \
  --address=$lb_addr --global \
  --target-http-proxy=http-lb-proxy \
  --ports=80

for r in $REGIONS; do
  ig=instance-group-$r
  zone=$r-b

  # Create instance group for each zone
  gcloud compute instance-groups unmanaged create $ig --zone=$zone

  # Add region's VM to instance group
  gcloud compute instance-groups unmanaged add-instances $ig \
    --instances=$r \
    --zone=$zone

  # Define HTTP service and map a named port
  gcloud compute instance-groups unmanaged set-named-ports $ig \
    --named-ports=http:80 \
    --zone=$zone

  # Add instance groups as backends to backend service
  gcloud compute backend-services add-backend backend-service \
    --balancing-mode UTILIZATION \
    --max-utilization 0.8 \
    --capacity-scaler 1 \
    --instance-group=$ig \
    --instance-group-zone=$zone \
    --global
done

# Ping LB IP until it gets a pong.
while true; do
  got=$(curl http://$lb_addr/ping 2>/dev/null | grep "pong")
  echo $got
  if [[ $got -eq "pong" ]]; then
    break
  fi
  sleep 10
done

echo "Load balance IP:" $lb_addr
