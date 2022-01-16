package main

import (
	"net/http"
	"os"
	"time"
)

func healthcheck() {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
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
