package api

import (
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	APIKeyheader string 
	APIKeyValue string
	HTTPClient *http.Client
}

func NewClient(baseURL, apiKeyHeader, apiKeyValue string) *Client {
	return &Client{
		BaseURL: baseURL,
		APIKeyheader: apiKeyHeader,
		APIKeyValue: apiKeyValue,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}