package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type CheckResult struct {
	URL     string
	Status  int
	Elapsed time.Duration
	Err     error
}

func checkURL(client *http.Client, url string) CheckResult {
	start := time.Now()
	resp, err := client.Get(url)
	elapsed := time.Since(start)
	if err != nil {
		return CheckResult{URL: url, Elapsed: elapsed, Err: err}
	}
	defer resp.Body.Close()
	return CheckResult{URL: url, Status: resp.StatusCode, Elapsed: elapsed}
}

func main() {
	urls := os.Args[1:]
