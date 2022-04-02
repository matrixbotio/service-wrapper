package main

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var client = http.Client{
	Timeout: 5 * time.Second,
}

func healthcheck() {
	if isProcessStoppedByWrapper() {
		println("process is stopped by wrapper")
		return
	}
	print("localhost GET 8080 /health ")
	start := time.Now().UnixMilli()
	resp, err := client.Get("http://localhost:8080/health")
	ms := strconv.FormatInt(time.Now().UnixMilli()-start, 10) + "ms"
	if err != nil {
		println("-1 " + ms)
		println(err.Error())
		os.Exit(1)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			println(err.Error())
			os.Exit(1)
		}
	}(resp.Body)

	println(strconv.Itoa(resp.StatusCode) + " " + ms)

	if resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}
}
