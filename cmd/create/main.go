package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/ImJasonH/gcping/pkg/util"
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

	if err := util.ForEachRegion(svc, *project, create); err != nil {
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

	out, err := util.TemplateExec(`gcloud compute instances create-with-container {{.region}}
--project={{.project}}
--zone={{.region}}-b
--machine-type=f1-micro
--container-image={{.image}}
--container-env=REGION={{.region}}
--tags=http-server
--address={{.addr}}
--network=network
--subnet=subnet
--maintenance-policy=MIGRATE
--boot-disk-size=10
--boot-disk-type=pd-standard
--boot-disk-device-name={{.region}}`, map[string]string{
		"project": *project,
		"region":  region,
		"addr":    addr,
		"image":   *image,
	})
	if err != nil {
		log.Println(string(out))
		if strings.Contains(string(out), "already exists") {
			return nil
		}
		return err
	}
	return nil
}
