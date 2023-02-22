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
	if len(urls) == 0 {
		urls = []string{
			"http://127.0.0.1:8000",
			"https://example.com",
		}
	}

	client := &http.Client{Timeout: 5 * time.Second}
	results := make([]CheckResult, len(urls))
	var wg sync.WaitGroup

	for i, u := range urls {
		wg.Add(1)
		go func(idx int, target string) {
			defer wg.Done()
			results[idx] = checkURL(client, target)
		}(i, u)
	}
	wg.Wait()

	for _, r := range results {
		if r.Err != nil {
			fmt.Printf("[FAIL] %-32s | unreachable (%v)\n", r.URL, r.Err)
		} else if r.Status != http.StatusOK {
			fmt.Printf("[WARN] %-32s | status %d (%dms)\n", r.URL, r.Status, r.Elapsed.Milliseconds())
		} else {
			fmt.Printf("[ OK ] %-32s | status %d (%dms)\n", r.URL, r.Status, r.Elapsed.Milliseconds())
		}
	}
}
