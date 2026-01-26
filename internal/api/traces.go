package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// ListTraces fetches traces for an agent
func (c *Client) ListTraces(orgName, projectName, agentName string, opts TraceListOptions) (*TraceListResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/traces"

	// Build query string
	params := url.Values{}
	params.Set("environment", opts.Environment)
	if opts.StartTime != "" {
		params.Set("startTime", opts.StartTime)
	}
	if opts.EndTime != "" {
		params.Set("endTime", opts.EndTime)
	}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}
	if opts.SortOrder != "" {
		params.Set("sortOrder", opts.SortOrder)
	}
	path += "?" + params.Encode()

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result TraceListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// GetTrace fetches a single trace by ID
func (c *Client) GetTrace(orgName, projectName, agentName, traceID, environment string) (*TraceDetailsResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/trace/" + traceID
	path += "?environment=" + url.QueryEscape(environment)

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result TraceDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// ExportTraces exports traces with full span details
func (c *Client) ExportTraces(orgName, projectName, agentName string, opts TraceListOptions) (*TraceExportResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/traces/export"

	// Build query string (same as ListTraces)
	params := url.Values{}
	params.Set("environment", opts.Environment)
	if opts.StartTime != "" {
		params.Set("startTime", opts.StartTime)
	}
	if opts.EndTime != "" {
		params.Set("endTime", opts.EndTime)
	}
	if opts.Limit > 0 {
		params.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Offset > 0 {
		params.Set("offset", strconv.Itoa(opts.Offset))
	}
	path += "?" + params.Encode()

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result TraceExportResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
