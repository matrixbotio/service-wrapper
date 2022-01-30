package main

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

func healthcheck() {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	if isProcessStoppedByWrapper() {
		println("process is stopped by wrapper")
		return
	}
	print("localhost GET 8080 /health ")
	start := time.Now().UnixMilli()
	resp, err := client.Get("http://localhost:8080/health")
	ms := strconv.FormatInt(time.Now().UnixMilli() - start, 10) + "ms"
	if err != nil {
		println("-1 " + ms)
		println(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	println(strconv.Itoa(resp.StatusCode) + " " + ms)

	if resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}
