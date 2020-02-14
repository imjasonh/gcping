package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	p := os.Getenv("PORT")
	if p == "" {
		p = "80"
	}
	log.Printf("Port is %s", p)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, "pong")
	})
	log.Fatal(http.ListenAndServe(":"+p, nil))
}
