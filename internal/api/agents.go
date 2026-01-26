package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ListAgents fetches agents in a project with pagination
func (c *Client) ListAgents(orgName, projectName string, opts ListOptions) ([]AgentResponse, int, error) {
	params := url.Values{}
	buildPaginationQuery(params, opts)

	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents"
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

	var listResp AgentListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, 0, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Agents, listResp.Total, nil
}

// GetAgent fetches a specific agent
func (c *Client) GetAgent(orgName, projectName, agentName string) (*AgentResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var agent AgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &agent, nil // Return pointer to the agent
}

// DeleteAgent deletes an agent from a project
func (c *Client) DeleteAgent(orgName, projectName, agentName string) error {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName
	return c.doDelete(path)
}

// CreateAgent creates a new agent in a project
func (c *Client) CreateAgent(orgName, projectName string, req CreateAgentRequest) (*AgentResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents"

	resp, err := c.doRequestWithBody("POST", path, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// API returns 202 Accepted for agent creation
	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var agent AgentResponse
	if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &agent, nil
}

// GenerateAgentToken generates a JWT token for an agent
func (c *Client) GenerateAgentToken(orgName, projectName, agentName string, req *TokenRequest) (*TokenResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/token"

	resp, err := c.doRequestWithBody("POST", path, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &tokenResp, nil
}

// GetAgentRuntimeLogs fetches runtime logs for a deployed agent
func (c *Client) GetAgentRuntimeLogs(orgName, projectName, agentName string, req RuntimeLogRequest) (*LogsResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/runtime-logs"

	resp, err := c.doRequestWithBody("POST", path, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var logsResp LogsResponse
	if err := json.NewDecoder(resp.Body).Decode(&logsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &logsResp, nil
}

// GetAgentMetrics fetches resource metrics for a deployed agent
func (c *Client) GetAgentMetrics(orgName, projectName, agentName string, req MetricsFilterRequest) (*MetricsResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/metrics"

	resp, err := c.doRequestWithBody("POST", path, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var metricsResp MetricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&metricsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &metricsResp, nil
}
