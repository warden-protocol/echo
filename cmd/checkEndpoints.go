package main

import (
	"context"
	"net/http"
	"sync"
	"time"
)

const (
	clientTimeout = 5
)

func checkEndpoints(endpoints []string) map[string]EndpointResult {
	results := make(map[string]EndpointResult)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, url := range endpoints {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			// Create a context with a timeout
			ctx, cancel := context.WithTimeout(context.Background(), clientTimeout*time.Second)
			defer cancel()

			// Create the HTTP request with context
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
			if err != nil {
				res := EndpointResult{
					StatusCode: 0,
					Reachable:  false,
					Error:      err.Error(),
				}
				mu.Lock()
				results[u] = res
				mu.Unlock()
				return
			}

			// Use the default HTTP client or customize if needed
			client := &http.Client{}

			// Perform the HTTP request
			resp, err := client.Do(req)
			var res EndpointResult
			if err != nil {
				res = EndpointResult{
					StatusCode: 0,
					Reachable:  false,
					Error:      err.Error(),
				}
			} else {
				res = EndpointResult{
					StatusCode: resp.StatusCode,
					Reachable:  resp.StatusCode >= 200 && resp.StatusCode < 300,
				}
				resp.Body.Close()
			}
			mu.Lock()
			results[u] = res
			mu.Unlock()
		}(url)
	}
	wg.Wait()
	return results
}
