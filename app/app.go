package app

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
	"cloud.google.com/go/bigquery"
)

const (
	dataset = "gcping"
	table = "gcping"
)

func init() {
	http.HandleFunc("/", redir)
	http.HandleFunc("/api", handler)
}

func redir(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://www.gcping.com", http.StatusMovedPermanently)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello")
}

func write(ctx context.Context) {
	client := bigquery.NewClient(ctx, appengine.AppID(ctx))
	u := client.Dataset(dataset).Table(table).Uploader()
	u.Put(ctx, nil)
}
