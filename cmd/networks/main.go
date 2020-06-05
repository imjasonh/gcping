package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ImJasonH/gcping/pkg/util"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

var (
	project = flag.String("project", "gcping-1369", "Project to use")
)

func main() {
	flag.Parse()

	svc, err := compute.NewService(context.Background())
	if err != nil {
		log.Fatalf("NewService: %v", err)
	}

	// Create subnet in each region.
	if err := util.ForEachRegion(svc, *project, createSubnet); err != nil {
		log.Fatal(err)
	}
}

func createSubnet(svc *compute.Service, region string) error {
	part := 40
	for i := 0; i < 40; i++ {
		start := time.Now()
		op, err := svc.Subnetworks.Insert(*project, region, &compute.Subnetwork{
			Name:        "subnet",
			Network:     fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/global/networks/%s", *project, "network"),
			IpCidrRange: fmt.Sprintf("10.%d.0.0/20", part),
		}).Do()
		if herr, ok := err.(*googleapi.Error); ok && herr.Code == http.StatusConflict {
			// Already exists.
			log.Printf("subnet.create (%s): already exists", region)
			return nil
		}
		if err != nil {
			log.Printf("subnet.insert (%s): %v", region, err)
			part += 2
			continue
		}
		if err := util.WaitForRegionOp(svc, *project, region, op); err != nil {
			return err
		}
		log.Printf("subnet.create (%s): ok, took %s", region, time.Since(start))
		break
	}
	return nil
}
