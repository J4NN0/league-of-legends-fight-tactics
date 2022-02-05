package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func Get(client *http.Client, url string, destination interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create API request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status not OK, got %d instead", resp.StatusCode)
	}

	if err := unmarshalToInterface(respBody, destination); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return nil
}

func unmarshalToInterface(body []byte, interfaceItem interface{}) error {
	marshalError := json.Unmarshal(body, &interfaceItem)
	if marshalError != nil {
		return fmt.Errorf("failed to unmarshal interface from body: %w", marshalError)
	}
	return nil
}
