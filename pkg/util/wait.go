package util

import (
	"fmt"

	compute "google.golang.org/api/compute/v1"
)

func WaitForZoneOp(svc *compute.Service, project, zone string, op *compute.Operation) error {
	var err error
	for {
		op, err = svc.ZoneOperations.Wait(project, zone, op.Name).Do()
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

func WaitForRegionOp(svc *compute.Service, project, region string, op *compute.Operation) error {
	var err error
	for {
		op, err = svc.RegionOperations.Wait(project, region, op.Name).Do()
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

func WaitForGlobalOp(svc *compute.Service, project string, op *compute.Operation) error {
	var err error
	for {
		op, err = svc.GlobalOperations.Wait(project, op.Name).Do()
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
