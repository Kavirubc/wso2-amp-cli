package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListAgents fetches all agents in a project
func (c *Client) ListAgents(orgName, projectName string) ([]AgentResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents"

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var listResp AgentListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return listResp.Agents, nil
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
