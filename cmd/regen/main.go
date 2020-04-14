package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	compute "google.golang.org/api/compute/v1"
)

var (
	project = flag.String("project", "gcping-1369", "Project to use")
	outFile = flag.String("out", "config.js", "Output file")
)

var tmpl = template.Must(template.New("name").Parse(`
var _URLS = {
{{range .}}  "{{.Region}}": "http://{{.IP}}/ping",
{{end}}};
`))

func main() {
	flag.Parse()

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
	svc, err := compute.NewService(context.Background())
	resp, err := svc.Addresses.AggregatedList(*project).Do()
	if err != nil {
		return nil, err
	}

	var addresses []address
	for reg, addrs := range resp.Items {
		reg = strings.TrimPrefix(reg, "regions/")
		addresses = append(addresses, address{reg, addrs.Addresses[0].Address})
	}
	sort.Slice(addresses, func(i, j int) bool { return addresses[i].Region < addresses[j].Region })
	return addresses, nil
}
