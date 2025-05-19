package adapters

type HTTPClient struct {
}

func NewHTTPClient() *HTTPClient {
	return &HTTPClient{}
}

func (client *HTTPClient) SendUpdateRequest() error {
	return nil
}
