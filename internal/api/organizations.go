package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ListOrganizations fetches organizations with pagination
func (c *Client) ListOrganizations(opts ListOptions) ([]OrganizationResponse, int, error) {
	params := url.Values{}
	buildPaginationQuery(params, opts)

	path := "/orgs"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var listResp OrganizationListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Organizations, listResp.Total, nil
}

// GetOrganization fetches a single organization by name
func (c *Client) GetOrganization(orgName string) (*OrganizationResponse, error) {
	resp, err := c.doRequest("GET", "/orgs/"+orgName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var org OrganizationResponse
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &org, nil
}

// CreateOrganization creates a new organization
func (c *Client) CreateOrganization(req CreateOrganizationRequest) (*OrganizationResponse, error) {
	resp, err := c.doRequestWithBody("POST", "/orgs", req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// API returns 202 Accepted for async organization creation
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var org OrganizationResponse
	if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &org, nil
}
