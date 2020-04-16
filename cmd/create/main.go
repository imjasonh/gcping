package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"golang.org/x/sync/errgroup"
	compute "google.golang.org/api/compute/v1"
)

var (
	project = flag.String("project", "gcping-1369", "Project to use")
	image   = flag.String("image", "", "Container image to run")
)

func main() {
	flag.Parse()

	if *image == "" {
		log.Fatal("-image is required")
	}

	svc, err := compute.NewService(context.Background())
	if err != nil {
		log.Fatalf("NewService: %v", err)
	}

	resp, err := svc.Regions.List(*project).Do()
	if err != nil {
		log.Fatalf("regions.list: %v", err)
	}
	var g errgroup.Group
	for _, r := range resp.Items {
		r := r.Name
		g.Go(func() error { return create(svc, r) })
	}
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func create(svc *compute.Service, region string) error {
	log.Println("Creating instance:", region)

	// Get the address for this region.
	resp, err := svc.Addresses.Get(*project, region, region).Do()
	if err != nil {
		log.Fatalf("Getting address: %v", err)
	}
	addr := resp.Address
	log.Printf("Address in %s: %s", region, addr)

	parts := strings.Split(strings.ReplaceAll(fmt.Sprintf(`gcloud compute instances create-with-container %s
--project=%s
--zone=%s-b
--machine-type=f1-micro
--container-image=%s
--container-env=REGION=%s
--tags=http-server
--address=%s
--network=network
--subnet=subnet
--maintenance-policy=MIGRATE
--boot-disk-size=10
--boot-disk-type=pd-standard
--boot-disk-device-name=%s`, region, *project, region, *image, region, addr, region), "\n", " "), " ")
	out, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
	if err != nil {
		log.Println(string(out))
		if strings.Contains(string(out), "already exists") {
			return nil
		}
		return err
	}
	return nil
}
