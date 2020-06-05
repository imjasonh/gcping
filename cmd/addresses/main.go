package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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

	var regions []string
	resp, err := svc.Regions.List(*project).Do()
	if err != nil {
		log.Fatalf("regions.list: %v", err)
	}
	for _, r := range resp.Items {
		regions = append(regions, r.Name)
	}

	// Write regions.txt
	if err := writeFile(regions); err != nil {
		log.Fatal(err)
	}

	if err := ensureAddrs(svc, regions); err != nil {
		log.Fatal(err)
	}
}

func writeFile(regions []string) error {
	f, err := os.Create("regions.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprintln(f, strings.Join(regions, "\n"))
	return nil
}

func ensureAddrs(svc *compute.Service, regions []string) error {
	want := map[string]struct{}{}
	for _, r := range regions {
		want[r] = struct{}{}
	}

	// For every region, ensure an address exists at that location.
	resp, err := svc.Addresses.AggregatedList(*project).Do()
	if err != nil {
		return err
	}
	have := map[string]struct{}{}
	for r, i := range resp.Items {
		have[strings.TrimPrefix(r, "regions/")] = struct{}{}
		for _, a := range i.Addresses {
			log.Printf("Address (%s): %s", r, a.Address)
		}
	}

	for h := range have {
		delete(want, h)
	}

	for w := range want {
		log.Println("want address for", w)
	}
	if len(want) == 0 {
		return nil
	}

	var g errgroup.Group
	for r := range want {
		log.Println("missing address in region", r)
		r := r
		g.Go(func() error {
			op, err := svc.Addresses.Insert(*project, r, &compute.Address{Name: r}).Do()
			if err != nil {
				return err
			}
			for {
				op, err = svc.RegionOperations.Wait(*project, r, op.Name).Do()
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
		})
	}
	return g.Wait()
}
