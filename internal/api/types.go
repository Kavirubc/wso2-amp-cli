package api

import "time"

type AgentResponse struct {
	UUID           string         `json:"uuid,omitempty"`
	Name           string         `json:"name"`
	DisplayName    string         `json:"displayName,omitempty"`
	Description    string         `json:"description,omitempty"`
	ProjectName    string         `json:"projectName"`
	Status         string         `json:"status,omitempty"`
	CreatedAt      time.Time      `json:"createdAt"`
	Provisioning   *Provisioning  `json:"provisioning,omitempty"`
	AgentType      *AgentTypeInfo `json:"agentType,omitempty"`
	RuntimeConfigs *RuntimeConfig `json:"runtimeConfigs,omitempty"`
	Language       string         `json:"language,omitempty"`
}

// Provisioning contains agent source configuration
type Provisioning struct {
	Type       string            `json:"type"`
	Repository *RepositoryConfig `json:"repository,omitempty"`
}

// RepositoryConfig holds git repository details for agent source
type RepositoryConfig struct {
	URL     string `json:"url"`
	Branch  string `json:"branch"`
	AppPath string `json:"appPath,omitempty"`
}

// AgentTypeInfo describes the agent's type classification
type AgentTypeInfo struct {
	Type    string `json:"type"`
	SubType string `json:"subType,omitempty"`
}

// RuntimeConfig holds runtime environment settings
type RuntimeConfig struct {
	Language        string                `json:"language,omitempty"`
	LanguageVersion string                `json:"languageVersion,omitempty"`
	RunCommand      string                `json:"runCommand,omitempty"`
	Env             []EnvironmentVariable `json:"env,omitempty"`
}

// OrganizationResponse represents an organization
type OrganizationResponse struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName,omitempty"`
	Description string    `json:"description,omitempty"`
	Namespace   string    `json:"namespace,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ProjectResponse represents a project
type ProjectResponse struct {
	UUID               string    `json:"uuid,omitempty"`
	Name               string    `json:"name"`
	OrgName            string    `json:"orgName"`
	DisplayName        string    `json:"displayName,omitempty"`
	Description        string    `json:"description,omitempty"`
	DeploymentPipeline string    `json:"deploymentPipeline,omitempty"`
	CreatedAt          time.Time `json:"createdAt"`
}

type ProjectListResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
	Total    int               `json:"total"`
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

// BuildDetailsResponse extends BuildResponse with step information and progress
type BuildDetailsResponse struct {
	BuildResponse              // Embed base build fields
	Percent         float64     `json:"percent,omitempty"`
	Steps           []BuildStep `json:"steps,omitempty"`
	DurationSeconds int         `json:"durationSeconds,omitempty"`
}

// BuildStep represents a single step in the build process
type BuildStep struct {
	Type       string `json:"type"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	StartedAt  string `json:"startedAt,omitempty"`
	FinishedAt string `json:"finishedAt,omitempty"`
}

// BuildLogsResponse contains build logs with metadata
type BuildLogsResponse struct {
	Logs       []LogEntry `json:"logs"`
	TotalCount int        `json:"totalCount"`
	TookMs     float64    `json:"tookMs"`
}

// LogEntry represents a single log line
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Log       string `json:"log"`
	LogLevel  string `json:"logLevel"`
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

// OrganizationListResponse wraps the paginated orgs response
type OrganizationListResponse struct {
	Organizations []OrganizationResponse `json:"organizations"`
	Limit         int                    `json:"limit"`
	Offset        int                    `json:"offset"`
	Total         int                    `json:"total"`
}

// DeploymentPipelineResponse represents a deployment pipeline
type DeploymentPipelineResponse struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"displayName,omitempty"`
	Description string    `json:"description,omitempty"`
	OrgName     string    `json:"orgName"`
	CreatedAt   time.Time `json:"createdAt"`
}

// DeploymentPipelineListResponse wraps paginated pipelines response
type DeploymentPipelineListResponse struct {
	DeploymentPipelines []DeploymentPipelineResponse `json:"deploymentPipelines"`
	Limit               int                          `json:"limit"`
	Offset              int                          `json:"offset"`
	Total               int                          `json:"total"`
}

// --- Request Types ---

// CreateProjectRequest for POST /orgs/{org}/projects
type CreateProjectRequest struct {
	Name               string  `json:"name"`
	DisplayName        string  `json:"displayName"`
	Description        *string `json:"description,omitempty"`
	DeploymentPipeline string  `json:"deploymentPipeline"`
}

// DeployAgentRequest for POST /orgs/{org}/projects/{proj}/agents/{agent}/deployments
type DeployAgentRequest struct {
	ImageId string                `json:"imageId"`
	Env     []EnvironmentVariable `json:"env,omitempty"`
}

// EnvironmentVariable represents a key-value pair for deployment config
type EnvironmentVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateAgentRequest for POST /orgs/{org}/projects/{proj}/agents
type CreateAgentRequest struct {
	Name           string          `json:"name"`
	DisplayName    string          `json:"displayName"`
	Description    string          `json:"description,omitempty"`
	Provisioning   Provisioning    `json:"provisioning"`
	AgentType      AgentTypeInfo   `json:"agentType"`
	RuntimeConfigs *RuntimeConfig  `json:"runtimeConfigs,omitempty"`
	InputInterface *InputInterface `json:"inputInterface,omitempty"`
}

// InputInterface defines the agent's input interface configuration
type InputInterface struct {
	Type     string        `json:"type"`
	Port     int           `json:"port,omitempty"`
	BasePath string        `json:"basePath,omitempty"`
	Schema   *SchemaConfig `json:"schema,omitempty"`
}

// SchemaConfig holds OpenAPI schema configuration
type SchemaConfig struct {
	Path string `json:"path"`
}

// TokenRequest for POST /orgs/{org}/projects/{proj}/agents/{agent}/token
type TokenRequest struct {
	ExpiresIn string `json:"expiresIn,omitempty"`
}

// TokenResponse represents a generated agent token
type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expiresAt"`
	IssuedAt  int64  `json:"issuedAt"`
	TokenType string `json:"tokenType"`
}
