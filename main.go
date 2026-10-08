package main

import (
	"fmt"
	"net/http"
	"time"
)

type HealthResult struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}
type HealthChecker interface {
	Check(url string) HealthResult
}
type HTTPClient struct {
	Client *http.Client
}

func (h HTTPClient) Check(url string) HealthResult {
	start := time.Now()

	response, err := h.Client.Get(url)

	duration := time.Since(start)

	result := HealthResult{
		URL:      url,
		Duration: duration,
		Err:      err,
	}

	if err != nil {
		return result
	}

	defer response.Body.Close()

	result.StatusCode = response.StatusCode

	return result
}
func main() {
	client := HTTPClient{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	var checker HealthChecker = client

	urls := []string{
		"https://go.dev",
		"https://github.com",
		"https://httpbin.org/status/404",
		"https://example.com",
		"https://invalid-fake-host-test.org",
	}

	fmt.Println("=== HTTP URL Health Checker ===")
	fmt.Println()

	for _, url := range urls {
		result := checker.Check(url)

		if result.Err != nil {
			fmt.Printf(
				"-> %s: Помилка (%v, час: %v)\n",
				result.URL,
				result.Err,
				result.Duration,
			)
			continue
		}

		fmt.Printf(
			"-> %s: Status %d (%s, час: %v)\n",
			result.URL,
			result.StatusCode,
			http.StatusText(result.StatusCode),
			result.Duration,
		)
	}
}
