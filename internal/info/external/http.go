package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

const (
	httpTimeout    = 15 * time.Second
	maxRetries     = 3
	retryBaseDelay = 500 * time.Millisecond
)

var httpClient = &http.Client{ //nolint:gochecknoglobals
	Timeout: httpTimeout,
}

// getJSON fetches the given URL and unmarshals the JSON body into result.
// It retries up to maxRetries times with exponential backoff on transient errors.
func getJSON(url string, result interface{}) error {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(math.Pow(2, float64(attempt-1))) * retryBaseDelay
			time.Sleep(delay)
		}

		lastErr = doGetJSON(url, result)
		if lastErr == nil {
			return nil
		}
	}

	return lastErr
}

// getBytes fetches the given URL and returns the raw response body.
// It retries up to maxRetries times with exponential backoff on transient errors.
func getBytes(url string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(math.Pow(2, float64(attempt-1))) * retryBaseDelay
			time.Sleep(delay)
		}

		data, err := doGetBytes(url)
		if err == nil {
			return data, nil
		}

		lastErr = err
	}

	return nil, lastErr
}

func doGetJSON(url string, result interface{}) error {
	data, err := doGetBytes(url)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(data, result); err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return nil
}

func doGetBytes(url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), httpTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unsuccessful status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}
