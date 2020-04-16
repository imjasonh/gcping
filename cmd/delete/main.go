package main

import (
	"context"
	"flag"
	"log"

	"github.com/ImJasonH/gcping/pkg/util"
	compute "google.golang.org/api/compute/v1"
)

var project = flag.String("project", "gcping-1369", "Project to use")

func main() {
	flag.Parse()

	svc, err := compute.NewService(context.Background())
	if err != nil {
		log.Fatalf("NewService: %v", err)
	}

	// Delete all instances in parallel.
	if err := util.ForEachRegion(svc, *project, deleteVM); err != nil {
		log.Fatal(err)
	}
}

func deleteVM(svc *compute.Service, region string) error {
	zone := region + "-b"
	log.Println("Deleting instance:", region)
	op, err := svc.Instances.Delete(*project, zone, region).Do()
	if err != nil {
		return err
	}
	if err := util.WaitForZoneOp(svc, *project, zone, op); err != nil {
		return err
	}
	log.Printf("instances.delete (%s): ok", region)
	return err
}
