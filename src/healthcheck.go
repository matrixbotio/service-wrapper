package main

import (
	"net/http"
	"os"
)

func healthcheck() {
	var client http.Client
	resp, err := client.Get("http://localhost/health")
	if err != nil {
		println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		println("Invalid HTTP status")
		os.Exit(1)
	}
}
