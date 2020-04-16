package util

import (
	"fmt"

	"golang.org/x/sync/errgroup"
	compute "google.golang.org/api/compute/v1"
)

func ForEachRegion(svc *compute.Service, project string, fn func(svc *compute.Service, region string) error) error {
	resp, err := svc.Regions.List(project).Do()
	if err != nil {
		return fmt.Errorf("regions.list: %v", err)
	}
	var g errgroup.Group
	for _, r := range resp.Items {
		r := r.Name
		g.Go(func() error { return fn(svc, r) })
	}
	return g.Wait()
}
