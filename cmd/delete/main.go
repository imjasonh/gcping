package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	project = flag.String("project", "gcping-1369", "Project to use")
	tok     = flag.String("tok", "", "Auth token")
)

func main() {
	flag.Parse()

	if *tok == "" {
		log.Fatalf("Must provide -tok")
	}

	// List all instances globally.
	urls, err := list()
	if err != nil {
		log.Fatalf("Listing instances: %v", err)
	}

	// Delete all instances in parallel.
	var g errgroup.Group
	for _, u := range urls {
		u := u
		g.Go(func() error { return delete(u) })
	}
	if err := g.Wait(); err != nil {
		log.Fatalf("Delete: %v", err)
	}

	// Poll until all instances are deleted.
	for {
		urls, err := list()
		if err != nil {
			log.Fatalf("Listing after delete: %v", err)
		}
		log.Printf("Found %d instances...", len(urls))
		if len(urls) == 0 {
			break
		}
		time.Sleep(10 * time.Second)
	}
}

func list() ([]string, error) {
	url := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/aggregated/instances", *project)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("NewRequest: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+*tok)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		all, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(all))
	}
	var response struct {
		Items map[string]struct {
			Instances []struct {
				SelfLink string `json:"selfLink"`
			} `json:"instances"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("json.Decode: %v", err)
	}

	var instances []string
	for _, i := range response.Items {
		for _, ii := range i.Instances {
			instances = append(instances, ii.SelfLink)
		}
	}
	sort.Strings(instances)
	return instances, nil
}

func delete(url string) error {
	log.Println("Deleting", url)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("NewRequest: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+*tok)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		all, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(all))
	}
	return nil
}
