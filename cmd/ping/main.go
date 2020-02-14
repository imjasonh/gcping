package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	log.Printf("Port is %s", port)

	region := os.Getenv("REGION")
	if region == "" {
		region = "pong"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-store")
		w.Header().Add("Access-Control-Allow-Origin", "*")
		fmt.Fprintln(w, region)
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
