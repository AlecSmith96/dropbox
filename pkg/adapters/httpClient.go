package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type HTTPClient struct {
	client  *http.Client
	baseURL string
}

func NewHTTPClient(client *http.Client, baseURL string) *HTTPClient {
	return &HTTPClient{
		client:  client,
		baseURL: baseURL,
	}
}

func (c *HTTPClient) IsServerLive() bool {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/v1/health/live", c.baseURL), nil)
	if err != nil {
		slog.Debug("error creating request", "err", err)
		return false
	}

	response, err := c.client.Do(req)
	if err != nil {
		slog.Debug("error sending create request", "err", err)
		return false
	}

	if response.StatusCode != 200 {
		return false
	}

	return true
}

func (c *HTTPClient) SendCreateRequest(path string, data []byte, isDirectory bool) error {
	type CreateRequestBody struct {
		Path        string `json:"path"`
		Data        []byte `json:"data"`
		IsDirectory bool   `json:"isDirectory"`
	}
	requestBody := CreateRequestBody{
		Path:        path,
		Data:        data,
		IsDirectory: isDirectory,
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	if err != nil {
		slog.Debug("unable to marshal request body to byte array", "err", err)
		return err
	}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/file", c.baseURL), bytes.NewReader(requestBodyBytes))
	if err != nil {
		slog.Debug("error creating request", "err", err)
		return err
	}

	response, err := c.client.Do(req)
	if err != nil {
		slog.Debug("error sending create request", "err", err)
		return err
	}

	// we can just check for != 200 here as we know the server doesnt return any other success codes (2**)
	if response.StatusCode != http.StatusOK {
		slog.Debug("request failed with status code", "statusCode", response.StatusCode)
		return fmt.Errorf("request failed with status code %d", response.StatusCode)
	}

	return nil
}

func (c *HTTPClient) SendDeleteRequest(path string) error {
	type CreateRequestBody struct {
		Path string
	}
	requestBody := CreateRequestBody{
		Path: path,
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	if err != nil {
		slog.Debug("unable to marshal request body to byte array", "err", err)
		return err
	}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/file", c.baseURL), bytes.NewReader(requestBodyBytes))
	if err != nil {
		fmt.Printf("c: could not create request: %s\n", err)
		os.Exit(1)
	}

	response, err := c.client.Do(req)
	if err != nil {
		slog.Debug("error sending create request", "err", err)
		return err
	}

	// we can just check for != 200 here as we know the server doesnt return any other success codes (2**)
	if response.StatusCode != http.StatusOK {
		slog.Error("request failed with status code", "statusCode", response.StatusCode)
		return errors.New("request failed")
	}

	return nil
}
