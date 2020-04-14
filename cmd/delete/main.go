package main

import (
	"context"
	"flag"
	"log"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
	compute "google.golang.org/api/compute/v1"
)

var project = flag.String("project", "gcping-1369", "Project to use")

func main() {
	flag.Parse()

	svc, err := compute.NewService(context.Background())
	if err != nil {
		log.Fatalf("NewService: %v", err)
	}

	// List all instances globally.
	instances, err := list(svc)
	if err != nil {
		log.Fatalf("Listing instances: %v", err)
	}

	// Delete all instances in parallel.
	var g errgroup.Group
	for _, inst := range instances {
		inst := inst
		g.Go(func() error { return delete(svc, inst) })
	}
	if err := g.Wait(); err != nil {
		log.Fatalf("Delete: %v", err)
	}

	// Poll until all instances are deleted.
	for {
		instances, err := list(svc)
		if err != nil {
			log.Fatalf("Listing after delete: %v", err)
		}
		log.Printf("Found %d instances...", len(instances))
		if len(instances) == 0 {
			break
		}
		time.Sleep(10 * time.Second)
	}
}

func list(svc *compute.Service) ([]*compute.Instance, error) {
	resp, err := svc.Instances.AggregatedList(*project).Do()
	if err != nil {
		return nil, err
	}
	var instances []*compute.Instance
	for _, inst := range resp.Items {
		for _, inst := range inst.Instances {
			instances = append(instances, inst)
		}
	}
	return instances, nil
}

func delete(svc *compute.Service, inst *compute.Instance) error {
	zone := inst.Zone[strings.LastIndex(inst.Zone, "/")+1:]
	log.Println("Deleting", zone, inst.Name)
	_, err := svc.Instances.Delete(*project, zone, inst.Name).Do()
	return err
}
