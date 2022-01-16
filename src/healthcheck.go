package main

import (
	"io"
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

	if resp.StatusCode == http.StatusOK {
		_, err := io.ReadAll(resp.Body)
		if err != nil {
			println(err)
			os.Exit(1)
		}
	}
}
