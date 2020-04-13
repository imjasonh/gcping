package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
		log.Fatalf("Must provide -tok")
	}

	of, err := os.Create(*outFile)
	if err != nil {
		log.Fatalf("os.Open(%s): %v", *outFile, err)
	}
	defer of.Close()

	addrs, err := addresses()
	if err != nil {
		log.Fatal(err)
	}
	if err := tmpl.Execute(io.MultiWriter(os.Stdout, of), addrs); err != nil {
		log.Fatalf("tmpl.Execute: %v", err)
	}
}

type address struct{ Region, IP string }

func addresses() ([]address, error) {
	// Get instances and IPs.
	url := fmt.Sprintf("https://www.googleapis.com/compute/v1/projects/%s/aggregated/addresses", *project)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("NewRequest: GET %s: %v", url, err)
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
			Addresses []struct {
				Address string `json:"address"`
			} `json:"addresses"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("json.Decode: %v", err)
	}

	var addresses []address
	for reg, addrs := range response.Items {
		reg = strings.TrimPrefix(reg, "regions/")
		addresses = append(addresses, address{reg, addrs.Addresses[0].Address})
	}
	sort.Slice(addresses, func(i, j int) bool { return addresses[i].Region < addresses[j].Region })
	return addresses, nil
}
