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

type RequestClient struct {
	client  HttpClient
	baseURL string
}

var _ RequestSender = &RequestClient{}

func NewHTTPClient(client HttpClient, baseURL string) *RequestClient {
	return &RequestClient{
		client:  client,
		baseURL: baseURL,
	}
}

// HttpClient is an interface used for mocking the actual http calls for testing
//
//go:generate mockgen --build_flags=--mod=mod -destination=../../mocks/httpClient.go  . "HttpClient"
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (c *RequestClient) IsServerLive() bool {
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

func (c *RequestClient) SendCreateRequest(path string, data []byte, isDirectory bool) error {
	type createRequestBody struct {
		Path        string `json:"path"`
		Data        []byte `json:"data"`
		IsDirectory bool   `json:"isDirectory"`
	}
	requestBody := createRequestBody{
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

func (c *RequestClient) SendDeleteRequest(path string) error {
	type deleteRequestBody struct {
		Path string
	}
	requestBody := deleteRequestBody{
		Path: path,
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	if err != nil {
		slog.Debug("unable to marshal request body to byte array", "err", err)
		return err
	}
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/v1/file", c.baseURL), bytes.NewReader(requestBodyBytes))
	if err != nil {
		slog.Debug("error creating request", "err", err)
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

func (c *RequestClient) SendRenameRequest(oldPath, newPath string) error {
	type renameRequestBody struct {
		Path         string
		PreviousPath string
	}
	requestBody := renameRequestBody{
		Path:         newPath,
		PreviousPath: oldPath,
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	if err != nil {
		slog.Debug("unable to marshal request body to byte array", "err", err)
		return err
	}
	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/v1/file", c.baseURL), bytes.NewReader(requestBodyBytes))
	if err != nil {
		slog.Debug("error creating request", "err", err)
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

func (c *RequestClient) SendUpdateRequest(path string, data []byte) error {
	type updateRequestBody struct {
		Path string
		Data []byte
	}
	requestBody := updateRequestBody{
		Path: path,
		Data: data,
	}

	requestBodyBytes, err := json.Marshal(&requestBody)
	if err != nil {
		slog.Debug("unable to marshal request body to byte array", "err", err)
		return err
	}
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/v1/file", c.baseURL), bytes.NewReader(requestBodyBytes))
	if err != nil {
		slog.Debug("error creating request", "err", err)
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
