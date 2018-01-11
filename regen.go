package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"
)

var (
	project = flag.String("project", "gcping-1369", "Project to use")
	tok     = flag.String("tok", "", "Auth token")
	outFile = flag.String("out", "config.js", "Output file")
)

var tmpl = template.Must(template.New("name").Parse(`
var _URLS = {
{{range .}}  "{{.Region}}": "http://{{.IP}}/ping",
{{end}}};
`))

func main() {
	flag.Parse()

	if *tok == "" {
		exitf("Must provide -tok")
	}

	// Get instances and IPs.
	url := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/aggregated/addresses", *project)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		exitf("NewRequest: GET %s: %v", url, err)
	}
	req.Header.Set("Authorization", "Bearer "+*tok)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		exitf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()
	var response struct {
		Items map[string]struct {
			Addresses []struct {
				Address string `json:"address"`
			} `json:"addresses"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		exitf("json.Decode: %v", err)
	}

	type region struct{ Region, IP string }
	var regions []region
	for reg, addrs := range response.Items {
		reg = strings.TrimPrefix(reg, "regions/")
		regions = append(regions, region{reg, addrs.Addresses[0].Address})
	}
	sort.Slice(regions, func(i, j int) bool { return regions[i].Region < regions[j].Region })

	of, err := os.Create(*outFile)
	if err != nil {
		exitf("os.Open(%s): %v", *outFile, err)
	}
	defer of.Close()
	if err := tmpl.Execute(io.MultiWriter(os.Stdout, of), regions); err != nil {
		exitf("tmpl.Execute: %v", err)
	}
}

func exitf(f string, args ...interface{}) {
	log.Printf(f, args...)
	os.Exit(1)
}
