package api

import "time"

type AgentResponse struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName,omitempty"`
	Description string    `json:"description,omitempty"`
	ProjectName string    `json:"projectName"`
	Status      string    `json:"status,omitempty"`
	Language    string    `json:"language,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// OrganizationResponse represents an organization
type OrganizationResponse struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	Status      string    `json:"status,omitempty"`
}

// ProjectResponse represents a project
type ProjectResponse struct {
	Name        string    `json:"name"`
	OrgName     string    `json:"orgName"`
	DisplayName string    `json:"displayName,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	Status      string    `json:"status,omitempty"`
}

// BuildResponse represents a build
type BuildResponse struct {
	Name        string     `json:"name"`
	AgentName   string     `json:"agentName"`
	ProjectName string     `json:"projectName"`
	CommitID    string     `json:"commitId"`
	Status      string     `json:"status"`
	Branch      string     `json:"branch,omitempty"`
	StartedAt   time.Time  `json:"startedAt"`
	EndedAt     *time.Time `json:"endedAt,omitempty"`
}

// DeploymentResponse represents a deployment
type DeploymentResponse struct {
	AgentName      string    `json:"agentName"`
	ProjectName    string    `json:"projectName"`
	Status         string    `json:"status"`
	Environment    string    `json:"environment"`
	LastDeployedAt time.Time `json:"lastDeployedAt"`
}

// AgentListResponse wraps the paginated agents response
type AgentListResponse struct {
	Agents []AgentResponse `json:"agents"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Total  int             `json:"total"`
}
