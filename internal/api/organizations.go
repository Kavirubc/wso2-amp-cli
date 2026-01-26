package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListOrganizations fetches all organizations
func (c *Client) ListOrganizations() ([]OrganizationResponse, error) {
	resp, err := c.doRequest("GET", "/orgs")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // Always close the body when done!

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Decode JSON response into slice of OrganizationResponse
	var listResp OrganizationListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Organizations, nil
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
