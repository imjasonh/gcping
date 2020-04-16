package main

import (
	"context"
	"flag"
	"fmt"
	"log"

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

	resp, err := svc.Regions.List(*project).Do()
	if err != nil {
		log.Fatalf("regions.list: %v", err)
	}

	// Create unmanaged IGs in each region.
	var g errgroup.Group
	for _, r := range resp.Items {
		r := r.Name
		g.Go(func() error { return create(svc, r) })
	}
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	// Update the global LB with all IG backends.
	bs, err := svc.BackendServices.Get(*project, "backend-service").Do()
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range resp.Items {
		zone := r.Name + "-b"
		ig := "instance-group-" + r.Name
		bs.Backends = append(bs.Backends, &compute.Backend{
			Group:          fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/zones/%s/instanceGroups/%s", *project, zone, ig),
			BalancingMode:  "UTILIZATION",
			MaxUtilization: 0.8,
			CapacityScaler: 1,
		})
	}
	op, err := svc.BackendServices.Update(*project, "backend-service", bs).Do()
	if err != nil {
		log.Fatal(err)
	}
	if err := waitForGlobalOp(svc, op); err != nil {
		log.Fatal(err)
	}
	log.Println("backendServices.update: ok")
}

func create(svc *compute.Service, region string) error {
	log.Println("Creating LB config for", region)

	// Create an (unmanaged) instance group.
	zone := region + "-b"
	ig := "instance-group-" + region
	op, err := svc.InstanceGroups.Insert(*project, zone, &compute.InstanceGroup{
		Name: ig,
	}).Do()
	if err != nil {
		return fmt.Errorf("instanceGroups.insert (%s): %v", region, err)
	}
	if err := waitForZoneOp(svc, zone, op); err != nil {
		return fmt.Errorf("wait for instanceGroups.insert (%s): %v", region, err)
	}
	log.Printf("instanceGroups.insert (%s): ok", region)

	// Add the instance to the IG.
	op, err = svc.InstanceGroups.AddInstances(*project, zone, ig, &compute.InstanceGroupsAddInstancesRequest{
		Instances: []*compute.InstanceReference{{
			Instance: fmt.Sprintf("zones/%s/instances/%s", zone, region),
		}},
	}).Do()
	if err != nil {
		return fmt.Errorf("instanceGroups.addInstances (%s): %v", region, err)
	}
	if err := waitForZoneOp(svc, zone, op); err != nil {
		return fmt.Errorf("wait for instanceGroups.addInstances (%s): %v", region, err)
	}
	log.Printf("instanceGroups.addInstances (%s): ok", region)

	// Set named ports.
	op, err = svc.InstanceGroups.SetNamedPorts(*project, zone, ig, &compute.InstanceGroupsSetNamedPortsRequest{
		NamedPorts: []*compute.NamedPort{{Name: "http", Port: 80}},
	}).Do()
	if err != nil {
		return fmt.Errorf("instanceGroups.setNamedPorts (%s): %v", region, err)
	}
	if err := waitForZoneOp(svc, zone, op); err != nil {
		return fmt.Errorf("wait for instanceGroups.setNamedPorts (%s): %v", region, err)
	}
	log.Printf("instanceGroups.setNamedPorts (%s): ok", region)
	return nil
}

func waitForZoneOp(svc *compute.Service, zone string, op *compute.Operation) error {
	var err error
	for {
		op, err = svc.ZoneOperations.Wait(*project, zone, op.Name).Do()
		if err != nil {
			return err
		}
		if op.Error != nil {
			return fmt.Errorf("Operation error: %v", op.Error)
		}
		if op.Status == "DONE" {
			return nil
		}
	}
}

func waitForGlobalOp(svc *compute.Service, op *compute.Operation) error {
	var err error
	for {
		op, err = svc.GlobalOperations.Wait(*project, op.Name).Do()
		if err != nil {
			return err
		}
		if op.Error != nil {
			return fmt.Errorf("Operation error: %v", op.Error)
		}
		if op.Status == "DONE" {
			return nil
		}
	}

}
