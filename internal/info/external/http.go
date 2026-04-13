package external

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const requestTimeout = 30 * time.Second

var httpClient = &http.Client{
	Timeout: requestTimeout,
}

func getHTTPResponse(url string, result interface{}) error {
	bodyBytes, err := getHTTPResponseBytes(url)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	return nil
}

func getHTTPResponseBytes(url string) ([]byte, error) {
	resp, err := httpClient.Get(url) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unsuccessful status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return bodyBytes, nil
}
