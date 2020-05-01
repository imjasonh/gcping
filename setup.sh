#!/bin/bash

set -x

export CLOUDSDK_CORE_PROJECT=gcping-1369
NETWORK_NAME=network
SUBNET_NAME=subnet

# Ensure a static IP address exists for each region.
# Also write list of regions to regions.txt
ensureAddrs() {
  go run cmd/addresses/main.go || exit 1
}

# Delete network and recreate it with subnets for each region.
recreateNetwork() {
  # Delete network and firewall rule.
  FWR_NAME=network-allow-http
  gcloud -q compute firewall-rules delete $FWR_NAME || true
  gcloud -q compute networks delete $NETWORK_NAME || true

  # Create network with subnets and firewall rule..
  gcloud compute networks create $NETWORK_NAME \
    --subnet-mode=custom \
    --description="Non-default network"
  gcloud compute firewall-rules create $FWR_NAME \
    --network=network \
    --allow=tcp:80 \
    --target-tags=http-server

  go run ./cmd/networks/main.go || exit 1
}

# Delete VMs in each region.
deleteVMs() {
  go run cmd/delete/main.go || exit 1
}

# Create VMs in each region.
createVMs() {
  # Build the image.
  image=$(KO_DOCKER_REPO=gcr.io/${CLOUDSDK_CORE_PROJECT} ko publish -B ./cmd/ping/)
  go run cmd/create/main.go -image=${image} || exit 1
}

# Delete and create the global loadbalancer with a instance group for each region. 
recreateLB() {
  gcloud -q compute forwarding-rules    delete http-content-rule --global
  gcloud -q compute target-http-proxies delete http-lb-proxy
  gcloud -q compute url-maps            delete web-map
  gcloud -q compute backend-services    delete backend-service --global
  gcloud -q compute http-health-checks  delete http-basic-check

  # Create LB.
  lb_addr=$(gcloud compute addresses describe global --global | grep "address: " | cut -d' ' -f2)
  if [[ -z $lb_addr ]]; then
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
    --address=$lb_addr \
    --global \
    --target-http-proxy=http-lb-proxy \
    --ports=80

  go run ./cmd/lb/main.go || exit 1

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
}

regenConfig() {
  go run cmd/regen/main.go || exit 1
}

uploadPages() {
  BUCKET=gs://www.gcping.com
  gsutil cp config.js ${BUCKET}
  gsutil cp descs.js ${BUCKET}
  gsutil cp index.html ${BUCKET}
  gsutil cp icon.png ${BUCKET}
  gsutil acl ch -u AllUsers:R ${BUCKET}/*
  gsutil web set -m index.html ${BUCKET}
}

ensureAddrs
deleteVMs
recreateNetwork
createVMs
recreateLB
regenConfig
uploadPages

