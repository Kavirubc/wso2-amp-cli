package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// ListDeployments fetches all deployments for an agent (returns simple list)
func (c *Client) ListDeployments(orgName, projectName, agentName string) ([]DeploymentResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/deployments"

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var deployments []DeploymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return deployments, nil
}

// GetDeploymentsMap fetches deployments as a map of environment name to details
func (c *Client) GetDeploymentsMap(orgName, projectName, agentName string) (map[string]DeploymentDetails, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/deployments"

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var deployments map[string]DeploymentDetails
	if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return deployments, nil
}

// DeployAgent deploys an agent to an environment
func (c *Client) DeployAgent(orgName, projectName, agentName string, req DeployAgentRequest) error {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/deployments"

	resp, err := c.doRequestWithBody("POST", path, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetAgentEndpoints fetches endpoints for an agent in a specific environment
func (c *Client) GetAgentEndpoints(orgName, projectName, agentName, environment string) ([]EndpointResponse, error) {
	path := "/orgs/" + orgName + "/projects/" + projectName + "/agents/" + agentName + "/endpoints"
	if environment != "" {
		path += "?environment=" + url.QueryEscape(environment)
	}

	resp, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var endpoints []EndpointResponse
	if err := json.NewDecoder(resp.Body).Decode(&endpoints); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return endpoints, nil
}
